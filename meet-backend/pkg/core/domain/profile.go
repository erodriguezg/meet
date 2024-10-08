package domain

const (
	ProfileCodeAdministrator = 1
	ProfileCodeUser          = 2
	ProfileCodeModel         = 3
	ProfileCodeModerator     = 4
)

type Profile struct {
	Code             int    `json:"code" bson:"code"`
	Name             string `json:"name" bson:"name"`
	PermissionsCodes []int  `json:"permissionsCodes" bson:"permissionsCodes"`
}
