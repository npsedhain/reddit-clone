package messages

import "github.com/asynkron/protoactor-go/actor"

// User related messages
type RegisterUser struct {
	Username string
	Password string
}

type RegisterUserResponse struct {
	Success bool
	Error   string
	UserId  string
	ActorPID *actor.PID
}

type LoginUser struct {
	Username string
	Password string
	ActorPID *actor.PID
}

type LoginUserResponse struct {
	Success bool
	Error   string
	Token   string
}
