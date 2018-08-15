package api

type CreateUserScheme struct {
	UserID string `json:"user_id"`
}

type AddUserChatScheme struct {
	ChatID int64 `json:"chat_id"`
}
