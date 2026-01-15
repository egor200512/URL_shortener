package pkg

import (
	"fmt"
	"net/mail"
	"strings"
)

const maxLen = 50

func ValidateEmailPassword(email, password string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("email is empty")
	}
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("password is empty")
	}

	if len(email) > maxLen {
		return fmt.Errorf("email too long (max %d chars)", maxLen)
	}
	if len(password) > maxLen {
		return fmt.Errorf("password too long (max %d chars)", maxLen)
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format: %s", err.Error())
	}

	return nil
}
