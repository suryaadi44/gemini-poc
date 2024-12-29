package validator

import (
	"reflect"
	"strings"

	"gemini-poc/utils/custom"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validation "github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
)

type Validator interface {
	ValidateJSON(s interface{}) *custom.ErrorValues
	ValidateQuery(s interface{}) *custom.ErrorValues
	ValidateMeta(s interface{}) *custom.ErrorValues
	ValidateStruct(s interface{}) *custom.ErrorValues
	ValidateVar(field interface{}, filedName string, tag string) *custom.ErrorValues
}

type CustomValidator struct {
	Validate *validation.Validate
	trans    ut.Translator
}

func NewValidator() Validator {
	validate := validation.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_trans.RegisterDefaultTranslations(validate, trans)

	customValidator := &CustomValidator{
		Validate: validate,
		trans:    trans,
	}

	// register custom validation functions
	// _ = customValidator.Validate.RegisterValidation("start_alpha", startWithLetter)
	// _ = customValidator.Validate.RegisterValidation("end_alphanum", endWithLetterOrNumber)
	// _ = customValidator.Validate.RegisterValidation("username", containOnlyLettersNumbersUnderscores)

	// register custom translation
	// customValidator.addTranslation("start_alpha", "must start with a letter")
	// customValidator.addTranslation("end_alphanum", "must end with a letter or number")
	// customValidator.addTranslation("username", "must contain only letters, numbers, and underscores")

	return customValidator
}

func (c *CustomValidator) ValidateVar(field interface{}, filedName string, tag string) *custom.ErrorValues {
	err := c.Validate.Var(field, tag)
	if err == nil {
		return nil
	}

	var errors custom.ErrorValues
	for _, err := range err.(validation.ValidationErrors) {
		errors = append(errors, *custom.NewErrorValue(
			filedName,
			err.Translate(c.trans),
		))
	}

	return &errors
}

func (c *CustomValidator) ValidateStruct(s interface{}) *custom.ErrorValues {
	err := c.Validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors custom.ErrorValues
	for _, err := range err.(validation.ValidationErrors) {
		errors = append(errors, *custom.NewErrorValue(
			err.Field(),
			err.Translate(c.trans),
		))
	}

	return &errors
}

func (c *CustomValidator) ValidateJSON(s interface{}) *custom.ErrorValues {
	c.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 1)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return c.ValidateStruct(s)
}

func (c *CustomValidator) ValidateQuery(s interface{}) *custom.ErrorValues {
	c.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("url"), ",", 1)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return c.ValidateStruct(s)
}

func (c *CustomValidator) ValidateMeta(s interface{}) *custom.ErrorValues {
	c.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("meta"), ",", 1)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return c.ValidateStruct(s)
}

// func (c *CustomValidator) addTranslation(tag string, message string) {
// 	registerFn := func(ut ut.Translator) error {
// 		return ut.Add(tag, message, true)
// 	}

// 	transFn := func(ut ut.Translator, fe validation.FieldError) string {
// 		param := fe.Param()
// 		tag := fe.Tag()

// 		t, err := ut.T(tag, fe.Field(), param)
// 		if err != nil {
// 			return fe.(error).Error()
// 		}
// 		return t
// 	}

// 	_ = c.Validate.RegisterTranslation(tag, c.trans, registerFn, transFn)
// }

// func startWithLetter(fl validation.FieldLevel) bool {
// 	rgx, _ := regexp.Compile(`^[a-zA-Z]`)
// 	return rgx.MatchString(fl.Field().String())
// }

// func endWithLetterOrNumber(fl validation.FieldLevel) bool {
// 	rgx := regexp.MustCompile(`[a-zA-Z0-9]$`)
// 	return rgx.MatchString(fl.Field().String())
// }

// func containOnlyLettersNumbersUnderscores(fl validation.FieldLevel) bool {
// 	rgx := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
// 	return rgx.MatchString(fl.Field().String())
// }
