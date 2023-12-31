package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate
var validatorStack map[string]func(validator.FieldLevel) bool
var response func(field, tag, param, errormessage string) any
var errorBody func(ctx *fiber.Ctx, errs []interface{})

func JsonValidation[T any](o T) []interface{} {
	if validate == nil {
		load()
	}
	var errors []interface{}

	err := validate.Struct(o)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			rsf := reflect.TypeOf(&o).Elem()
			field, _ := rsf.FieldByName(err.Field())
			resp := response(err.Field(), err.Tag(), err.Param(), field.Tag.Get("errmsg"))
			errors = append(errors, resp)
		}

	}
	return errors

}

func ValidateBodyAs[T any](body T) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
		}

		errs := JsonValidation(body)
		if len(errs) > 0 {
			if errorBody == nil {

				return c.Status(fiber.StatusBadRequest).JSON(Response{Error: errs})
			}
			errorBody(c, errs)
			return nil
		}
		return c.Next()
	}
}

func load() {
	validate = validator.New()
	if validatorStack == nil {
		validatorStack = make(map[string]func(validator.FieldLevel) bool)
	}
	for k, v := range validatorStack {

		validate.RegisterValidation(k, v)
	}
	if response == nil {
		SetErrorBuilder(func(field, tag, param, errormessage string) any {
			var el ValidationError
			el.Message = errormessage

			el.Field = field
			el.Rule = tag
			el.Param = param
			return el
		})
	}
}
func RegisterValidation(tag string, fn func(validator.FieldLevel) bool) {
	if validatorStack == nil {
		validatorStack = make(map[string]func(validator.FieldLevel) bool)
	}
	validatorStack[tag] = fn
}
func SetErrorBuilder(f func(field, tag, param, errormessage string) any) {
	response = f
}

func SetResponseBody(fn func(ctx *fiber.Ctx, errs []interface{})) {
	errorBody = fn
}
