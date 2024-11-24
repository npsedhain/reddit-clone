package actors

import (
	"fmt"
	"math/rand"
	"reddit/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type ClientActor struct {
    enginePID     *actor.PID
		controllerPID *actor.PID
    rand          *rand.Rand
    username      string
    token         string
    mySubreddits  []string    // subreddits created/joined by this user
    myPosts       []string    // posts created by this user
    myComments    []string    // comments created by this user
    actionDelay   time.Duration
		startTime     time.Time
}

func NewClientActor(enginePID *actor.PID, controllerPID *actor.PID) *ClientActor {
    return &ClientActor{
        enginePID:    enginePID,
				controllerPID: controllerPID,
        rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
        mySubreddits: make([]string, 0),
        myPosts:      make([]string, 0),
        myComments:   make([]string, 0),
        actionDelay:  time.Duration(rand.Intn(1000)) * time.Millisecond,
				startTime:    time.Now(),
    }
}

func (state *ClientActor) generateContent() string {
    return fmt.Sprintf("content_%d", state.rand.Intn(10000))
}

func (state *ClientActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *actor.Started:
        // Register user when actor starts
        state.register(context)

    case *messages.RegisterUserResponse:
			metricsMsg := &messages.MetricsMessage{
				Action:       "register",
				Success:      msg.Success,
				ResponseTime: time.Since(state.startTime),
				Error:        msg.Error,
			}
			context.Send(state.controllerPID, metricsMsg)
			if msg.Success {
					state.username = msg.UserId
					// Automatically login after successful registration
					state.login(context)
			}

    case *messages.LoginUserResponse:
        if msg.Success {
            state.token = msg.Token
            // Start periodic actions after successful login
            context.SetReceiveTimeout(state.actionDelay)
        }

    case *actor.ReceiveTimeout:
        if state.token != "" {
            state.performRandomAction(context)
        }
        context.SetReceiveTimeout(state.actionDelay)

    case *messages.CreateSubredditResponse:
        if msg.Success {
            state.mySubreddits = append(state.mySubreddits, msg.SubId)
        }

    case *messages.CreatePostResponse:
        if msg.Success {
            state.myPosts = append(state.myPosts, msg.PostId)
        }

    case *messages.CreateCommentResponse:
        if msg.Success {
            state.myComments = append(state.myComments, msg.CommentID)
        }
    }
}

func (state *ClientActor) performRandomAction(context actor.Context) {
    action := state.rand.Intn(6)

    switch action {
    case 0:
        state.createSubreddit(context)
    case 1:
        state.createPost(context)
    case 2:
        state.createComment(context)
    case 3:
        state.sendDirectMessage(context)
    case 4:
        state.voteOnPost(context)
    case 5:
        state.voteOnComment(context)
    }
}

// Individual actions
func (state *ClientActor) register(context actor.Context) {
		state.startTime = time.Now()
    msg := &messages.RegisterUser{
        Username: fmt.Sprintf("user_%d", state.rand.Intn(100000)),
        Password: fmt.Sprintf("pass_%d", state.rand.Intn(100000)),
    }
    context.Request(state.enginePID, msg)
}

func (state *ClientActor) login(context actor.Context) {
    msg := &messages.LoginUser{
        Username: state.username,
        Password: state.username, // Simplified for demo
    }
    context.Request(state.enginePID, msg)
}

func (state *ClientActor) createSubreddit(context actor.Context) {
    msg := &messages.CreateSubreddit{
        Name:        fmt.Sprintf("subreddit_%d", state.rand.Intn(10000)),
        Description: state.generateContent(),
        CreatorId:   state.username,
    }
    context.Request(state.enginePID, msg)
}

func (state *ClientActor) createPost(context actor.Context) {
    if len(state.mySubreddits) == 0 {
        return
    }

    subreddit := state.mySubreddits[state.rand.Intn(len(state.mySubreddits))]
    msg := &messages.Post{
        Title:         fmt.Sprintf("post_%d", state.rand.Intn(10000)),
        Content:       state.generateContent(),
        AuthorId:      state.username,
        SubredditName: subreddit,
    }
    context.Request(state.enginePID, msg)
}

func (state *ClientActor) createComment(context actor.Context) {
    if len(state.myPosts) == 0 {
        return
    }

    postID := state.myPosts[state.rand.Intn(len(state.myPosts))]
    msg := &messages.CreateComment{
        PostID:   postID,
        Content:  state.generateContent(),
        AuthorID: state.username,
    }
    context.Request(state.enginePID, msg)
}

func (state *ClientActor) sendDirectMessage(context actor.Context) {
    msg := &messages.SendDirectMessage{
        FromUserID: state.username,
        ToUserID:   fmt.Sprintf("user_%d", state.rand.Intn(1000)), // Random user
        Content:    state.generateContent(),
    }
    context.Request(state.enginePID, msg)
}

func (state *ClientActor) voteOnPost(context actor.Context) {
    if len(state.myPosts) == 0 {
        return
    }

    postID := state.myPosts[state.rand.Intn(len(state.myPosts))]
    msg := &messages.Vote{
        UserID:    state.username,
        TargetID:  postID,
        IsUpvote:  state.rand.Float32() > 0.5,
				Type:      "post",
    }
    context.Request(state.enginePID, msg)
}

func (state *ClientActor) voteOnComment(context actor.Context) {
    if len(state.myComments) == 0 {
        return
    }

    commentID := state.myComments[state.rand.Intn(len(state.myComments))]
    msg := &messages.Vote{
        UserID:    state.username,
        TargetID:  commentID,
        IsUpvote:  state.rand.Float32() > 0.5,
				Type:      "comment",
    }
    context.Request(state.enginePID, msg)
}
