package services

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dixonwille/busybee"
)

//Hipchat is of type StatusService that can be used to get and update the status of a user.
type Hipchat struct {
	host   string
	token  string
	client *http.Client
}

//NewHipchat returns a new instance of Hipchat.
func NewHipchat(host string, accessToken string) *Hipchat {
	client := new(http.Client)
	return &Hipchat{
		host:   host,
		token:  accessToken,
		client: client,
	}
}

//UpdateStatus will update the status of the user to the status specified.
func (h *Hipchat) UpdateStatus(uid string, status busybee.Status) error {
	return nil
}

//CurrentStatus will return the current status of the user.
func (h *Hipchat) CurrentStatus(uid string) (busybee.Status, error) {
	return busybee.Available, nil
}

//NewRequest creates a new request with the appropriate headers to make a hipchat call.
func (h *Hipchat) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.token))
	return req, nil
}
