package busybee

//Userer implements the Statuser and Calendarer interfaces.
//Mainly helpers just so you don't have to pass in the Unique Identifiers for each service.
type Userer interface {
	UserStatus
	UserCalendar
}

//User is the link between the calendar and status.
type User struct {
	CalendarUID     string
	StatusUID       string
	statusService   StatusService
	calendarService CalendarService
}

//NewUser creates a new user link betweent the statuser and calendarer.
func NewUser(calendarUID, statusUID string, statusService StatusService, calendarService CalendarService) *User {
	return &User{
		CalendarUID:     calendarUID,
		StatusUID:       statusUID,
		statusService:   statusService,
		calendarService: calendarService,
	}
}

//UpdateStatus updates the status for the user.
func (u *User) UpdateStatus(status Status) error {
	return u.statusService.UpdateStatus(u.StatusUID, status)
}

//CurrentStatus returns the current status for the user.
func (u *User) CurrentStatus() (Status, error) {
	return u.statusService.CurrentStatus(u.StatusUID)
}

//InEvent returns true if the user is in an event and false otherwise.
func (u *User) InEvent() (bool, error) {
	return u.calendarService.InEvent(u.CalendarUID)
}
