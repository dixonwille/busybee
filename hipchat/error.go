package hipchat

import (
	"encoding/json"
	"io"
)

//Error is a response from hipchat that returns an error.
type Error struct {
	Error *ErrorBody `json:"error"`
}

//NewError returns a new instance of HipchatError.
func NewError() *Error {
	return &Error{
		Error: new(ErrorBody),
	}
}

//Decode reads from reader and updates the model with the data.
func (he *Error) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(he)
}

//ErrorBody is the main body of the error response.
type ErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}
