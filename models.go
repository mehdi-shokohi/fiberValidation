package validation


type ValidationError struct {
	Field string
	Rule  string
	Param string
	Message string
}
