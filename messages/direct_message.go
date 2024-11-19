package messages

// SendDirectMessage represents a request to send a DM
type SendDirectMessage struct {
    FromUserID string
    ToUserID   string
    Content    string
		ParentID   string
}

// SendDirectMessageResponse represents the response to a DM request
type SendDirectMessageResponse struct {
    Success   bool
    MessageID string
    Error     string
}

// GetUserMessages represents a request to get all DMs for a user
type GetUserMessages struct {
    UserID string
}

// GetUserMessagesResponse represents the response containing user's DMs
type GetUserMessagesResponse struct {
    Success   bool
    Messages  []DirectMessage
    Error     string
}

// DirectMessage represents a single DM
type DirectMessage struct {
    MessageID  string
    FromUserID string
    ToUserID   string
    Content    string
    Timestamp  int64
		ParentID   string
}

