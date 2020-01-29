package lib

import (
	"github.com/satori/go.uuid"
)

func GenUUID() (id string) {
	u1 := uuid.NewV4()
	id = u1.String()
	return
}
