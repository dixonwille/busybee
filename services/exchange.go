package services

import (
	"io"
	"net/http"

	ntlmssp "github.com/Azure/go-ntlmssp"
)

//Exchange is of type CalendarService and is used to communicate with an Exchange server.
type Exchange struct {
	host     string
	username string
	password string
	client   *http.Client
}

//NewExchange creates a new instance to connect to the Exchange server with.
func NewExchange(host, username, password string) *Exchange {
	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}
	return &Exchange{
		host:     host,
		username: username,
		password: password,
		client:   client,
	}
}

//InEvent returns whether the specified uid is in an event or not.
func (e *Exchange) InEvent(uid string) (bool, error) {
	return false, nil
}

//NewRequest creates a new request with the appropriate headers to communicate with Exchange.
func (e *Exchange) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(e.username, e.password)
	return req, err
}
