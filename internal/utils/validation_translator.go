package utils

import (
	"errors"
	"fmt"
	"snack-store-api/internal/messages"
	"sort"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var translatorCache sync.Map

func InitTranslator(v *validator.Validate) ut.Translator {
	if v == nil {
		return nil
	}

	if cached, ok := translatorCache.Load(v); ok {
		return cached.(ut.Translator)
	}

	uni := ut.New(en.New())
	enTrans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(v, enTrans)
	translatorCache.Store(v, enTrans)

	return enTrans
}

func TranslateValidationError(v *validator.Validate, err error) string {
	enTrans := InitTranslator(v)
	if enTrans == nil {
		return messages.FailedValidationOccurred
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return messages.FailedValidationOccurred
	}

	return joinMessages(validationErrors.Translate(enTrans))
}

func joinMessages(msgs map[string]string) string {
	if len(msgs) == 0 {
		return ""
	}

	keys := make([]string, 0, len(msgs))
	for k := range msgs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, msgs[k])
	}

	return strings.Join(parts, ", ")
}

func WrapMessageAsError(msg string, err ...error) error {
	if len(err) > 0 && err[0] != nil {
		return fmt.Errorf("%s: %w", msg, err[0])
	}

	return errors.New(msg)
}
