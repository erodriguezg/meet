package dto

import "github.com/erodriguezg/meet/pkg/util/datetime"

type RoomDTO struct {
	RoomHash            string        `json:"roomHash"`
	Owner               RoomOwnerDTO  `json:"owner"`
	CreationDate        datetime.Date `json:"creationDate"`
	AnonymousAccess     bool          `json:"anonymousAccess"`
	LastInteractionDate datetime.Date `json:"lastInteractionDate"`
}

type RoomOwnerDTO struct {
	PersonId  string `json:"personId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CreateRoomDTO struct {
	OwnerPersonId   string `json:"ownerPersonId"`
	AnonymousAccess bool   `json:"anonymousAccess"`
}

type ChangeRoomVisibilityRoomDTO struct {
	OwnerPersonId      string `json:"ownerPersonId"`
	RoomHash           string `json:"roomHash"`
	NewAnonymousAccess bool   `json:"newAnonymousAccess"`
}
