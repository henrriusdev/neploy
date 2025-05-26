package service

type Services struct {
	Application   Application
	Gateway       Gateway
	HealthChecker HealthChecker
	Metadata      Metadata
	Onboard       Onboard
	Role          Role
	TechStack     TechStack
	Trace         Trace
	User          User
	Visitor       Visitor
}
