package models

import "time"

type LogModel struct {
	LogId     uint
	AppId     uint
	Level     string
	Message   string
	Timestamp time.Time
	Metadata  string
}
