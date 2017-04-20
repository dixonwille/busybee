package busybee

//ConfigCreator is a function that should return a pointer to an empty configuration for that service.
type ConfigCreator func() interface{}
