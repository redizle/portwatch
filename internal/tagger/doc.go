// Package tagger resolves human-readable service tags for port numbers.
//
// It ships with a built-in table of well-known IANA port-to-service mappings
// (ssh, http, https, mysql, redis, etc.) and supports user-defined overrides
// that take precedence over the defaults.
//
// Usage:
//
//	t := tagger.New(map[int]string{8080: "my-api"})
//	fmt.Println(t.Tag(8080))  // "my-api"
//	fmt.Println(t.Tag(22))    // "ssh"
//	fmt.Println(t.Tag(9999))  // "unknown"
//
// Tagger is safe for concurrent use.
package tagger
