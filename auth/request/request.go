package request

type SignupRequest struct {
	Username string `json:"username" validate:"required,min=2"`
	Fullname string `json:"fullname" validate:"required"`
	Password string `json:"password" validate:"required,min=5"`
}
