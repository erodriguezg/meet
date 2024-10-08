package domain

import "time"

type PersonFilter struct {
	Id            *string    `json:"id"`
	NameLike      *string    `json:"nameLike"`
	BirthdayLower *time.Time `json:"birthdayLower"`
	BirthdayUpper *time.Time `json:"birthdayUpper"`
}
