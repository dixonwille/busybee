package busybee

//Service is used for any services that can be registered to busybee
type Service struct {
	CreateConfig ConfigCreator
}
