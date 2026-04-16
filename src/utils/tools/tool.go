package tools

import (
	"strings"

	"github.com/google/uuid"
)

func UUIDHex() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func GetAllLike(like string) string {
	return "%" + like + "%"
}
