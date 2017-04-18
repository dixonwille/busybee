package hipchat

import (
	"encoding/json"
	"io"
)

//User is a user in hipchat.
type User struct {
	ID          int           `json:"id,omitempty"`
	Name        string        `json:"name"`
	MentionName string        `json:"mention_name"`
	Presence    *UserPresence `json:"presence"`
	Email       string        `json:"email"`
	Roles       []string      `json:"roles"`
	Title       string        `json:"title"`
	GroupAdmin  bool          `json:"is_group_admin"`
	TimeZone    string        `json:"timezone"`
}

//NewUser creates a new instance of a hipchat user.
func NewUser() *User {
	return &User{
		Presence: new(UserPresence),
	}
}

//Encode takes the model and writes it to writer.
func (hu *User) Encode(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(hu)
}

//Decode reads from reader and updates the model with the data.
func (hu *User) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(hu)
}

//UserPresence is part of the body in a hipchat user.
//This states whether the user is busy or not "show" and a reason "status".
type UserPresence struct {
	Status string `json:"status"`
	Show   string `json:"show"`
}
