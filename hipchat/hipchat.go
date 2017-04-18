/*Package hipchat is used for implementing a StatusService in BusyBee.

It is not recommended to use this package outside of BusyBee.
To use with BusyBee make sure to import this package.
You can do so by adding the following:
	import _ "github.com/dixonwille/busybee/hipchat"
*/
package hipchat

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"net/url"

	"bytes"

	"github.com/dixonwille/busybee"
)

func init() {
	busybee.RegisterStatusService("hipchat", New)
}

//Hipchat is of type StatusService that can be used to get and update the status of a user.
type Hipchat struct {
	host   string
	token  string
	client *http.Client
}

//New returns a new instance of Hipchat.
func New(conf map[string]string) (busybee.StatusService, error) {
	host, ok := conf["host"]
	if !ok || host == "" {
		return nil, errors.New("host is required to create a Hipchat Status Service")
	}
	token, ok := conf["token"]
	if !ok || token == "" {
		return nil, errors.New("token is required to create a Hipchat Status Service")
	}
	client := new(http.Client)
	hc := &Hipchat{
		host:   host,
		token:  token,
		client: client,
	}
	return hc, hc.valid()
}

//UpdateStatus will update the status of the user to the status specified.
func (h *Hipchat) UpdateStatus(uid string, status busybee.Status) error {
	if status == busybee.StatusUnknown {
		return errors.New("Cannot update status of user to unkown")
	}
	strStatus := convertStatus(status)
	user, err := h.GetUser(uid)
	if err != nil {
		return err
	}
	if strings.ToLower(user.Presence.Show) == strStatus {
		return nil
	}
	user.Presence.Show = strStatus
	user.ID = 0 //Need to omit ID from body in update
	body := new(bytes.Buffer)
	err = user.Encode(body)
	if err != nil {
		return err
	}
	req, err := h.newRequest(http.MethodPut, fmt.Sprintf("/v2/user/%s", uid), body)
	if err != nil {
		return err
	}
	res, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		errModel := NewError()
		err = errModel.Decode(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Could not update the status of the user: Code: %d Message: %s", errModel.Error.Code, errModel.Error.Message)
	}
	return nil
}

//GetUser gets the hipchat user with uid.
func (h *Hipchat) GetUser(uid string) (*User, error) {
	req, err := h.newRequest(http.MethodGet, fmt.Sprintf("/v2/user/%s", uid), nil)
	if err != nil {
		return nil, err
	}
	res, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		errModel := NewError()
		err = errModel.Decode(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Could not get the user: Code: %d Message: %s", errModel.Error.Code, errModel.Error.Message)
	}
	user := NewUser()
	err = user.Decode(res.Body)
	return user, err
}

func (h *Hipchat) valid() error {
	req, err := h.newRequest(http.MethodGet, "/v2/user", nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Set("auth_test", "true")
	req.URL.RawQuery = q.Encode()
	res, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		errModel := NewError()
		err = errModel.Decode(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Hipchat did not accept this client: Code: %d Message: %s", errModel.Error.Code, errModel.Error.Message)
	}
	return nil
}

//newRequest creates a new request with the appropriate headers to make a hipchat call.
func (h *Hipchat) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(h.host)
	if err != nil {
		return nil, err
	}
	u.Path = path
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.token))
	return req, nil
}

func convertStatus(busybeeStatus busybee.Status) string {
	switch busybeeStatus {
	case busybee.StatusAvailable:
		return "chat"
	case busybee.StatusBusy:
		return "dnd"
	default:
		return "chat"
	}
}
