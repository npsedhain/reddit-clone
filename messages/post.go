package messages

import "github.com/asynkron/protoactor-go/actor"

// Post related messages
type CreatePost struct {
	Title     string
	Content   string
	SubId     string
	AuthorId  string
}

type CreatePostResponse struct {
	Success bool
	Error   string
	PostId  string
	ActorPID *actor.PID
}

type Post struct {
	PostId        string
	Title         string
	Content       string
	AuthorId      string
	SubredditName string
	Timestamp     int64
	ActorPID      *actor.PID
}
