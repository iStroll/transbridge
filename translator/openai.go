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
	"transbridge/internal/utils"
)

// TranslationMetrics 翻译指标
type TranslationMetrics struct {
	InputTokens  int     `json:"input_tokens"`  // 输入token数
	OutputTokens int     `json:"output_tokens"` // 输出token数
	TotalTokens  int     `json:"total_tokens"`  // 总token数
	ModelLatency float64 `json:"model_latency"` // 模型处理延迟（毫秒）
}

// OpenAITranslator 实现 OpenAI 的翻译器
type OpenAITranslator struct {
	Provider    string
	ApiURL      string
	ApiKey      string
	Model       string
	Timeout     int
	MaxTokens   int
	Temperature float32
	Client      *http.Client
	LastMetrics TranslationMetrics
}

// 确保 OpenAITranslator 实现了 Translator 接口
var _ Translator = (*OpenAITranslator)(nil)

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
		Provider:    provider,
		ApiURL:      apiURL,
		ApiKey:      apiKey,
		Model:       model,
		Timeout:     timeout,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		Client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// Translate 实现翻译功能
func (t *OpenAITranslator) Translate(promptTemplate, text, sourceLang, targetLang string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Second)
	defer cancel()

	return t.TranslateWithContext(ctx, promptTemplate, text, sourceLang, targetLang)
}

// TranslateWithContext 支持上下文的翻译方法
func (t *OpenAITranslator) TranslateWithContext(ctx context.Context, promptTemplate, text, sourceLang, targetLang string) (string, error) {
	log.Println(t.ApiURL, t.Model)

	slang, _ := utils.GetLanguageName(sourceLang)
	tlang, _ := utils.GetLanguageName(targetLang)

	prompt := utils.ApplyPromptTemplate(promptTemplate, text, slang, tlang)

	log.Println(prompt)

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
		Model:       t.Model,
		Messages:    messages,
		Temperature: t.Temperature,
		MaxTokens:   t.MaxTokens,
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", t.ApiURL, bytes.NewBuffer(reqData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.ApiKey))

	// 发送请求
	resp, err := t.Client.Do(req)
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

// GetProvider 获取提供商名称
func (t *OpenAITranslator) GetProvider() string {
	return t.Provider
}

// GetAPIURL 获取 API URL
func (t *OpenAITranslator) GetAPIURL() string {
	return t.ApiURL
}

// GetModel 获取模型名称
func (t *OpenAITranslator) GetModel() string {
	return t.Model
}

// GetMetrics 获取最近一次请求的指标
func (t *OpenAITranslator) GetMetrics() TranslationMetrics {
	return t.LastMetrics
}

// Close 实现清理接口
func (t *OpenAITranslator) Close() error {
	// OpenAI 客户端当前不需要特别的清理操作
	return nil
}

// ValidateConfig 验证配置是否有效
func (t *OpenAITranslator) ValidateConfig() error {
	if t.Provider == "" {
		return fmt.Errorf("provider is required")
	}
	if t.Model == "" {
		return fmt.Errorf("model is required")
	}
	if t.Client == nil {
		return fmt.Errorf("client is not initialized")
	}
	return nil
}

// String 实现 Stringer 接口
func (t *OpenAITranslator) String() string {
	return fmt.Sprintf("%s/%s", t.Provider, t.Model)
}

// OpenAIChatCompletion 提供 OpenAI 聊天完成功能
type OpenAIChatCompletion struct {
	*OpenAITranslator
}

// NewOpenAIChatCompletion 创建新的 OpenAI 聊天完成实例
func NewOpenAIChatCompletion(translator *OpenAITranslator) *OpenAIChatCompletion {
	return &OpenAIChatCompletion{
		OpenAITranslator: translator,
	}
}

// CreateChatCompletion 提供原生的ChatCompletion接口
func (t *OpenAIChatCompletion) CreateChatCompletion(ctx context.Context, oaiRequest openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	reqData, err := json.Marshal(oaiRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Println(string(reqData))

	req, err := http.NewRequestWithContext(ctx, "POST",
		t.ApiURL,
		bytes.NewBuffer(reqData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.ApiKey))

	resp, err := t.Client.Do(req)
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
