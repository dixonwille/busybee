/*Package hipchat is used for implementing a StatusService in BusyBee.

It is not recommended to use this package outside of BusyBee.
To use with BusyBee make sure to import this package.
You can do so by adding the following:
	import _ "github.com/dixonwille/busybee/hipchat"
But since you will need the Configuration Struct it may just be easier to use it like normal.
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
	"github.com/dixonwille/busybee/util"
)

func init() {
	hipchatService := new(busybee.StatusService)
	hipchatService.Create = New
	hipchatService.CreateConfig = NewConf
	busybee.RegisterStatusService("hipchat", hipchatService)
}

//Hipchat is of type StatusService that can be used to get and update the status of a user.
type Hipchat struct {
	host    string
	token   string
	client  *http.Client
	busybee *busybee.BusyBee
}

//Conf holds all the needed information to create a new hipchat service.
type Conf struct {
	Host  string `quest:"What is the hipchat host?"`
	Token string `quest:"What is your hipchat token?,encrypt"`
}

//NewConf creates a new Hipchat configuration.
func NewConf() interface{} {
	return new(Conf)
}

//New returns a new instance of Hipchat.
func New(conf interface{}, bb *busybee.BusyBee) (busybee.UpdateStatuser, error) {
	hcConf, ok := conf.(*Conf)
	if !ok {
		return nil, errors.New("Must use a configuration struct from the hipchat service")
	}
	client := new(http.Client)
	hc := &Hipchat{
		host:    util.CleanHost(hcConf.Host),
		token:   hcConf.Token,
		client:  client,
		busybee: bb,
	}
	return hc, hc.valid()
}

//UpdateStatus will update the status of the user to the status specified.
//Will not update if the user is already on specified status.
//Updates the show if in meeting to DnD always.
//Updates the show if DnD to Available but never Away to Available.
//Updates the Status message to "In a meeting (BusyBee)" if changing to DnD.
//Updates the Status message to "I'm free (BusyBee)" if changing to Available.
func (h *Hipchat) UpdateStatus(uid string, status busybee.Status) error {
	if status == busybee.StatusUnknown {
		return errors.New("Cannot update status of user to unkown")
	}
	user, err := h.GetUser(uid)
	if err != nil {
		return err
	}
	strStatus := convertStatus(status)
	//User is offline don't show them as online
	if user.Presence == nil {
		return nil
	}
	curStatus := strings.ToLower(user.Presence.Show)
	switch {
	case curStatus == strStatus:
		fallthrough
	case (curStatus == "xa" || curStatus == "away") && status == busybee.StatusAvailable:
		return nil
	case status == busybee.StatusBusy:
		user.Presence.Status = "In a meeting (BusyBee)"
	case status == busybee.StatusAvailable:
		user.Presence.Status = "I'm free (BusyBee)"
	}
	user.ID = 0 //Need to omit ID from body in update
	user.Presence.Show = strStatus
	body := new(bytes.Buffer)
	err = user.Encode(body)
	if err != nil {
		return err
	}
	req, err := h.newRequest(http.MethodPut, fmt.Sprintf("/v2/user/%s", util.CleanMention(uid)), body)
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
//uid is the users @Mention name.
func (h *Hipchat) GetUser(uid string) (*User, error) {
	req, err := h.newRequest(http.MethodGet, fmt.Sprintf("/v2/user/%s", util.CleanMention(uid)), nil)
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
	token, err := h.busybee.Decrypt(h.token, "")
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
