package messages

// User related messages
type RegisterUser struct {
	Username string
	Password string
}

type RegisterUserResponse struct {
	Success bool
	Error   string
	UserId  string
}

type LoginUser struct {
	Username string
	Password string
}

type LoginUserResponse struct {
	Success bool
	Error   string
	Token   string
}
