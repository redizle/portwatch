// Package portname resolves port numbers to human-readable service names.
//
// It combines a built-in well-known port table with user-defined overrides,
// allowing callers to get friendly names like "http" or "postgres" for
// common ports while still supporting custom labels for private services.
//
// Usage:
//
//	 resolver := portname.New(nil)
//	 name := resolver.Resolve(443) // returns "https"
//
//	 resolver.SetOverride(9000, "my-api")
//	 name = resolver.Resolve(9000) // returns "my-api"
package portname
