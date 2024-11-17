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

// Subreddit related messages
type CreateSubreddit struct {
    Name        string
    Description string
    CreatorId   string
}

type CreateSubredditResponse struct {
    Success bool
    Error   string
    SubId   string
}

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
}

type JoinSubreddit struct {
    UserId    string
    SubredditName string
}

type JoinSubredditResponse struct {
    Success bool
    Error   string
}

type GetSubredditMembers struct {
    SubredditName string
}

type GetSubredditMembersResponse struct {
    Members []string
    Success bool
    Error   string
}

type LeaveSubreddit struct {
    UserId        string
    SubredditName string
}

type LeaveSubredditResponse struct {
    Success bool
    Error   string
}
