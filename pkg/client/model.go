package client

type Permission string

const ViewPortfolioPermission = Permission("VIEW_PORTFOLIO")

type User struct {
	Username            string `json:"username,omitempty"`
	Email               string `json:"email,omitempty"`
	NewPassword         string `json:"newPassword,omitempty"`
	ConfirmPassword     string `json:"confirmPassword,omitempty"`
	Fullname            string `json:"fullname,omitempty"`
	Suspended           bool   `json:"suspended,omitempty"`
	ForcePasswordChange bool   `json:"forcePasswordChange,omitempty"`
	NonExpiryPassword   bool   `json:"nonExpiryPassword,omitempty"`
}

type AdminUsers struct {
	Users []AdminUser `json:"users,omitempty" yaml:"users,omitempty"`
}

type AdminUser struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Team struct {
	Uuid      string `json:"uuid,omitempty"`
	Name      string `json:"name,omitempty"`
	OidcUsers []User `json:"oidcUsers,omitempty"`
}
