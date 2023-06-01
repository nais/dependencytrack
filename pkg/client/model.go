package client

type Permission string

const (
	AccessManagementPermission        = Permission("ACCESS_MANAGEMENT")
	PolicyManagementPermission        = Permission("POLICY_MANAGEMENT")
	PolicyViolationAnalysisPermission = Permission("POLICY_VIOLATION_ANALYSIS")
	SystemConfigurationPermission     = Permission("SYSTEM_CONFIGURATION")
	ViewPortfolioPermission           = Permission("VIEW_PORTFOLIO")
	ViewVulnerabilityPermission       = Permission("VIEW_VULNERABILITY")
	ViewPolicyViolationPermission     = Permission("VIEW_POLICY_VIOLATION")
)

type User struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type NewUser struct {
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
	Uuid      string   `json:"uuid,omitempty"`
	Name      string   `json:"name,omitempty"`
	OidcUsers []User   `json:"oidcUsers,omitempty"`
	ApiKeys   []ApiKey `json:"apiKeys,omitempty"`
}

type ApiKey struct {
	Key string `json:"key,omitempty"`
}

type BomSubmitRequest struct {
	ProjectName    string `json:"projectName"`
	ProjectVersion string `json:"projectVersion"`
	AutoCreate     bool   `json:"autoCreate"`
	Bom            string `json:"bom"`
}

type Project struct {
	Active     bool     `json:"active"`
	Author     string   `json:"author"`
	Classifier string   `json:"classifier"`
	Group      string   `json:"group"`
	Name       string   `json:"name"`
	Publisher  string   `json:"publisher"`
	Tags       []Tag    `json:"tags"`
	Uuid       string   `json:"uuid"`
	Version    string   `json:"version"`
	Parent     *Project `json:"parent"`
}

type Tag struct {
	Name string `json:"name"`
}

type Tags struct {
	Tags []Tag `json:"tags"`
}
