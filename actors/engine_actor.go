package actors

import (
	"fmt"
	"reddit/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type EngineActor struct {
	userActors           []*actor.PID
	postActors           []*actor.PID
	subredditActors      []*actor.PID
	directMessageActors  []*actor.PID
	commentActors        []*actor.PID
	currentUserActor     int
}

func NewEngineActor() *EngineActor {
	return &EngineActor{
		userActors:           make([]*actor.PID, 0),
		postActors:           make([]*actor.PID, 0),
		subredditActors:      make([]*actor.PID, 0),
		directMessageActors:  make([]*actor.PID, 0),
		commentActors:        make([]*actor.PID, 0),
		currentUserActor:     0,
	}
}

func (state *EngineActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
			// Create 10 instances of each actor type
			for i := 0; i < 10; i++ {
					// Create user actor
					userProps := actor.PropsFromProducer(func() actor.Actor { return NewUserActor() })
					userPID, _ := context.SpawnNamed(userProps, fmt.Sprintf("user_actor_%d", i))
					state.userActors = append(state.userActors, userPID)

					// Create post actor
					postProps := actor.PropsFromProducer(func() actor.Actor { return NewPostActor() })
					postPID, _ := context.SpawnNamed(postProps, fmt.Sprintf("post_actor_%d", i))
					state.postActors = append(state.postActors, postPID)

					// Create subreddit actor
					subredditProps := actor.PropsFromProducer(func() actor.Actor { return NewSubredditActor() })
					subredditPID, _ := context.SpawnNamed(subredditProps, fmt.Sprintf("subreddit_actor_%d", i))
					state.subredditActors = append(state.subredditActors, subredditPID)

					// Create direct message actor
					dmProps := actor.PropsFromProducer(func() actor.Actor { return NewDirectMessageActor() })
					dmPID, _ := context.SpawnNamed(dmProps, fmt.Sprintf("dm_actor_%d", i))
					state.directMessageActors = append(state.directMessageActors, dmPID)

					// Create comment actor
					commentProps := actor.PropsFromProducer(func() actor.Actor { return NewCommentActor() })
					commentPID, _ := context.SpawnNamed(commentProps, fmt.Sprintf("comment_actor_%d", i))
					state.commentActors = append(state.commentActors, commentPID)
			}

		case *messages.RegisterUser:
			// Forward to next user actor in round-robin fashion
			userActor := state.userActors[state.currentUserActor]
			state.currentUserActor = (state.currentUserActor + 1) % len(state.userActors)
			context.RequestWithCustomSender(userActor, msg, context.Sender())

	case *messages.LoginUser:
			userActor := state.userActors[state.currentUserActor]
			state.currentUserActor = (state.currentUserActor + 1) % len(state.userActors)
			context.RequestWithCustomSender(userActor, msg, context.Sender())

	case *messages.CreateSubreddit:
			subredditActor := state.subredditActors[state.currentUserActor]
			state.currentUserActor = (state.currentUserActor + 1) % len(state.subredditActors)
			context.RequestWithCustomSender(subredditActor, msg, context.Sender())

	case *messages.Post:
			postActor := state.postActors[state.currentUserActor]
			state.currentUserActor = (state.currentUserActor + 1) % len(state.postActors)
			context.RequestWithCustomSender(postActor, msg, context.Sender())

	case *messages.CreateComment:
			commentActor := state.commentActors[state.currentUserActor]
			state.currentUserActor = (state.currentUserActor + 1) % len(state.commentActors)
			context.RequestWithCustomSender(commentActor, msg, context.Sender())

	case *messages.SendDirectMessage:
			dmActor := state.directMessageActors[state.currentUserActor]
			state.currentUserActor = (state.currentUserActor + 1) % len(state.directMessageActors)
			context.RequestWithCustomSender(dmActor, msg, context.Sender())

	case *messages.Vote:
		switch msg.Type {
		case "post":
				postActor := state.postActors[state.currentUserActor]
				state.currentUserActor = (state.currentUserActor + 1) % len(state.postActors)
				context.RequestWithCustomSender(postActor, msg, context.Sender())
		case "comment":
				commentActor := state.commentActors[state.currentUserActor]
				state.currentUserActor = (state.currentUserActor + 1) % len(state.commentActors)
				context.RequestWithCustomSender(commentActor, msg, context.Sender())
		}
	}
}
