package validators

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidaString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n > maxLength || n < minLength {
		return fmt.Errorf("Must containen from %d to %d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	err := ValidaString(value, 3, 100)
	if err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("Username only can contain lower case letters, diggits or underscore")
	}
	return nil
}

func ValidateFullname(value string) error {
	err := ValidaString(value, 3, 100)
	if err != nil {
		return err
	}
	if !isValidFullname(value) {
		return fmt.Errorf("Fullname only can contain letters or spaces")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidaString(value, 8, 256)
}

func ValidateEmail(value string) error {
	err := ValidaString(value, 3, 200)
	if err != nil {
		return err
	}
	_, err = mail.ParseAddress(value)
	if err != nil {
		return fmt.Errorf("%s is not a valid email address", value)
	}
	return nil
}
