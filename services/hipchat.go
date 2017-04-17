package services

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"net/url"

	"encoding/json"

	"bytes"

	"github.com/dixonwille/busybee"
	"github.com/dixonwille/busybee/models"
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
	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(user)
	if err != nil {
		return err
	}
	req, err := h.NewRequest(http.MethodPut, fmt.Sprintf("/v2/user/%s", uid), body)
	if err != nil {
		return err
	}
	res, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		return errors.New("Could not update the status of the user")
	}
	return nil
}

//GetUser gets the hipchat user with uid.
func (h *Hipchat) GetUser(uid string) (*models.HipchatUser, error) {
	req, err := h.NewRequest(http.MethodGet, fmt.Sprintf("/v2/user/%s", uid), nil)
	if err != nil {
		return nil, err
	}
	res, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Could not get the user")
	}
	user := models.NewHipchatUser()
	err = json.NewDecoder(res.Body).Decode(user)
	return user, err
}

//NewRequest creates a new request with the appropriate headers to make a hipchat call.
func (h *Hipchat) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
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
