package client

import "time"

type (
	Permission string
	TagPrefix  string
)

const (
	AccessManagementPermission        = Permission("ACCESS_MANAGEMENT")
	PolicyManagementPermission        = Permission("POLICY_MANAGEMENT")
	PolicyViolationAnalysisPermission = Permission("POLICY_VIOLATION_ANALYSIS")
	SystemConfigurationPermission     = Permission("SYSTEM_CONFIGURATION")
	ViewPolicyViolationPermission     = Permission("VIEW_POLICY_VIOLATION")
	ViewPortfolioPermission           = Permission("VIEW_PORTFOLIO")
	ViewVulnerabilityPermission       = Permission("VIEW_VULNERABILITY")
	WorkloadTagPrefix                 = TagPrefix("workload:")
	EnvironmentTagPrefix              = TagPrefix("env:")
	TeamTagPrefix                     = TagPrefix("team:")
	ProjectTagPrefix                  = TagPrefix("project:")
	ImageTagPrefix                    = TagPrefix("image:")
	VersionTagPrefix                  = TagPrefix("version:")
	RekorTagPrefix                    = TagPrefix("rekor:")
	DigestTagPrefix                   = TagPrefix("digest:")
)

func (t TagPrefix) String() string {
	return string(t)
}

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
	ParentUuid     string `json:"parentUUID,omitempty"`
	AutoCreate     bool   `json:"autoCreate"`
	Bom            string `json:"bom"`
}

type Project struct {
	Active                 bool           `json:"active"`
	Author                 string         `json:"author"`
	Classifier             string         `json:"classifier"`
	Group                  string         `json:"group"`
	Name                   string         `json:"name"`
	LastBomImportFormat    string         `json:"lastBomImportFormat,omitempty"`
	LastBomImport          int64          `json:"lastBomImport,omitempty"`
	LastInheritedRiskScore float64        `json:"lastInheritedRiskScore,omitempty"`
	Publisher              string         `json:"publisher"`
	Tags                   []Tag          `json:"tags"`
	Uuid                   string         `json:"uuid"`
	Version                string         `json:"version"`
	Parent                 *Project       `json:"parent"`
	Metrics                *ProjectMetric `json:"metrics,omitempty"`
}

type ProjectMetric struct {
	Critical             int     `json:"critical"`
	High                 int     `json:"high"`
	Medium               int     `json:"medium"`
	Low                  int     `json:"low"`
	Unassigned           int     `json:"unassigned"`
	Vulnerabilities      int     `json:"vulnerabilities"`
	VulnerableComponents int     `json:"vulnerableComponents"`
	Components           int     `json:"components"`
	Suppressed           int     `json:"suppressed"`
	FindingsTotal        int     `json:"findingsTotal"`
	FindingsAudited      int     `json:"findingsAudited"`
	FindingsUnaudited    int     `json:"findingsUnaudited"`
	InheritedRiskScore   float64 `json:"inheritedRiskScore"`
	FirstOccurrence      int64   `json:"firstOccurrence"`
	LastOccurrence       int64   `json:"lastOccurrence"`
}

type Tag struct {
	Name string `json:"name"`
}

type Tags struct {
	Tags []Tag `json:"tags"`
}

type ConfigProperty struct {
	GroupName     string `json:"groupName"`
	PropertyName  string `json:"propertyName"`
	PropertyValue string `json:"propertyValue"`
	PropertyType  string `json:"propertyType"`
	Description   string `json:"description"`
}

type Component struct {
	UUID    string `json:"uuid"`
	PURL    string `json:"purl"`
	Project string `json:"project"`
	Name    string `json:"name"`
}

type Vulnerability struct {
	VulnId       string  `json:"vulnId"`
	Severity     string  `json:"severity"`
	SeverityRank int     `json:"severityRank"`
	Source       string  `json:"source"`
	Aliases      []Alias `json:"aliases"`
	Title        string  `json:"title"`
}

type Alias struct {
	CveId  string `json:"cveId"`
	GhsaId string `json:"ghsaId"`
}

type Finding struct {
	Component     Component     `json:"component"`
	Vulnerability Vulnerability `json:"vulnerability"`
}

type AnalysisRequest struct {
	Project               string `json:"project"`
	Component             string `json:"component"`
	Vulnerability         string `json:"vulnerability"`
	AnalysisState         string `json:"analysisState"`
	AnalysisJustification string `json:"analysisJustification"`
	AnalysisResponse      string `json:"analysisResponse"`
	AnalysisDetails       string `json:"analysisDetails"`
	Comment               string `json:"comment"`
	IsSuppressed          bool   `json:"isSuppressed"`
}

type Analysis struct {
	AnalysisState         string            `json:"analysisState"`
	AnalysisJustification string            `json:"analysisJustification"`
	AnalysisResponse      string            `json:"analysisResponse"`
	AnalysisDetails       string            `json:"analysisDetails"`
	AnalysisComments      []AnalysisComment `json:"analysisComments"`
	IsSuppressed          bool              `json:"isSuppressed"`
}

type AnalysisComment struct {
	Timestamp time.Time `json:"timestamp"`
	Comment   string    `json:"comment"`
	Commenter string    `json:"commenter"`
}
