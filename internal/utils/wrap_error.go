package utils

import "fmt"

// WrapError wraps an error with a custom action message.
func WrapError(action string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("error %s: %w", action, err)
}
