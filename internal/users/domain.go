package users

import "time"

type User struct {
	ID        int64
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
}
