package messages

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

type Post struct {
	PostId        string
	Title         string
	Content       string
	AuthorId      string
	SubredditName string
	Timestamp     int64
}
