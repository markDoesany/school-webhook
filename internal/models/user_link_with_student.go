package models

type UserLinkWithStudent struct {
	UserLink
	Student *StudentProfile `json:"student,omitempty"`
}
