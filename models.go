package validation

type ValidationError struct {
	Field     string `json:"field"`
	Rule      string `json:"rule"`
	NameSpace string `json:"ns"`
	Param     string `json:"param"`
	Message   string `json:"message"`
}

type Response struct {
	Error interface{}
}
