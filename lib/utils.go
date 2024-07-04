package lib

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckDeletedRecord(value time.Time, ok bool) error {
	var timeDeletedAt *time.Time
	if ok {
		timeDeletedAt = &value
	}

	if timeDeletedAt.IsZero() {
		timeDeletedAt = nil
	}

	if timeDeletedAt != nil {
		return errors.New("record not found")
	}

	return nil
}
