package user

import "time"

type User struct {
	ID        string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FirstName string `sql:"type:varchar(255)"`
	LastName  string `sql:"type:varchar(255)"`
	Login     string `sql:"type:varchar(255)"`
	Password  string `sql:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
