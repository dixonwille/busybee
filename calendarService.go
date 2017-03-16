package busybee

//CalendarService is any service that can be used to get meetings from.
type CalendarService interface {
	InEvent(uid string) (bool, error)
}

//UserCalendar is how an individual can call the Calendar service.
type UserCalendar interface {
	InEvent() (bool, error)
}
