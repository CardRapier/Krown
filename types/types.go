package types

type RegisterUserPayload struct {
	UserName string `json:"username"`
	Email string  `json:"email"`
	Password string `json:"password"`
}