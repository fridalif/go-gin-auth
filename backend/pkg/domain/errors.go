package domain

type MyError struct {
	ErrorBase error
	Module    string
}

func NewError(err error, module string) *MyError {
	return &MyError{
		ErrorBase: err,
		Module:    module,
	}
}
