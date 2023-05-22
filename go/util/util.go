package util

import "github.com/google/uuid"

func Getrandomstr(length int) string {
	id := uuid.New()
	str := id.String()[:length]
	return str
}
