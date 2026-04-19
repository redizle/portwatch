package portrelay_test

import (
	"fmt"

	"portwatch/internal/portrelay"
)

func ExampleRelay_Dispatch() {
	r := portrelay.New()

	_ = r.Register(portrelay.Target{
		Name: "stdout",
		Handler: func(port int, event string) error {
			fmt.Printf("port=%d event=%s\n", port, event)
			return nil
		},
	})

	_ = r.Dispatch(8080, "open")
	// Output:
	// port=8080 event=open
}

func ExampleRelay_Len() {
	r := portrelay.New()
	fmt.Println(r.Len())

	_ = r.Register(portrelay.Target{
		Name:    "a",
		Handler: func(int, string) error { return nil },
	})
	fmt.Println(r.Len())
	// Output:
	// 0
	// 1
}
