package throttle_test

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func ExampleThrottle_Allow() {
	th := throttle.New(50 * time.Millisecond)

	if th.Allow() {
		fmt.Println("scan started")
	}

	// immediate second call is throttled
	if !th.Allow() {
		fmt.Println("scan skipped")
	}

	time.Sleep(60 * time.Millisecond)

	if th.Allow() {
		fmt.Println("scan started again")
	}

	// Output:
	// scan started
	// scan skipped
	// scan started again
}

func ExampleThrottle_Reset() {
	th := throttle.New(1 * time.Hour)
	th.Allow() // consume the first slot

	th.Reset()

	if th.Allow() {
		fmt.Println("allowed after reset")
	}

	// Output:
	// allowed after reset
}
