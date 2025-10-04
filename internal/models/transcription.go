package models

import "github.com/aws/aws-sdk-go-v2/service/transcribe/types"

type Language types.Type

const (
	LanguageMalayalam Language = "ml-IN"
)

func MapLanguageToAWS(lang Language) types.LanguageCode {
	switch lang {
	case LanguageMalayalam:
		return types.LanguageCodeMlIn
	default:
		return types.LanguageCodeEnIn
	}
}
