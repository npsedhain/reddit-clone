package messages

type CreateComment struct {
    PostID      string
    ParentID    string // Empty for top-level comments
    Content     string
    AuthorID    string
}

type CreateCommentResponse struct {
    Success   bool
    CommentID string
    Error     string
}

type GetPostComments struct {
    PostID string
}

type GetPostCommentsResponse struct {
    Success  bool
    Comments []Comment
    Error    string
}

type Comment struct {
    ID        string
    PostID    string
    ParentID  string
    Content   string
    AuthorID  string
    Children  []Comment
    CreatedAt int64
}
