package service

import (
	"fmt"
	"regexp"
)

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

var semverRegex = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

func validateVersionTag(version string) error {
	if !semverRegex.MatchString(version) {
		return fmt.Errorf("invalid version tag format, expected 'vX.Y.Z'")
	}
	return nil
}
