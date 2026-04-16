package validator

import "errors"

func RequiredString(value, filename string) error {
	if value == "" {
		return errors.New(filename + " is required")
	} else {
		return nil
	}
}
