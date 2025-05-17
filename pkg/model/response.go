package model

type LoginResponse struct {
	Token     string   `json:"token"`
	User      User     `json:"user"`
	RoleIDs   []string `json:"roles"`
	RoleNames []string `json:"roleNames"`
}

type OAuthResponse struct {
	Provider Provider `json:"provider"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
}

type UserResponse struct {
	Email    string   `json:"email"`
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Provider string   `json:"provider"`
	Roles    []string `json:"roles"`
}

type TeamMemberResponse struct {
	ID         string      `json:"id"`
	Username   string      `json:"username"`
	Email      string      `json:"email"`
	FirstName  string      `json:"firstName"`
	LastName   string      `json:"lastName"`
	Provider   string      `json:"provider"`
	Roles      []Role      `json:"roles"`
	TechStacks []TechStack `json:"techStacks"`
}

type TeamResponse struct {
	UserRoles     []UserRoles     `json:"userRoles"`
	UserTechStack []UserTechStack `json:"userTechStack"`
}

type FullApplication struct {
	Application
	Stats     []ApplicationStat `json:"stat"`
	TechStack TechStack         `json:"techStack"`
	Status    string            `json:"status"`
}

type FullGateway struct {
	Gateway
	Application Application `json:"application"`
}

type FullUser struct {
	User
	Roles      []Role      `json:"roles"`
	TechStacks []TechStack `json:"techStacks"`
}

type RoleWithUsers struct {
	Role
	Users []User `json:"users"`
}

type TechStackWithApplications struct {
	TechStack
	Applications []Application `json:"applications"`
}

type ApplicationDockered struct {
	Application
	CpuUsage       float64              `json:"cpuUsage"`
	MemoryUsage    float64              `json:"memoryUsage"`
	Uptime         string               `json:"uptime"`
	RequestsPerMin int                  `json:"requestsPerMin"`
	Logs           []string             `json:"logs"`
	Versions       []ApplicationVersion `json:"versions"`
}

type RequestStat struct {
	Hour       string `db:"hour" json:"name"`
	Successful int    `db:"successful" json:"successful"`
	Errors     int    `db:"errors" json:"errors"`
}

type TechStat struct {
	Name  string `json:"name" db:"name"`
	Count uint   `json:"value" db:"count"`
}
