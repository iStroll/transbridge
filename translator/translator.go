// translator/translator.go
package translator

// Translator 定义翻译器接口
type Translator interface {
	Translate(text, sourceLang, targetLang string) (string, error)
	GetAPIURL() string
	GetModel() string
	GetProvider() string
	Close() error
}
