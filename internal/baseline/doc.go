// Package baseline provides functionality for defining and comparing
// against an expected set of open ports on localhost.
//
// Ports added to the baseline are considered "normal". Any port that
// is active but not in the baseline is flagged as unexpected, which
// can be used to trigger alerts or log warnings.
//
// The baseline can be persisted to and loaded from a JSON file so that
// it survives daemon restarts.
package baseline
