package service

import "errors"

type Services struct {
	User     User
	Role     Role
	Onboard  Onboard
	Metadata Metadata
}

func NotFound(entity string) error {
	return errors.New(entity + " not found")
}

func AlreadyExists(entity string) error {
	return errors.New(entity + " already exists")
}

func Invalid(entity string) error {
	return errors.New("invalid " + entity)
}

func InvalidField(field string) error {
	return errors.New("invalid " + field)
}

func InvalidValue(value string) error {
	return errors.New("invalid value: " + value)
}

func InvalidLength(field string, length int) error {
	return errors.New("invalid length for " + field + ": " + string(length))
}

func InvalidEmail() error {
	return errors.New("invalid email")
}

func InvalidPassword() error {
	return errors.New("invalid password")
}

func InvalidToken() error {
	return errors.New("invalid token")
}

func InvalidCredentials() error {
	return errors.New("invalid credentials")
}

func Unauthorized() error {
	return errors.New("unauthorized")
}

func Forbidden() error {
	return errors.New("forbidden")
}

func InternalError() error {
	return errors.New("internal error")
}

func NotImplemented() error {
	return errors.New("not implemented")
}

func NotImplementedYet() error {
	return errors.New("not implemented yet")
}

func NotImplementedFor(entity string) error {
	return errors.New("not implemented for " + entity)
}

func NotImplementedForField(field string) error {
	return errors.New("not implemented for " + field)
}

func NoResults() error {
	return errors.New("no results")
}

func NoResultsFor(entity string) error {
	return errors.New("no results for " + entity)
}

func NoResultsForFilter(filter string) error {
	return errors.New("no results for filter: " + filter)
}

func IsSqlNoRows(err error) bool {
	return errors.Is(err, NoResults())
}
