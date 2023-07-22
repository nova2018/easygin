package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
	"sync"
)

type contextShouldBind interface {
	ShouldBind(obj interface{}) error
}
type contextShouldBindWith interface {
	ShouldBindWith(obj interface{}, b binding.Binding) error
}

type validateConfig struct {
	lock            sync.Mutex
	listTranslation []*validateTranslationConfig
	mapValidation   map[string]*Validation
	listValidate    []*defaultValidate
}

type defaultValidate struct {
	once             sync.Once
	validate         *validator.Validate
	translator       *ut.UniversalTranslator
	translateTagName string
	listTranslator   []locales.Translator
	mapTranslation   map[string]RegisterValidatorTranslations
	local            string
	isInit           bool
}

type validateTranslationConfig struct {
	local       string
	translation RegisterValidatorTranslations
	override    bool
}

type RegisterValidatorTranslations func(v *validator.Validate, trans ut.Translator) error

type Validation struct {
	Tag                      string
	Validation               validator.Func
	CallValidationEvenIfNull bool
	Translation              []*Translation
}

type Translation struct {
	Locale              string
	RegisterTranslation validator.RegisterTranslationsFunc
	Translation         validator.TranslationFunc
	IsDone              bool
}

func NewTranslation(locale string, reg validator.RegisterTranslationsFunc, trans ...validator.TranslationFunc) *Translation {
	var tr validator.TranslationFunc
	if len(trans) > 0 {
		tr = trans[0]
	}
	return &Translation{
		Locale:              locale,
		RegisterTranslation: reg,
		Translation:         tr,
	}
}

var (
	_validate         = NewValidate(GetBindingValidator())
	_validateCfg      = newValidateConfig()
	_translateTagName = "zh"
)

func newValidateConfig() *validateConfig {
	v := &validateConfig{
		listTranslation: make([]*validateTranslationConfig, 0),
		mapValidation:   make(map[string]*Validation, 0),
		listValidate:    make([]*defaultValidate, 0, 1),
	}
	return v
}

func NewValidate(vd *validator.Validate) *defaultValidate {
	v := new(defaultValidate)
	v.validate = vd
	v.listTranslator = []locales.Translator{
		en.New(), //english
		zh.New(), //chinese
	}
	uni := ut.New(v.listTranslator[0], v.listTranslator...)
	v.mapTranslation = make(map[string]RegisterValidatorTranslations, 0)
	v.translator = uni
	v.local = "zh"
	v.translateTagName = _translateTagName
	vd.RegisterTagNameFunc(validatorTagNameTranslate(v))
	return v
}

func GetBindingValidator() *validator.Validate {
	if x, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return x
	}
	panic("get validate failure")
}

func AddTranslator(trans locales.Translator, translation RegisterValidatorTranslations, override bool) error {
	return _validate.AddTranslator(trans, translation, override)
}

func SetDefaultLocal(local string) error {
	_, err := _validate.GetTranslator(local)
	if err != nil {
		return err
	}
	_validate.local = local
	return nil
}

func (c *validateConfig) RegisterValidation(validation *Validation) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.mapValidation[validation.Tag] = validation

	for _, vd := range c.listValidate {
		err := vd.RegisterValidation(validation)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *validateConfig) AddTranslator(local string, translation RegisterValidatorTranslations, override bool) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, v := range c.listValidate {
		trans, err := v.GetTranslator(local)
		if err != nil {
			return err
		}
		err = v.AddTranslator(trans, translation, override)
		if err != nil {
			return err
		}
	}
	c.listTranslation = append(c.listTranslation, &validateTranslationConfig{
		local:       local,
		translation: translation,
		override:    override,
	})
	return nil
}

