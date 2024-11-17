package messages

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

type GetSubredditPosts struct {
	SubredditName string
}

type GetSubredditPostsResponse struct {
	Success bool
	Error   string
	Posts   []Post
}
