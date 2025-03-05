// translator/translator.go
package translator

type Translator interface {
	Translate(text, sourceLang, targetLang string) (string, error)
}
