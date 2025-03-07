package utils

import "fmt"

func AsError(err error, msg string, reason string) error {
	return fmt.Errorf("%s: %w: %s", msg, err, reason)
}
