package validation

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/go-playground/validator/v10"
)

type InputFrom struct {
	Title string `json:"title" validate:"required,min=3,max=32"`
	Name  string `json:"name" validate:"required,min=3,max=25" `
	Data  Spec `json:"data" validate:"required"`
}
type Spec struct{
	Attr string  `json:"attr" validate:"inernalId"`
	Param []string `json:"param" validate:"required"`
}
type MyValidateError struct {
	Field string
	Tag  string
	Message string
}
func TestVal(t *testing.T) {
	ts := InputFrom{Name: "hi", Data:Spec{Attr: "e1",Param: []string{"n1","n2"}}}
	fv:=NewFiberValidation(WithResponseCast(func(errs []ValidationError) any {
		errList:=[]map[string]any{}
		for _,el:=range errs{
			errList=append(errList, map[string]any{
				"field":el.Field,
				"ns":el.NameSpace,
				"message":el.Message,


			})
		}
		return map[string]any{"data":nil,"errors":errList}
	}))
	// fv:=GetValidator()
	fv.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		if valuer, ok := fl.Field().Interface().(string); ok {
			//checking username existence in database
			if valuer=="mate jason"{
				return false
			}
		}
		match, _ := regexp.MatchString("^[a-zA-Z0-9]*[-]?[a-zA-Z0-9]*$", fl.Field().String())
		return match
	},"field {0} must be compatible with pattern {1}","username","1-no space")
	fv.RegisterValidation("inernalId", func(fl validator.FieldLevel) bool {
		match, err := regexp.MatchString("^/d*$", fl.Field().String())
		if err!=nil{
			return false
		}
		return match
	},"field {0} must be compatible with pattern {1}","inernalId","1-no space")


	
	errrs := fv.JsonValidation(ts)
	for _, v := range errrs {
		b,_:=json.Marshal(v)
		fmt.Println(string(b))
	}

}
