package translator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"time"
)

// TranslationMetrics 翻译指标
type TranslationMetrics struct {
	InputTokens  int     `json:"input_tokens"`  // 输入token数
	OutputTokens int     `json:"output_tokens"` // 输出token数
	TotalTokens  int     `json:"total_tokens"`  // 总token数
	ModelLatency float64 `json:"model_latency"` // 模型处理延迟（毫秒）
}

type OpenAITranslator struct {
	provider    string
	apiURL      string
	apiKey      string
	model       string
	timeout     int
	maxTokens   int
	temperature float32
	client      *http.Client
	lastMetrics TranslationMetrics
}

// NewOpenAITranslator 创建新的OpenAI翻译器实例
func NewOpenAITranslator(provider, apiURL, apiKey, model string, timeout, maxTokens int, temperature float32) *OpenAITranslator {
	// 确保默认值合理
	if timeout <= 0 {
		timeout = 30 // 默认30秒超时
	}
	if temperature <= 0 {
		temperature = 0.3 // 默认温度值
	}
	if maxTokens <= 0 {
		maxTokens = 2000 // 默认最大token数
	}

	return &OpenAITranslator{
		provider:    provider,
		apiURL:      apiURL,
		apiKey:      apiKey,
		model:       model,
		timeout:     timeout,
		maxTokens:   maxTokens,
		temperature: temperature,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// Translate 实现翻译功能
func (t *OpenAITranslator) Translate(text, sourceLang, targetLang string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.timeout)*time.Second)
	defer cancel()

	return t.TranslateWithContext(ctx, text, sourceLang, targetLang)
}

// TranslateWithContext 支持上下文的翻译方法
func (t *OpenAITranslator) TranslateWithContext(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	log.Println(t.apiURL, t.model)
	// 构造翻译提示
	prompt := fmt.Sprintf("Translate the following text from %s to %s. Only return the translated text without any explanations:\n\n%s",
		sourceLang, targetLang, text)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "You are a professional translator. Translate the text accurately while maintaining its original style and meaning.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	// 构造请求
	reqBody := openai.ChatCompletionRequest{
		Model:       t.model,
		Messages:    messages,
		Temperature: t.temperature,
		MaxTokens:   t.maxTokens,
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", t.apiURL, bytes.NewBuffer(reqData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.apiKey))

	// 发送请求
	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result openai.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// 检查响应是否包含翻译结果
	if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("no translation result in response")
	}

	return result.Choices[0].Message.Content, nil
}

// CreateChatCompletion 提供原生的ChatCompletion接口
func (t *OpenAITranslator) CreateChatCompletion(ctx context.Context, oaiRequest openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {

	reqData, err := json.Marshal(oaiRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Println(string(reqData))

	req, err := http.NewRequestWithContext(ctx, "POST",
		t.apiURL,
		bytes.NewBuffer(reqData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.apiKey))

	resp, err := t.client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result openai.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetModelInfo 获取模型配置信息
func (t *OpenAITranslator) GetModelInfo() map[string]interface{} {
	return map[string]interface{}{
		"provider": t.provider,
		"model":    t.model,
		"api_url":  t.apiURL,
	}
}

// GetProvider 获取提供商名称
func (t *OpenAITranslator) GetProvider() string {
	return t.provider
}
func (t *OpenAITranslator) GetAPIURL() string {
	return t.apiURL
}

func (t *OpenAITranslator) GetModel() string {
	return t.model
}

// GetMetrics 获取最近一次请求的指标
func (t *OpenAITranslator) GetMetrics() TranslationMetrics {
	return t.lastMetrics
}

// Close 实现清理接口
func (t *OpenAITranslator) Close() error {
	// OpenAI 客户端当前不需要特别的清理操作
	return nil
}

// ValidateConfig 验证配置是否有效
func (t *OpenAITranslator) ValidateConfig() error {
	if t.provider == "" {
		return fmt.Errorf("provider is required")
	}
	if t.model == "" {
		return fmt.Errorf("model is required")
	}
	if t.client == nil {
		return fmt.Errorf("client is not initialized")
	}
	return nil
}

// String 实现 Stringer 接口
func (t *OpenAITranslator) String() string {
	return fmt.Sprintf("%s/%s", t.provider, t.model)
}
