package model

type BaseEntity struct {
	ID        string `json:"id" db:"id" goqu:"skipinsert,skipupdate"`
	CreatedAt Date   `json:"createdAt" db:"created_at" goqu:"skipinsert,skipupdate"`
	UpdatedAt Date   `json:"updatedAt" db:"updated_at" goqu:"skipinsert,skipupdate,omitzero"`
	DeletedAt *Date  `json:"deletedAt" db:"deleted_at" goqu:"skipinsert,skipupdate,omitnil"`
}

type BaseRelation struct {
	CreatedAt Date `json:"createdAt" db:"created_at"`
	UpdatedAt Date `json:"updatedAt" db:"updated_at"`
	DeletedAt Date `json:"deletedAt" db:"deleted_at"`
}

type User struct {
	BaseEntity
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	DOB       Date   `json:"dob" db:"dob"`
	Address   string `json:"address" db:"address"`
	Phone     string `json:"phone" db:"phone"`
	Provider  string `json:"provider" db:"-"`
}

type Role struct {
	BaseEntity
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Icon        string `json:"icon" db:"icon"`
	Color       string `json:"color" db:"color"`
}

type Application struct {
	BaseEntity
	AppName         string  `json:"appName" db:"app_name"`
	Description     string  `json:"description" db:"description"`
	StorageLocation string  `json:"storageLocation" db:"storage_location"`
	DeployLocation  string  `json:"deployLocation" db:"deploy_location"`
	TechStackID     *string `json:"techStackId" db:"tech_stack_id" goqu:"omitempty,omitnil,skipinsert"`
}

type TechStack struct {
	BaseEntity
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type UserRoles struct {
	BaseRelation
	UserID string `json:"userId" db:"user_id"`
	RoleID string `json:"roleId" db:"role_id"`
	User   *User  `json:"user" db:"user"`
	Role   *Role  `json:"role" db:"role"`
}

type UserTechStack struct {
	BaseRelation
	UserID      string     `json:"userId" db:"user_id"`
	TechStackID string     `json:"techStackId" db:"tech_stack_id"`
	User        *User      `json:"user" db:"user"`
	TechStack   *TechStack `json:"techStack" db:"tech_stack"`
}

type Trace struct {
	BaseEntity
	UserID          string `json:"userId" db:"user_id"`
	Type            string `json:"type" db:"type"`
	Action          string `json:"action" db:"action"`
	ActionTimestamp Date   `json:"actionTimestamp" db:"action_timestamp"`
	SqlStatement    string `json:"sqlStatement" db:"sql_statement"`
}

type Gateway struct {
	BaseEntity
	Name            string `json:"name" db:"name"`
	EndpointURL     string `json:"endpointUrl" db:"endpoint_url"`
	EndpointType    string `json:"endpointType" db:"endpoint_type"` // "subdomain" or "path"
	Domain          string `json:"domain" db:"domain"`
	Subdomain       string `json:"subdomain" db:"subdomain"`
	Path            string `json:"path" db:"path"`
	Port            string `json:"port" db:"port"`
	Stage           string `json:"stage" db:"stage"`
	HttpMethod      string `json:"httpMethod" db:"http_method"`
	IntegrationType string `json:"integrationType" db:"integration_type"`
	LoggingLevel    string `json:"loggingLevel" db:"logging_level"`
	ApplicationID   string `json:"applicationId" db:"application_id"`
	Status          string `json:"status" db:"status"` // "active", "inactive", "error"
}

type RefreshToken struct {
	BaseEntity
	UserID string `json:"user_id" db:"user_id"`
	Token  string `json:"token" db:"token"` // hello
}

type ApplicationStat struct {
	BaseEntity
	ApplicationID       string `json:"application_id" db:"application_id"`
	Date                Date   `json:"date" db:"date"`
	Requests            int    `json:"requests" db:"requests"`
	Errors              int    `json:"errors" db:"errors"`
	AverageResponseTime int    `json:"average_response_time" db:"average_response_time"`
	DataTransfered      int    `json:"data_transfered" db:"data_transfered"`
	UniqueVisitors      int    `json:"unique_visitors" db:"unique_visitors"`
	Healthy             bool   `json:"healthy" db:"healthy"`
}

type VisitorInfo struct {
	BaseEntity
	IpAddress string `json:"ip_address" db:"ip_address"`
	Location  string `json:"location" db:"location"`
	Device    string `json:"device" db:"device"`
	Browser   string `json:"browser" db:"browser"`
	Os        string `json:"os" db:"os"`
	VisitedAt Date   `json:"visited_at" db:"visited_at"`
}

type VisitorTrace struct {
	BaseEntity
	ApplicationID string `json:"application_id" db:"application_id"`
	VisitorID     string `json:"visitor_id" db:"visitor_id"`
	PageVisited   string `json:"page_visited" db:"page_visited"`
	VisitDuration int    `json:"visit_duration" db:"visit_duration"`
	VisitedAt     Date   `json:"visited_at" db:"visited_at"`
}

type UserOAuth struct {
	BaseEntity
	UserID   string   `json:"user_id" db:"user_id"`
	Provider Provider `json:"provider" db:"provider"`
	OAuthID  string   `json:"oauth_id" db:"oauth_id"`
}

type Metadata struct {
	BaseEntity
	TeamName string `json:"team_name" db:"team_name"`
	LogoURL  string `json:"logo_url" db:"logo_url"`
	Language string `json:"language" db:"language"`
}

type ApplicationUser struct {
	BaseRelation
	ApplicationID string `json:"application_id" db:"application_id"`
	UserID        string `json:"user_id" db:"user_id"`
}

type Invitation struct {
	BaseEntity
	Email      string `json:"email" db:"email"`
	Role       string `json:"role" db:"role"`
	Token      string `json:"token" db:"token"`
	ExpiresAt  Date   `json:"expires_at" db:"expires_at"`
	AcceptedAt *Date  `json:"accepted_at,omitempty" db:"accepted_at"`
}

type GatewayConfig struct {
	BaseEntity
	DefaultVersioningType VersioningType `json:"defaultVersioningType" db:"default_versioning_type" json:"defaultVersioningType,omitempty"`
	DefaultVersion        VersionType    `json:"defaultVersion" db:"default_version" json:"defaultVersion,omitempty"`
	LoadBalancer          bool           `json:"loadBalancer" db:"load_balancer"`
}

type ApplicationVersion struct {
	BaseEntity
	VersionTag      string `json:"versionTag" db:"version_tag"`
	Description     string `json:"description" db:"description"`
	Status          string `json:"status" db:"status"`                    // Running, Stopped, etc.
	StorageLocation string `json:"StorageLocation" db:"storage_location"` // Aquí va la ruta final al binario/despliegue
	ApplicationID   string `json:"applicationId" db:"application_id"`
}
