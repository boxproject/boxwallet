package util

import "github.com/satori/go.uuid"

func GetUUid() string {
	for {
		uu, err := uuid.NewV4()
		if err == nil {
			return uu.String()
		}
	}
}
