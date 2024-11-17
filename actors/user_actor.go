package actors

import (
	"reddit/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type UserActor struct {
    // In-memory storage
    users map[string]string // username -> password
    karma map[string]int // userID -> karma
}

func NewUserActor() *UserActor {
    return &UserActor{
        users: make(map[string]string),
        karma: make(map[string]int),
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
                response.Token = "reddit-token-" + msg.Username // Simple token for demo
            } else {
                response.Success = false
                response.Error = "Invalid credentials"
            }

            context.Respond(response)

        case *messages.UpdateKarma:
            if _, exists := state.users[msg.UserID]; exists {
                state.karma[msg.UserID] += msg.Change
            }

        case *messages.GetKarma:
            if _, exists := state.users[msg.UserID]; !exists {
                context.Respond(&messages.GetKarmaResponse{
                    Success: false,
                    Error:   "User not found",
                })
                return
            }

            context.Respond(&messages.GetKarmaResponse{
                Success: true,
                Karma:   state.karma[msg.UserID],
            })
    }
}
