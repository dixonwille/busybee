package busybee

import (
	"errors"
	"fmt"
)

var eventServices = make(map[string]*EventService)

//InEventer is any service that can tell if UID is in an event.
type InEventer interface {
	InEvent(uid string) (bool, error)
}

//EventServiceCreator is a function that will create a new instance of an InEventer.
type EventServiceCreator func(interface{}) (InEventer, error)

//EventService is the foundation of a service that can retrieve whether someone is in an event or not.
type EventService struct {
	Create       EventServiceCreator
	CreateConfig ConfigCreator
}

//RegisterEventService registers an EventService instance with BusyBee.
func RegisterEventService(name string, event *EventService) error {
	if name == "" {
		return errors.New("Must supply a name to register a calendar service")
	}
	if event == nil {
		return errors.New("Must supply a Event Service to register")
	}
	if _, exists := eventServices[name]; exists {
		return fmt.Errorf("%s is already registered as an Event Service", name)
	}
	eventServices[name] = event
	return nil
}

//GetEventService returns an instance of a registered EventService by name.
func GetEventService(name string) (*EventService, error) {
	event, ok := eventServices[name]
	if !ok {
		return nil, fmt.Errorf("Could not find %s in the list of registered Event Services", name)
	}
	return event, nil
}

//ListEventService returns a slice of all registered event services.
func ListEventService() []string {
	tmp := make([]string, 0, len(eventServices))
	for k := range eventServices {
		tmp = append(tmp, k)
	}
	return tmp
}
