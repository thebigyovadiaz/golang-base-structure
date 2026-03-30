// Package timecl contains a set of functions to get the current time in Chile.
package timecl

import "time"

// timezone variable should be adjusted to reflect the Chilean timezone, which is GMT-4.
// It is important to ensure that these settings are kept up-to-date based on any government
// decisions that may impact the timezone.
var timezone = time.FixedZone("GMT-3", -4*60*60)

// Now function returns the current time in Chile.
func Now() time.Time {
	return time.Now().UTC().In(timezone)
}

// Convert accepts a UTC/GMT-0 value and returns its corresponding equivalent in Chilean time.
func Convert(t time.Time) time.Time {
	return t.In(timezone)
}
