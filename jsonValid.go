package validation

import (
	enLocale "github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var vfglob *FiberValidator

type FiberValidator struct {
	validate *validator.Validate
	response func(errs []ValidationError) any
	trans ut.Translator
}
func WithTranslator(ut ut.Translator)func(*FiberValidator){
	return func(fv *FiberValidator) {
		fv.trans = ut
	}
}
func WithResponseCast(castFn func(errs []ValidationError)any)func(*FiberValidator){
	return func(fv *FiberValidator) {
		fv.response = castFn
	}
}
func (fv *FiberValidator) SetTranslator(ut ut.Translator) {
	fv.trans = ut
}

func NewFiberValidation(opts ...func(*FiberValidator)) *FiberValidator {
	fv := &FiberValidator{
		validate: validator.New(),
	}
	for _,o:=range opts{
		o(fv)
	}
	if fv.response == nil {
		fv.response = func(errs []ValidationError) any {
			return map[string]any{"data": nil,"error":errs}
		}
	}
	if fv.trans==nil{
	en := enLocale.New()
	uni := ut.New(en, en)

	fv.trans, _ = uni.GetTranslator("en")
	}
	vfglob = fv
	return fv
}
func (fv *FiberValidator) JsonValidation(o any) (errors []ValidationError) {

	err := fv.validate.Struct(o)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field: err.Field(),
				Rule: err.Tag(),
				Message: err.Translate(fv.trans),
				Param: err.Param(),
				NameSpace: err.Namespace(),
			})
		}

	}
	return 

}
func GetValidator() *FiberValidator {
	if vfglob == nil {
		vfglob = NewFiberValidation()
	}
	return vfglob
}
func ValidateBodyAs[T any](body T) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(Response{Error: err.Error()})
		}
		v := GetValidator()
		errs := v.JsonValidation(body)
		if len(errs) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(v.response(errs))

		}
		return c.Next()
	}
}

func (fv *FiberValidator) RegisterValidation(tag string, fn func(validator.FieldLevel) bool, errMessagePattern string,extractErrField func(validator.FieldError)[]string) {

	fv.validate.RegisterValidation(tag, fn)
	fv.validate.RegisterTranslation(tag, fv.trans, func(ut ut.Translator) error {
		return ut.Add(
			tag,
			errMessagePattern,
			true,
		)

	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, extractErrField(fe)...)

		return t
	})

}
