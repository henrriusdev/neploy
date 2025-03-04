package model

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type OAuthResponse struct {
	Provider Provider `json:"provider"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
}

type UserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
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
	CpuUsage       float64  `json:"cpuUsage"`
	MemoryUsage    float64  `json:"memoryUsage"`
	Uptime         string   `json:"uptime"`
	RequestsPerMin int      `json:"requestsPerMin"`
	Logs           []string `json:"logs"`
}
