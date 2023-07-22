## gofiber validation
* add custom tag 
* impl custom error response
* define error message as tag 
  
in below , `mydata` and `username` are custom tag.
```go
type InputFrom struct {
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
	ts := InputFrom{Name: "mehdi shokohi", Data: "df"}

    // Register `username` tag for validation
	RegisterValidation("username", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString("^[a-zA-Z0-9]*[-]?[a-zA-Z0-9]*$", fl.Field().String())
		return match
	})


    // Register `mydata` tag for validation
	RegisterValidation("mydata", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString("^/d*$", fl.Field().String())
		return match
	})

	


	errrs := JsonValidation(ts)
	for _, v := range errrs {
		fmt.Println(v)
	}

}


```


### define custom error response

```go

	SetResponseBuilder(func(field, tag, param, errormessage string) any {
		var el MyValidateError
		el.Message = fmt.Sprintf("%s : %s",field,errormessage)
		el.Field = field
		return el
	})

```

### Go-Fiber Middleware

```go
func main() {
	// http server config
	app := fiber.New(fiber.Config{
		Prefork:     false,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Post("/:id", validation.ValidateBodyAs(models.InputFrom{}),structureForms)
.....
```
### Error Response In Your Api Model

```go
func SetResponseBody(fn func(ctx *fiber.Ctx,errs []interface{})){
    
    ctx.Json(MyResponseProto{Data:nil,Error:errs})
}
``