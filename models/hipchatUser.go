package models

//HipchatUser is a user in hipchat.
type HipchatUser struct {
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

//HipchatUserPresence is part of the body in a hipchat user.
//This states whether the user is busy or not "show" and a reason "status".
type HipchatUserPresence struct {
	Status string `json:"status"`
	Show   string `json:"show"`
}
