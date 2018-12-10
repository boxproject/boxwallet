package util

import (
	"log"
	"testing"
	"time"
)

func TestTimingWheel(t *testing.T) {
	w := NewTimingWheel(100*time.Millisecond, 10)

	go func() {
		for i := 0; ; i++ {
			select {
			case <-w.After(100 * time.Millisecond):
				log.Println(1, i)
			}
		}
	}()

	go func() {
		for i := 0; ; i++ {
			select {
			case <-w.After(200 * time.Millisecond):
				log.Println(2, i)
			}
		}
	}()
	go func() {
		for i := 0; ; i++ {
			select {
			case <-w.After(400 * time.Millisecond):
				log.Println(3, i)
			}
		}
	}()
	select {}
}
