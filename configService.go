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

//MainConfig are configuration values that are not dependent on the the plugins installed.
type MainConfig struct {
	EventUID  string `quest:"What is the unique identifier for the Event Service?"`
	StatusUID string `quest:"What is the unique identifier for the Status Service?"`
	Plugins   map[string]PluginConfig
}

//PluginConfig is the structure of a Plugin for configuration
type PluginConfig struct {
	Type   ServiceType
	Config interface{}
}
