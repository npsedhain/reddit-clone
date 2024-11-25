package messages

import "github.com/asynkron/protoactor-go/actor"

type CreateComment struct {
    PostID      string
    ParentID    string // Empty for top-level comments
    Content     string
    AuthorID    string
    ActorPID    *actor.PID
}

type CreateCommentResponse struct {
    Success   bool
    CommentID string
    Error     string
    ActorPID  *actor.PID
}

type GetPostComments struct {
    PostID string
    ActorPID *actor.PID
}

type GetPostCommentsResponse struct {
    Success  bool
    Comments []Comment
    Error    string
    ActorPID *actor.PID
}

type Comment struct {
    ID        string
    PostID    string
    ParentID  string
    Content   string
    AuthorID  string
    Children  []Comment
    CreatedAt int64
    ActorPID  *actor.PID
}
