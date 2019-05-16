package util

import (
	"errors"
	"fmt"
	"strconv"
)

// ConvertAndValidateID is a helper function that converts a string into a Job
// ID. A Job ID must be greater than 0.
func ConvertAndValidateID(str string) (uint64, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("'%q' is not an integer", str)
	} else if id <= 0 {
		return 0, errors.New("integer can only be greater than 0")
	}

	return id, nil
}
