package utils

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const defaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

func GenerateTempPassword() string {
	length := 12
	if envLength := os.Getenv("TEMP_PASSWORD_LENGTH"); envLength != "" {
		fmt.Sscanf(envLength, "%d", &length)
	}

	charset := os.Getenv("TEMP_PASSWORD_CHARSET")
	if charset == "" {
		charset = defaultCharset
	}

	rand.Seed(time.Now().UnixNano())
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}
