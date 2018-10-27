package api

type CreateUserScheme struct {
	UserName string `json:"user_name"`
}

type AddUserChatScheme struct {
	ChatID int64 `json:"chat_id"`
}
