package models

type Profile struct {
	UserID string
	Email  string
}

func (profile Profile) IsEmpty() bool {
	return profile.UserID == "" || profile.Email == ""
}
