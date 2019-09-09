package auth

import "errors"

var (
	errNoUsersYet = errors.New("No users yet")
)

func ErrNoUsersYet() error { return errNoUsersYet }

func ErrLoginAsNonExisting(username string) error {
	return errors.New("Attempt to login as non-existing user: " + username)
}
