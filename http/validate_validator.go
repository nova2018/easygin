package http

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
	"sync"
)

type regexValidation struct {
	mapRegex map[string]*regexp.Regexp
	lock     sync.RWMutex
}

var (
	_regexValidation         = &regexValidation{}
	_registerRegexValidation = &regexValidation{}
)

func (v *regexValidation) GetRegex(key, pattern string) *regexp.Regexp {
	v.lock.RLock()
	val, ok := v.mapRegex[key]
	v.lock.RUnlock()
	if !ok {
		v.lock.Lock()
		if v.mapRegex == nil {
			v.mapRegex = make(map[string]*regexp.Regexp, 1)
		}
		var err error
		val, err = regexp.Compile(pattern)
		if err != nil {
			v.lock.Unlock()
			panic(err)
		}
		v.mapRegex[key] = val
		v.lock.Unlock()
	}
	return val
}

func isRegex(fl validator.FieldLevel) bool {
	reg := _regexValidation.GetRegex(fl.Param(), fl.Param())
	return reg.MatchString(fl.Field().String())
}

func isRegisterRegex(fl validator.FieldLevel) bool {
	tag := fl.GetTag()
	reg := _registerRegexValidation.GetRegex(tag, "")
	return reg.MatchString(fl.Field().String())
}

func init() {
	_ = RegisterValidation(&Validation{
		Tag:        "regex",
		Validation: isRegex,
		Translation: []*Translation{
			NewTranslation("en", NewRegistrationFunc("regex", "{0} must match {1}", false), defaultParamTranslateFunc),
			NewTranslation("zh", NewRegistrationFunc("regex", "{0}必须匹配{1}", false), defaultParamTranslateFunc),
		},
	})
}

func NewRegistrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}
		return
	}
}

func defaultTranslateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		log.Printf("warning: error translating FieldError: %#v", fe)
		return fe.(error).Error()
	}

	return t
}

func defaultParamTranslateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
	if err != nil {
		log.Printf("warning: error translating FieldError: %#v", fe)
		return fe.(error).Error()
	}

	return t
}

func RegisterRegexValidator(tag, pattern, zh, en string) error {
	_ = _registerRegexValidation.GetRegex(tag, pattern)
	return RegisterValidation(&Validation{
		Tag:        tag,
		Validation: isRegisterRegex,
		Translation: []*Translation{
			NewTranslation("en", NewRegistrationFunc(tag, en, false), defaultTranslateFunc),
			NewTranslation("zh", NewRegistrationFunc(tag, zh, false), defaultTranslateFunc),
		},
	})
}
