package settings

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"os"
)

// ReloadLocalBundle ReloadThird 加载本地国际化文件
func ReloadLocalBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 加载翻译文件
	//bundle.MustLoadMessageFile("locales/active.en.toml")
	//bundle.MustLoadMessageFile("locales/active.zh.toml")
	loadMessageFiles(bundle)
	return bundle
}

func loadMessageFiles(b *i18n.Bundle) {
	files, err := os.ReadDir("locales")
	if err != nil {
		log.Fatalf("failed to read locales directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		b.MustLoadMessageFile("locales/" + file.Name())
	}
}

func LocalizeMsg(locale *i18n.Localizer, messageID string) string {
	return locale.MustLocalize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

func LocalizeMsgCount(locale *i18n.Localizer, messageID string, count string) string {
	return locale.MustLocalize(&i18n.LocalizeConfig{
		MessageID: messageID,
		TemplateData: map[string]interface{}{
			"Count": count,
		},
	})
}

func LocalizeMsgTemplateData(locale *i18n.Localizer, messageID string, templateData map[string]interface{}) string {
	return locale.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
}
