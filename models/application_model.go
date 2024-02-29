package models

type ApplicationModel struct {
	AppId       uint
	UserId      uint
	AppName     string
	CreatedAt   string
	Description string
	Logs        []LogModel
}
