package utils

import "github.com/google/uuid"

func GenerateUniqueToken() string {
	return uuid.NewString()
}
