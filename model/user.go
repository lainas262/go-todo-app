package model

type User struct {
	Id           int64   `db:"id"`
	FirstName    string  `db:"first_name"`
	LastName     string  `db:"last_name"`
	Email        string  `db:"email"`
	UserName     string  `db:"username"`
	PasswordHash *string `db:"password_hash"`
}
