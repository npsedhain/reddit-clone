package actors

import (
	"reddit/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type UserActor struct {
    // In-memory storage for demo purposes
    users map[string]string // username -> password
}

func NewUserActor() *UserActor {
    return &UserActor{
        users: make(map[string]string),
    }
}

func (state *UserActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *messages.RegisterUser:
        response := &messages.RegisterUserResponse{}

        if _, exists := state.users[msg.Username]; exists {
            response.Success = false
            response.Error = "Username already exists"
        } else {
            state.users[msg.Username] = msg.Password
            response.Success = true
            response.UserId = msg.Username // Using username as ID for simplicity
        }

        context.Respond(response)

    case *messages.LoginUser:
        response := &messages.LoginUserResponse{}

        if password, exists := state.users[msg.Username]; exists && password == msg.Password {
            response.Success = true
            response.Token = "dummy-token-" + msg.Username // Simple token for demo
        } else {
            response.Success = false
            response.Error = "Invalid credentials"
        }

        context.Respond(response)
    }
}
