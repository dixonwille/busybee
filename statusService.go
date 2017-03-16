package busybee

//Status is an integer representing possable status.
type Status int

const (
	//Busy is used to state that the user is busy.
	Busy Status = iota
	//Available is used to state that the user is Available.
	Available
)

//StatusService is any service that can be used to update your status.
type StatusService interface {
	UpdateStatus(uid string, status Status) error
	CurrentStatus(uid string) (Status, error)
}

//UserStatus is how an individual will call the Status Service.
type UserStatus interface {
	UpdateStatus(status Status) error
	CurrentStatus() (Status, error)
}
