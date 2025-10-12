package models

import "github.com/aws/aws-sdk-go-v2/service/transcribe/types"

type Language string

const (
	LanguageMalayalam Language = "ml-IN"
)

type TranscriptionResponse struct {
	LanguageCode    string `json:"language_code"`
	Text            string `json:"text"`
	TranscriptionId string `json:"transcription_id"`
}

var SupportedLanguagesMap map[string]bool

func InitSupportedLanguages(languages []Language) interface{} {
	SupportedLanguagesMap = make(map[string]bool)
	for index := 0; index < len(languages); index++ {
		SupportedLanguagesMap[string(languages[index])] = true
	}
	return SupportedLanguagesMap
}

func GetSupportedLanguages() []string {
	keys := make([]string, 0, len(SupportedLanguagesMap))

	// Iterate over the map and append each key (k) to the keys slice.
	for k := range SupportedLanguagesMap {
		keys = append(keys, k)
	}
	return keys
}

func IsSupportedLanguage(language string) bool {
	isSupported := SupportedLanguagesMap[language]
	return isSupported
}

func MapLanguageToAWS(lang Language) types.LanguageCode {
	switch lang {
	case LanguageMalayalam:
		return types.LanguageCodeMlIn
	default:
		return types.LanguageCodeEnIn
	}
}
