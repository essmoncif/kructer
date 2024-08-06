package core

import "time"

type Session struct {
	UserID    string
	Login     string
	Roles     []string
	Token     string
	ExpiresAt time.Time
}
