package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaTranslator 实现 Ollama 的翻译器
type OllamaTranslator struct {
	apiURL     string
	model      string
	timeout    time.Duration
	httpClient *http.Client
}

// OllamaRequest 定义 Ollama API 请求结构
type OllamaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// OllamaResponse 定义 Ollama API 响应结构
type OllamaResponse struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`
}

// Message 定义消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewOllamaTranslator 创建新的 Ollama 翻译器实例
func NewOllamaTranslator(apiURL, model string, timeout int) *OllamaTranslator {
	return &OllamaTranslator{
		apiURL:     apiURL,
		model:      model,
		timeout:    time.Duration(timeout) * time.Second,
		httpClient: &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

// Translate 实现翻译接口
func (t *OllamaTranslator) Translate(text, sourceLang, targetLang string) (string, error) {
	prompt := fmt.Sprintf("Translate the following text from %s to %s:\n\n%s", sourceLang, targetLang, text)

	reqBody := OllamaRequest{
		Model: t.model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", t.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Message.Content, nil
}

// GetAPIURL 返回 API URL
func (t *OllamaTranslator) GetAPIURL() string {
	return t.apiURL
}

// GetModel 返回模型名称
func (t *OllamaTranslator) GetModel() string {
	return t.model
}

// GetProvider 返回提供商名称
func (t *OllamaTranslator) GetProvider() string {
	return "ollama"
}

func (t *OllamaTranslator) Close() error {
	return nil
}
