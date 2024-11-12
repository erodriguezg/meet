package domain

const (
	PermissionCodeManageSystem   = 1
	PermissionCodeEditOwnProfile = 2
	PermissionCodeCreateRoom     = 3
)

type Permission struct {
	Code int    `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}
