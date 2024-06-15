package models

type User struct {
	ID    string `db:"id" json:"id"`
	Email string `db:"email" json:"email"`
	Roles Roles  `db:"roles" json:"roles"`
}
