package MultiLangValidator

import (
	"errors"

	"github.com/ghp3000/MultiLangValidator/translations/tr_en"
	"github.com/ghp3000/MultiLangValidator/translations/tr_zh"
	"github.com/ghp3000/MultiLangValidator/translations/tr_zh_tw"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type TranslationInf interface {
	RegisterDefaultTranslations(v *validator.Validate, trans ut.Translator) (err error)
	LoadFieldsFile(filename string) error
	Field(fe validator.FieldError) string
}

type Locale string

const (
	LocaleZh   Locale = "zh"
	LocaleEn   Locale = "en"
	LocaleZhTw Locale = "zh_tw"
)

type ValidatError struct {
	Field string
	Err   string
}

func (e ValidatError) Error() string {
	return e.Err
}

type Validator struct {
	validate    *validator.Validate
	uni         *ut.UniversalTranslator
	transMap    map[string]ut.Translator
	defaultLang string
}

// NewMultiLangValidator 新建实例,指定默认语言
func NewMultiLangValidator(defaultLang Locale) *Validator {
	m := &Validator{
		validate:    validator.New(),
		transMap:    make(map[string]ut.Translator),
		defaultLang: string(defaultLang),
	}
	_ = m.Register(defaultLang, "")
	return m
}

// Register 注册新的语言,不需要字段名翻译的fieldFilename置为空
func (v *Validator) Register(locale Locale, fieldFilename string) error {
	var t locales.Translator
	var tr TranslationInf
	var translatorName string
	switch locale {
	case LocaleZh:
		t = zh.New()
		tr = tr_zh.New()
		translatorName = string(locale)
	case LocaleEn:
		t = en.New()
		tr = tr_en.New()
		translatorName = string(locale)
	case LocaleZhTw:
		t = zh_Hant_TW.New()
		tr = tr_zh_tw.New()
		translatorName = "zh_Hant_TW"
	default:
		return errors.New("invalid locale")
	}
	if fieldFilename != "" {
		if err := tr.LoadFieldsFile(fieldFilename); err != nil {
			return err
		}
	}
	if v.uni == nil {
		v.uni = ut.New(t, t)
	} else {
		if err := v.uni.AddTranslator(t, true); err != nil {
			return err
		}
	}
	trans, found := v.uni.GetTranslator(translatorName)
	if !found {
		return errors.New("invalid locale")
	}
	if err := tr.RegisterDefaultTranslations(v.validate, trans); err != nil {
		return err
	}
	v.transMap[string(locale)] = trans
	return nil
}

// Validate 校验并返回遇到的第一个错误,指定的语言不存在的时使用默认语言
func (v *Validator) Validate(data interface{}, lang string) *ValidatError {
	trans, ok := v.transMap[lang]
	if !ok {
		trans = v.transMap[v.defaultLang]
	}
	err := v.validate.Struct(data)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				if trans != nil {
					return &ValidatError{
						Field: e.Field(),
						Err:   e.Translate(trans),
					}
				}
				return &ValidatError{
					Field: e.Field(),
					Err:   e.Error(),
				}
			}
		}
		return &ValidatError{
			Field: "",
			Err:   err.Error(),
		}
	}
	return nil
}

// Validates 校验并返回遇到的全部错误,指定的语言不存在的时使用默认语言
func (v *Validator) Validates(data interface{}, lang string) []ValidatError {
	trans, ok := v.transMap[lang]
	if !ok {
		trans = v.transMap[v.defaultLang]
	}
	err := v.validate.Struct(data)
	if err != nil {
		var ret []ValidatError
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				if trans != nil {
					ret = append(ret, ValidatError{
						Field: e.Field(),
						Err:   e.Translate(trans),
					})
				} else {
					ret = append(ret, ValidatError{
						Field: e.Field(),
						Err:   e.Error(),
					})
				}
			}
			return ret
		}
		return append(ret, ValidatError{
			Field: "",
			Err:   err.Error(),
		})
	}
	return nil
}
