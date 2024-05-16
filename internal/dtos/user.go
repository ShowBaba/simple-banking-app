package dtos

type UserDTO struct {
	Email          string
	FirstName      string
	LastName       string
	Password       string
	PhoneNumber    string
	Username       string
	ProfilePicture string
	Transactions   TransactionDTO
}
