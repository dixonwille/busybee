package models

import (
	"encoding/json"
	"io"
)

//HipchatUser is a user in hipchat.
type HipchatUser struct {
	ID          int                  `json:"id,omitempty"`
	Name        string               `json:"name"`
	MentionName string               `json:"mention_name"`
	Presence    *HipchatUserPresence `json:"presence"`
	Email       string               `json:"email"`
	Roles       []string             `json:"roles"`
	Title       string               `json:"title"`
	GroupAdmin  bool                 `json:"is_group_admin"`
	TimeZone    string               `json:"timezone"`
}

//NewHipchatUser creates a new instance of a hipchat user.
func NewHipchatUser() *HipchatUser {
	return &HipchatUser{
		Presence: new(HipchatUserPresence),
	}
}

//Encode takes the model and writes it to writer.
func (hu *HipchatUser) Encode(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(hu)
}

//Decode reads from reader and updates the model with the data.
func (hu *HipchatUser) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(hu)
}

//HipchatUserPresence is part of the body in a hipchat user.
//This states whether the user is busy or not "show" and a reason "status".
type HipchatUserPresence struct {
	Status string `json:"status"`
	Show   string `json:"show"`
}

//HipchatError is a response from hipchat that returns an error.
type HipchatError struct {
	Error *HipchatErrorBody `json:"error"`
}

//NewHipchatError returns a new instance of HipchatError.
func NewHipchatError() *HipchatError {
	return &HipchatError{
		Error: new(HipchatErrorBody),
	}
}

//Decode reads from reader and updates the model with the data.
func (he *HipchatError) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(he)
}

//HipchatErrorBody is the main body of the error response.
type HipchatErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}
