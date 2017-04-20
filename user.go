package busybee

//User is the link between the calendar and status.
type User struct {
	CalendarUID   string
	StatusUID     string
	statusService UpdateStatuser
	eventService  InEventer
}

//NewUser creates a new user link betweent the statuser and calendarer.
func NewUser(calendarUID, statusUID string, statusService UpdateStatuser, eventService InEventer) *User {
	return &User{
		CalendarUID:   calendarUID,
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
	return u.eventService.InEvent(u.CalendarUID)
}