func (c *validateConfig) FillValidate(vd *defaultValidate) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, tr := range c.listTranslation {
		trans, err := vd.GetTranslator(tr.local)
		if err != nil {
			return err
		}
		err = vd.AddTranslator(trans, tr.translation, tr.override)
		if err != nil {
			return err
		}
	}
	for _, t := range c.mapValidation {
		err := vd.RegisterValidation(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *validateConfig) BindValidate(vd *defaultValidate) error {
	err := c.FillValidate(vd)
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.listValidate = append(c.listValidate, vd)
	return nil
}

func (v *defaultValidate) AddTranslator(trans locales.Translator, translation RegisterValidatorTranslations, override bool) error {
	err := v.translator.AddTranslator(trans, override)
	if err == nil {
		v.mapTranslation[trans.Locale()] = translation
		v.listTranslator = append(v.listTranslator, trans)
		if v.isInit {
			err = v.RegisterTranslate(trans.Locale(), translation)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (v *defaultValidate) lazyInit() {
	if v.isInit {
		return
	}
	v.once.Do(func() {
		for _, trans := range v.listTranslator {
			err := v.RegisterTranslate(trans.Locale(), nil)
			if err != nil {
				fmt.Printf("validate lazy init failure! err=[%v]\n", err)
			}
		}
		v.isInit = true
	})
}

func (v *defaultValidate) RegisterTranslate(local string, translation RegisterValidatorTranslations) error {
	trans, err := v.GetTranslator(local)
	if err != nil {
		return err
	}

	if translation == nil {
		if x, ok := v.mapTranslation[local]; ok {
			translation = x
		} else {
			switch local {
			case "en":
				translation = enTranslations.RegisterDefaultTranslations
			case "zh":
				translation = chTranslations.RegisterDefaultTranslations
			default:
				translation = enTranslations.RegisterDefaultTranslations
			}
		}
	}

	err = translation(v.validate, trans)

	return err
}

func (v *defaultValidate) GetTranslator(local string) (ut.Translator, error) {
	var o bool
	trans, o := v.translator.GetTranslator(local)
	if !o {
		return nil, fmt.Errorf("uni.GetTranslator(%s) failed", local)
	}
	return trans, nil
}

func (v *defaultValidate) GetLocal(local ...string) string {
	currentLocal := v.local
	if len(local) > 0 {
		currentLocal = local[0]
	}
	return currentLocal
}

func (v *defaultValidate) Validate(fn func() error, local ...string) (bool, []error) {
	v.lazyInit()
	currentLocal := v.GetLocal(local...)

	if err := fn(); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return false, []error{err}
		}
		trans, err := v.GetTranslator(currentLocal)
		if err != nil {
			return false, []error{err}
		}
		listTranslation := errs.Translate(trans)
		listErr := make([]error, 0, len(listTranslation))
		for _, msg := range listTranslation {
			listErr = append(listErr, errors.New(msg))
		}
		return false, listErr
	}
	return true, nil
}

func (v *defaultValidate) RegisterValidation(validation *Validation) error {
	vd := v.validate
	err := vd.RegisterValidation(validation.Tag, validation.Validation, validation.CallValidationEvenIfNull)
	if err != nil {
		return err
	}

	if validation.Translation != nil && len(validation.Translation) > 0 {
		for _, tr := range validation.Translation {
			trans, err := v.GetTranslator(tr.Locale)
			if err != nil {
				return err
			}
			if tr.Translation == nil {
				tr.Translation = defaultTranslateFunc
			}
			err = vd.RegisterTranslation(validation.Tag, trans, tr.RegisterTranslation, tr.Translation)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validatorTagNameTranslate(v *defaultValidate) func(field reflect.StructField) string {
	return func(field reflect.StructField) string {
		var listTag []string
		if v.local == "zh" {
			listTag = []string{v.translateTagName, "json"}
		} else {
			listTag = []string{"json"}
		}
		for _, s := range listTag {
			if s == "" {
				continue
			}
			tagVal := field.Tag.Get(s)
			if tagVal == "" {
				continue
			}
			parts := strings.SplitN(tagVal, ",", 2)
			if parts[0] != "-" && parts[0] != "" {
				return parts[0]
			}
		}
		return ""
	}
}

func BindValidate(validate *defaultValidate) error {
	return _validateCfg.BindValidate(validate)
}

func RegisterValidation(validation *Validation) error {
	return _validateCfg.RegisterValidation(validation)
}

func FillValidate(vd *defaultValidate) error {
	return _validateCfg.FillValidate(vd)
}

func ValidateWithCustom(fn func() error, local ...string) (bool, []error) {
	return _validate.Validate(fn, local...)
}

func SetValidateRemarkTagName(tagName string) {
	_validate.translateTagName = tagName
}

func Validate(ctx contextShouldBind, obj interface{}, local ...string) (bool, []error) {
	return ValidateWithCustom(func() error {
		return ctx.ShouldBind(obj)
	}, local...)
}

func ValidateWithShould(fn func(obj interface{}) error, obj interface{}, local ...string) (bool, []error) {
	return ValidateWithCustom(func() error {
		return fn(obj)
	}, local...)
}

func ValidateWithBinding(ctx contextShouldBindWith, obj interface{}, binding binding.Binding, local ...string) (bool, []error) {
	return ValidateWithCustom(func() error {
		return ctx.ShouldBindWith(obj, binding)
	}, local...)
}

func init() {
	_ = BindValidate(_validate)
}
