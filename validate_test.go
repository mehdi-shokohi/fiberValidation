package validation

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/go-playground/validator/v10"
)

type FactoryModelInputFrom struct {
	Title string `json:"title" validate:"required,min=3,max=32"`
	Name  string `json:"name" validate:"username,required,min=3,max=25" errmsg:"invalid username"`
	Data  string `json:"data" validate:"required,mydata" errmsg:"invalid data"`
}
type MyValidateError struct {
	Field string
	Tag  string
	Message string
}
func TestVal(t *testing.T) {
	ts := FactoryModelInputFrom{Name: "mehdi shokohi", Data: "df"}
	RegisterValidation("username", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString("^[a-zA-Z0-9]*[-]?[a-zA-Z0-9]*$", fl.Field().String())
		return match
	})
	RegisterValidation("mydata", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString("^/d*$", fl.Field().String())
		return match
	})

	
	SetResponseBuilder(func(field, tag, param, errormessage string) any {
		var el MyValidateError
		el.Message = fmt.Sprintf("%s : %s",field,errormessage)
		el.Field = field
		return el
	})

	errrs := JsonValidation(ts)
	for _, v := range errrs {
		fmt.Println(v)
	}

}
