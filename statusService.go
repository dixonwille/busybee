package busybee

import (
	"errors"
	"fmt"
)

var statusServiceFactories = make(map[string]StatusServiceFactory)

//Status is an integer representing possable status.
type Status int

const (
	//StatusUnknown is used if a service could not figure out what the status of the user was.
	StatusUnknown Status = iota
	//StatusBusy is used to state that the user is busy.
	StatusBusy
	//StatusAvailable is used to state that the user is Available.
	StatusAvailable
)

//StatusService is any service that can be used to update your status.
type StatusService interface {
	UpdateStatus(uid string, status Status) error
}

//UserStatus is how an individual will call the Status Service.
type UserStatus interface {
	UpdateStatus(status Status) error
}

//StatusServiceFactory is a function that will create a new instance of a StatusService.
type StatusServiceFactory func(map[string]string) (StatusService, error)

//RegisterStatusService registers a StatusService instance with BusyBee
func RegisterStatusService(name string, factory StatusServiceFactory) error {
	if name == "" {
		return errors.New("Must supply a name to register a status service")
	}
	if factory == nil {
		return errors.New("Must supply a StatusServiceFactory to register")
	}
	if _, exists := statusServiceFactories[name]; exists {
		return fmt.Errorf("%s is already registered as a StatusService", name)
	}
	statusServiceFactories[name] = factory
	return nil
}

//CreateStatusService creates a new status service instance based on the supplied configuration values.
//service must be supplied in conf so that it knows which registered service you want an instance of.
func CreateStatusService(conf map[string]string) (StatusService, error) {
	if conf == nil {
		return nil, errors.New("Must supply a map for configuration values")
	}
	service, ok := conf["service"]
	if !ok {
		return nil, errors.New("Must supply a service you want to create")
	}
	status, ok := statusServiceFactories[service]
	if !ok {
		return nil, fmt.Errorf("Could not find %s in the list of registered status services", service)
	}
	return status(conf)
}
