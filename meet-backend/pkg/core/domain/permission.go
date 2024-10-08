package domain

const (
	PermissionCodeRegisterModel = 1
	PermissionCodeEditAllModels = 2
)

type Permission struct {
	Code int    `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}
