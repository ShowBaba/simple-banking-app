package dtos

type LoginDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUpDTO struct {
	Email       string
	FirstName   string
	LastName    string
	Password    string
	PhoneNumber string
	Username    string
}

type LoginResp struct {
	Token string `json:"token"`
}
