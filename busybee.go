package busybee

//ServiceType are the different types of services BusyBee supports.
type ServiceType int

const (
	//ServiceTypeEvent is for event services.
	ServiceTypeEvent ServiceType = iota
	//ServiceTypeStatus is for status services.
	ServiceTypeStatus
)

//ConfigCreator is a function that should return a pointer to an empty configuration for that service.
type ConfigCreator func() interface{}

//BusyBee are configuration values that are not dependent on the the plugins installed.
type BusyBee struct {
	EventUID   string `quest:"What is the unique identifier for the Event Service?"`
	StatusUID  string `quest:"What is the unique identifier for the Status Service?"`
	PrivateKey string
	Plugins    map[string]PluginConfig
}

//PluginConfig is the structure of a Plugin for configuration
type PluginConfig struct {
	Type   ServiceType
	Config interface{}
}

//User is the link between the calendar and status.
type User struct {
	busyBee       *BusyBee
	statusService UpdateStatuser
	eventService  InEventer
}

//NewUser creates a new user link betweent the statuser and calendarer.
func (bb *BusyBee) NewUser(eventService InEventer, statusService UpdateStatuser) *User {
	return &User{
		busyBee:       bb,
		statusService: statusService,
		eventService:  eventService,
	}
}

//UpdateStatus updates the status for the user.
func (u *User) UpdateStatus(status Status) error {
	return u.statusService.UpdateStatus(u.busyBee.StatusUID, status)
}

//InEvent returns true if the user is in an event and false otherwise.
func (u *User) InEvent() (bool, error) {
	return u.eventService.InEvent(u.busyBee.EventUID)
}
