package actors

import (
	"reddit/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/uuid"
)

type DirectMessageActor struct {
    messages map[string][]messages.DirectMessage // UserID -> Messages
}

func NewDirectMessageActor() *DirectMessageActor {
    return &DirectMessageActor{
        messages: make(map[string][]messages.DirectMessage),
    }
}

func (state *DirectMessageActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *messages.SendDirectMessage:
        messageID := uuid.New().String()
        dm := messages.DirectMessage{
            MessageID:  messageID,
            FromUserID: msg.FromUserID,
            ToUserID:   msg.ToUserID,
            Content:    msg.Content,
            ParentID:   msg.ParentID,
            Timestamp:  time.Now().Unix(),
        }

        // Store message for both sender and receiver
        state.messages[msg.FromUserID] = append(state.messages[msg.FromUserID], dm)
        state.messages[msg.ToUserID] = append(state.messages[msg.ToUserID], dm)

        context.Respond(&messages.SendDirectMessageResponse{
            Success:   true,
            MessageID: messageID,
        })

    case *messages.GetUserMessages:
        userMessages := state.messages[msg.UserID]
        context.Respond(&messages.GetUserMessagesResponse{
            Success:  true,
            Messages: userMessages,
        })
    }
}
