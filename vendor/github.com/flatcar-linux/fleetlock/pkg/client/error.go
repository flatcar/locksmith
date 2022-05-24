package client

import "fmt"

// Error is the type returned by FleetLock in case
// StatusCode is different from 200.
type Error struct {
	// Kind is the error type identifier.
	Kind string `json:"kind"`
	// Value is a human-friendly error description.
	Value string `json:"value"`
}

// String returns a high-level error formatted for human
// readability.
func (e *Error) String() string {
	return fmt.Sprintf("%s (%s)", e.Value, e.Kind)
}
