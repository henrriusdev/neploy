package model

type BaseEntity struct {
	ID        string `json:"id" db:"id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
	DeletedAt string `json:"deleted_at" db:"deleted_at"`
}

type BaseRelation struct {
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
	DeletedAt string `json:"deleted_at" db:"deleted_at"`
}

type User struct {
	BaseEntity
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	DOB       string `json:"dob" db:"dob"`
	Address   string `json:"address" db:"address"`
	Phone     string `json:"phone" db:"phone"`
}

type Role struct {
	BaseEntity
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type Application struct {
	BaseEntity
	AppName         string `json:"app_name" db:"app_name"`
	StorageLocation string `json:"storage_location" db:"storage_location"`
	DeployLocation  string `json:"deploy_location" db:"deploy_location"`
	TechStackID     string `json:"tech_stack_id" db:"tech_stack_id"`
}

type TechStack struct {
	BaseEntity
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type UserRoles struct {
	BaseRelation
	UserID string `json:"user_id" db:"user_id"`
	RoleID string `json:"role_id" db:"role_id"`
}

type UserTechStack struct {
	BaseRelation
	UserID      string `json:"user_id" db:"user_id"`
	TechStackID string `json:"tech_stack_id" db:"tech_stack_id"`
}

type Traces struct {
	BaseEntity
	UserID          string `json:"user_id" db:"user_id"`
	Type            string `json:"type" db:"type"`
	Action          string `json:"action" db:"action"`
	ActionTimestamp Date   `json:"action_timestamp" db:"action_timestamp"`
	SqlStatement    string `json:"sql_statement" db:"sql_statement"`
}

type Gateway struct {
	BaseEntity
	Name            string `json:"name" db:"name"`
	EndpointURL     string `json:"endpoint_url" db:"endpoint_url"`
	EndpointType    string `json:"endpoint_type" db:"endpoint_type"`
	Stage           string `json:"stage" db:"stage"`
	HttpMethod      string `json:"http_method" db:"http_method"`
	IntegrationType string `json:"integration_type" db:"integration_type"`
	LoggingLevel    string `json:"logging_level" db:"logging_level"`
	ApplicationID   string `json:"application_id" db:"application_id"`
}

type Environment struct {
	BaseEntity
	Name string `json:"name" db:"name"`
}

type ApplicationEnvironment struct {
	BaseRelation
	ApplicationID string `json:"application_id" db:"application_id"`
	EnvironmentID string `json:"environment_id" db:"environment_id"`
}

type RefreshToken struct {
	BaseEntity
	UserID string `json:"user_id" db:"user_id"`
	Token  string `json:"token" db:"token"`
}

type ApplicationStats struct {
	BaseEntity
	ApplicationID       string `json:"application_id" db:"application_id"`
	EnvironmentID       string `json:"environment_id" db:"environment_id"`
	Date                Date   `json:"date" db:"date"`
	Requests            int    `json:"requests" db:"requests"`
	Errors              int    `json:"errors" db:"errors"`
	AverageResponseTime int    `json:"average_response_time" db:"average_response_time"`
	DataTransfered      int    `json:"data_transfered" db:"data_transfered"`
	UniqueVisitors      int    `json:"unique_visitors" db:"unique_visitors"`
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
	UserID      string `json:"user_id" db:"user_id"`
	Provider    string `json:"provider" db:"provider"`
	OAuthID     string `json:"oauth_id" db:"oauth_id"`
	AccessToken string `json:"access_token" db:"access_token"`
}
