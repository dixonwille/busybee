package busybee

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
