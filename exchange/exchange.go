/*Package exchange is used as a CalendarService in BusyBee.

It is not recommended to use this package outside of BusyBee.
To use with BusyBee make sure to import this package.
You can do so by adding the following:
	import _ "github.com/dixonwille/busybee/exchange"
*/
package exchange

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bytes"

	"fmt"

	"errors"

	ntlmssp "github.com/Azure/go-ntlmssp"
	"github.com/dixonwille/busybee"
)

func init() {
	busybee.RegisterCalendarService("exchange", New)
}

//Exchange is of type CalendarService and is used to communicate with an Exchange server.
type Exchange struct {
	host     string
	username string
	password string
	client   *http.Client
}

//New creates a new Exchange Calendar service for BusyBee for consumption.
//conf holds information on how to create the new instance.
//Keys that must exist are as follows:
//
//* host - holds the host for the exchange service
//* user - username to sign in to exchange with
//* pass - password for the username to sign in to exchange
func New(conf map[string]string) (busybee.CalendarService, error) {
	host, ok := conf["host"]
	if !ok || host == "" {
		return nil, errors.New("host is required to create an Exchange Calendar Service")
	}
	user, ok := conf["user"]
	if !ok || user == "" {
		return nil, errors.New("user is required to create an Exchange Calendar Service")
	}
	pass, ok := conf["pass"]
	if !ok || pass == "" {
		return nil, errors.New("pass is required to create an Exchange Calendar Service")
	}
	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}
	return &Exchange{
		host:     host,
		username: user,
		password: pass,
		client:   client,
	}, nil
}

//InEvent returns whether the specified uid is in an event or not.
func (e *Exchange) InEvent(uid string) (bool, error) {
	now := time.Now()
	req := NewRequestEnvelope(now, now.AddDate(0, 0, 1), uid)
	reqBuff := new(bytes.Buffer)
	err := req.Encode(reqBuff)
	if err != nil {
		return false, err
	}
	request, err := e.newRequest(http.MethodPost, "/ews/exchange.asmx", reqBuff)
	if err != nil {
		return false, err
	}
	res, err := e.client.Do(request)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Could not see if user was busy: Code: %d Status: %s", res.StatusCode, res.Status)
	}
	response := new(ResponseEnvelope)
	err = response.Decode(res.Body)
	if err != nil {
		return false, err
	}
	if len(response.FreeBusyResponses) < 1 {
		return false, errors.New("Did not get a proper response from the server")
	}
	for _, event := range response.FreeBusyResponses[0].CalendarEvents {
		n := time.Now() //Want to be as close as possable to the current time
		start, err := time.ParseInLocation(DateTimeFormat, event.StartTime, time.Local)
		if err != nil {
			return false, err
		}
		end, err := time.ParseInLocation(DateTimeFormat, event.EndTime, time.Local)
		if err != nil {
			return false, err
		}
		if n.After(start) && n.Before(end) && strings.ToLower(event.BusyType) == "busy" {
			return true, nil
		}
	}
	return false, nil
}

//newRequest creates a new request with the appropriate headers to communicate with Exchange.
func (e *Exchange) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(e.host)
	if err != nil {
		return nil, err
	}
	u.Path = path
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(e.username, e.password)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	return req, err
}
