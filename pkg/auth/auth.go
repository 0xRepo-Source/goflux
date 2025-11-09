package auth

// Package auth contains authentication helpers and interfaces for goflux.

// Credential represents an authentication credential (placeholder).
type Credential struct {
	Type  string // e.g. "ssh", "token", "jwt"
	Value string
}

// Validate checks the credential (placeholder).
func (c *Credential) Validate() bool {
	return c != nil && c.Value != ""
}
