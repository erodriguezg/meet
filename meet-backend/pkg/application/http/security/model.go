package security

type Identity struct {
	User            User    `json:"user"`
	Company         Company `json:"company"`
	UserGroupName   *string `json:"userGroupName"`
	Profile         Profile `json:"profile"`
	PermissionCodes []int   `json:"permissionCodes"`
}

type User struct {
	Rut   int    `json:"rut"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Company struct {
	Rut  int    `json:"rut"`
	Name string `json:"name"`
}

type Profile struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

type BypassRequest struct {
	Email    string  `json:"email"`
	Password *string `json:"password"`
}

type LoginSsoRequest struct {
	Token string `json:"token"`
}
