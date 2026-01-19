package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"
)

// PasswordPolicy defines the password requirements
type PasswordPolicy struct {
	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireDigit     bool
	RequireSpecial   bool
	CheckPwned       bool
}

// DefaultPasswordPolicy returns the default strict password policy
func DefaultPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:        12,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireDigit:     true,
		RequireSpecial:   true,
		CheckPwned:       true,
	}
}

// PasswordValidationError contains details about password validation failures
type PasswordValidationError struct {
	Errors []string
}

func (e *PasswordValidationError) Error() string {
	return strings.Join(e.Errors, "; ")
}

// ValidatePassword validates a password against the policy
func (p *PasswordPolicy) ValidatePassword(password string) error {
	var errors []string

	// Check minimum length
	if len(password) < p.MinLength {
		errors = append(errors, fmt.Sprintf("Password must be at least %d characters long", p.MinLength))
	}

	// Check for uppercase
	if p.RequireUppercase && !containsUppercase(password) {
		errors = append(errors, "Password must contain at least one uppercase letter")
	}

	// Check for lowercase
	if p.RequireLowercase && !containsLowercase(password) {
		errors = append(errors, "Password must contain at least one lowercase letter")
	}

	// Check for digit
	if p.RequireDigit && !containsDigit(password) {
		errors = append(errors, "Password must contain at least one digit")
	}

	// Check for special character
	if p.RequireSpecial && !containsSpecial(password) {
		errors = append(errors, "Password must contain at least one special character (!@#$%^&*()-_=+[]{}|;:',.<>?/)")
	}

	if len(errors) > 0 {
		return &PasswordValidationError{Errors: errors}
	}

	// Check against haveibeenpwned (only if basic validation passes)
	if p.CheckPwned {
		pwned, err := IsPasswordPwned(password)
		if err != nil {
			// Log error but don't block registration if API is unavailable
			// In production, you might want to handle this differently
		} else if pwned {
			return &PasswordValidationError{
				Errors: []string{"This password has been found in data breaches and cannot be used. Please choose a different password."},
			}
		}
	}

	return nil
}

// containsUppercase checks if string contains uppercase letter
func containsUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// containsLowercase checks if string contains lowercase letter
func containsLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// containsDigit checks if string contains digit
func containsDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// containsSpecial checks if string contains special character
func containsSpecial(s string) bool {
	specialChars := "!@#$%^&*()-_=+[]{}|;:',.<>?/`~\"\\"
	for _, r := range s {
		if strings.ContainsRune(specialChars, r) {
			return true
		}
	}
	return false
}

// IsPasswordPwned checks if a password has been found in data breaches
// using the haveibeenpwned API with k-anonymity (only sends first 5 chars of SHA1 hash)
func IsPasswordPwned(password string) (bool, error) {
	// Hash the password with SHA1
	hash := sha1.New()
	hash.Write([]byte(password))
	hashBytes := hash.Sum(nil)
	hashStr := strings.ToUpper(hex.EncodeToString(hashBytes))

	// Split hash: first 5 characters for API, rest for comparison
	prefix := hashStr[:5]
	suffix := hashStr[5:]

	// Query haveibeenpwned API
	url := fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", prefix)
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to check password: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("haveibeenpwned API returned status %d", resp.StatusCode)
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	// Check if our hash suffix is in the response
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts) >= 1 && parts[0] == suffix {
			return true, nil // Password found in breach database
		}
	}

	return false, nil // Password not found in breach database
}

// ValidatePasswordWithPolicy is a convenience function using the default policy
func ValidatePasswordWithPolicy(password string) error {
	policy := DefaultPasswordPolicy()
	return policy.ValidatePassword(password)
}
