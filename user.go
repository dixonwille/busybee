package busybee

//User is the link between the calendar and status.
type User struct {
	EventUID      string
	StatusUID     string
	statusService UpdateStatuser
	eventService  InEventer
}

//NewUser creates a new user link betweent the statuser and calendarer.
func NewUser(eventUID, statusUID string, eventService InEventer, statusService UpdateStatuser) *User {
	return &User{
		EventUID:      eventUID,
		StatusUID:     statusUID,
		statusService: statusService,
		eventService:  eventService,
	}
}

//UpdateStatus updates the status for the user.
func (u *User) UpdateStatus(status Status) error {
	return u.statusService.UpdateStatus(u.StatusUID, status)
}

//InEvent returns true if the user is in an event and false otherwise.
func (u *User) InEvent() (bool, error) {
	return u.eventService.InEvent(u.EventUID)
}
