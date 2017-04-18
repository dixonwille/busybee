package busybee

import (
	"errors"
	"fmt"
)

var calendarServiceFactories = make(map[string]CalendarServiceFactory)

//CalendarService is any service that can be used to get meetings from.
type CalendarService interface {
	InEvent(uid string) (bool, error)
}

//UserCalendar is how an individual can call the Calendar service.
type UserCalendar interface {
	InEvent() (bool, error)
}

//CalendarServiceFactory is a function that will create a new instance of a CalendarService.
type CalendarServiceFactory func(map[string]string) (CalendarService, error)

//RegisterCalendarService registers a CalendarService instance with BusyBee
func RegisterCalendarService(name string, factory CalendarServiceFactory) error {
	if name == "" {
		return errors.New("Must supply a name to register a calendar service")
	}
	if factory == nil {
		return errors.New("Must supply a CalendarServiceFactory to register")
	}
	if _, exists := calendarServiceFactories[name]; exists {
		return fmt.Errorf("%s is already registered as a CalendarService", name)
	}
	calendarServiceFactories[name] = factory
	return nil
}

//CreateCalendarService creates a new calendar instance based on the supplied configuration values.
//service must be supplied in conf so that it knows which registered service you want an instance of.
func CreateCalendarService(conf map[string]string) (CalendarService, error) {
	if conf == nil {
		return nil, errors.New("Must supply a map for configuration values")
	}
	service, ok := conf["service"]
	if !ok {
		return nil, errors.New("Must supply a service you want to create")
	}
	calendar, ok := calendarServiceFactories[service]
	if !ok {
		return nil, fmt.Errorf("Could not find %s in the list of registered calendar services", service)
	}
	return calendar(conf)
}
