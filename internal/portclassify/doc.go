// Package portclassify provides tier-based classification of TCP port numbers.
//
// Ports are divided into three standard tiers:
//
//   - system     (1–1023):     well-known ports assigned by IANA
//   - registered (1024–49151): registered application ports
//   - dynamic    (49152–65535): ephemeral / private ports
//
// Custom overrides can be applied per-port to override the default range
// logic, which is useful when a service runs on a non-standard port but
// should still be treated as a system-level concern.
//
// Example:
//
//	c := portclassify.New()
//	tier, err := c.Classify(8080) // => registered
package portclassify
