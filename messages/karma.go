package messages

type UpdateKarma struct {
    UserID string
    Change int
}

type GetKarma struct {
    UserID string
}

type GetKarmaResponse struct {
    Success bool
    Karma   int
    Error   string
}
