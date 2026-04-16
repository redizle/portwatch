package debounce_test

import (
	"fmt"
	"time"

	"portwatch/internal/debounce"
)

func ExampleDebouncer_Push() {
	results := make(chan string, 1)

	d := debounce.New(20*time.Millisecond, func(key string) {
		results <- key
	})

	// Rapid pushes — only the last one should fire.
	d.Push("3000")
	d.Push("3000")
	d.Push("3000")

	select {
	case key := <-results:
		fmt.Println("fired:", key)
	case <-time.After(200 * time.Millisecond):
		fmt.Println("timeout")
	}

	// Output:
	// fired: 3000
}
