package auth

import "errors"

// Claims represents the JWT claims structure
type Claims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Validate implements jwt.ClaimsValidator
func (c *Claims) Validate() error {
	if c.Username == "" {
		return errors.New("username is required")
	}
	return nil
}

func (c *Claims) GetUsername() string {
	return c.Username
}
