package utils

import "fmt"

func AsError(err error, reason string) error {
	return fmt.Errorf("%s: %w", reason, err)
}
