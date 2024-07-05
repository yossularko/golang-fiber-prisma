package lib

import (
	"errors"
	"fmt"
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("err cek deleted record:", r.(error).Error())
			return
		}
	}()

	var timeDeletedAt *time.Time
	if ok {
		timeDeletedAt = &value
	}

	if timeDeletedAt != nil {
		if timeDeletedAt.IsZero() {
			timeDeletedAt = nil
		}
	}

	if timeDeletedAt != nil {
		return errors.New("record not found")
	}

	return nil
}
