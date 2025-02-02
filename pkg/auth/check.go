package auth

import (
	"github.com/elliotxx/osp/pkg/log"
)

// CheckAuth checks if user is authenticated and prompts to login if not
func CheckAuth() error {
	if _, err := GetToken(); err != nil {
		log.Debug("Failed to get token: %v", err)
		log.Error("You are not logged in. Please run 'osp auth login' to authenticate.")
		return ErrNotAuthenticated
	}
	return nil
}
