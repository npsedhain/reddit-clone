package messages

type Vote struct {
    TargetID string // Can be either post ID or comment ID
    UserID   string
    IsUpvote bool
    Type     string // "post" or "comment"
}

type VoteResponse struct {
    Success bool
    Error   string
}
