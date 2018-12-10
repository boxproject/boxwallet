package util

import "log"

func CatchPanic() {
	if err := recover(); err != nil {
		log.Println("recover from ", err)
	}
}
