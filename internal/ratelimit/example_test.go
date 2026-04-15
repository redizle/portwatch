package ratelimit_test

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func ExampleLimiter_Allow() {
	l := ratelimit.New(1 * time.Hour)

	if l.Allow("port:8080") {
		fmt.Println("alert sent")
	}

	if !l.Allow("port:8080") {
		fmt.Println("alert suppressed")
	}

	// Output:
	// alert sent
	// alert suppressed
}

func ExampleLimiter_Reset() {
	l := ratelimit.New(1 * time.Hour)
	l.Allow("port:443")
	l.Reset("port:443")

	if l.Allow("port:443") {
		fmt.Println("allowed after reset")
	}

	// Output:
	// allowed after reset
}
