package util

import (
	"errors"
	"fmt"
	"strconv"
)

func ConvertAndValidateID(str string) (uint64, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("'%s' is not an integer\n", str)
	} else if id <= 0 {
		return 0, errors.New("integer can only be greater than 0")
	}

	return id, nil
}
