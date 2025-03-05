// translator/model_manager.go
package translator

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"transbridge/config"
)

// 模型标识结构
type ModelIdentifier struct {
	Provider string // 服务提供商
	Model    string // 模型名称
	APIURL   string // 模型服务地址
}

func (m ModelIdentifier) String() string {
	return fmt.Sprintf("%s/%s", m.Provider, m.Model)
}

// ModelManager 管理多个服务提供商和其模型
type ModelManager struct {
	translators  map[ModelIdentifier]*OpenAITranslator
	modelWeights map[ModelIdentifier]int
	defaultModel ModelIdentifier
	mu           sync.RWMutex
}

func NewModelManager(providers []config.ProviderConfig) (*ModelManager, error) {
	if len(providers) == 0 {
		return nil, errors.New("no providers configured")
	}

	mm := &ModelManager{
		translators:  make(map[ModelIdentifier]*OpenAITranslator),
		modelWeights: make(map[ModelIdentifier]int),
	}

	var defaultFound bool
	for _, provider := range providers {
		// 获取提供商的默认超时时间
		defaultTimeout := provider.Timeout

		for _, modelCfg := range provider.Models {
			// 确定模型的超时时间
			timeout := defaultTimeout
			if modelCfg.Timeout != nil {
				timeout = *modelCfg.Timeout
			}

			identifier := ModelIdentifier{
				Provider: provider.Provider,
				Model:    modelCfg.Name,
			}

			// 创建翻译器实例
			translator := NewOpenAITranslator(
				provider.Provider,
				provider.APIURL,
				provider.APIKey,
				modelCfg.Name,
				timeout,
				modelCfg.MaxTokens,
				modelCfg.Temperature,
			)

			//			log.Println("Translator created:", *translator)

			mm.translators[identifier] = translator
			mm.modelWeights[identifier] = modelCfg.Weight

			// 如果是默认提供商的第一个模型，设为默认模型
			if provider.IsDefault && !defaultFound {
				mm.defaultModel = identifier
				defaultFound = true
			}
		}
	}

	if !defaultFound {
		// 如果没有设置默认模型，使用第一个可用的模型
		for identifier := range mm.translators {
			mm.defaultModel = identifier
			break
		}
	}

	return mm, nil
}

// GetModel 获取指定提供商和模型的翻译器
func (mm *ModelManager) GetModel(provider, model string) (*OpenAITranslator, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	identifier := ModelIdentifier{Provider: provider, Model: model}
	if translator, exists := mm.translators[identifier]; exists {
		return translator, nil
	}
	return nil, fmt.Errorf("model %s not found for provider %s", model, provider)
}

// GetDefaultModel 获取默认模型
func (mm *ModelManager) GetDefaultModel() *OpenAITranslator {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	return mm.translators[mm.defaultModel]
}

func (mm *ModelManager) GetRandomModel() *OpenAITranslator {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	var totalWeight int
	for _, weight := range mm.modelWeights {
		totalWeight += weight
	}

	r := rand.Intn(totalWeight)
	for identifier, weight := range mm.modelWeights {
		r -= weight
		if r <= 0 {
			return mm.translators[identifier]
		}
	}

	return mm.translators[mm.defaultModel]
}

// ListModels 列出所有可用的模型
func (mm *ModelManager) ListModels() []ModelIdentifier {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	models := make([]ModelIdentifier, 0, len(mm.translators))
	for identifier := range mm.translators {
		models = append(models, identifier)
	}
	return models
}

// GetModelsByProvider 获取指定提供商的所有模型
func (mm *ModelManager) GetModelsByProvider(provider string) []string {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	var models []string
	for identifier := range mm.translators {
		if identifier.Provider == provider {
			models = append(models, identifier.Model)
		}
	}
	return models
}
