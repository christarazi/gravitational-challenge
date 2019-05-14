package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// ConvertAndValidateID is a helper function that converts a string into a Job
// ID. A Job ID must be greater than 0.
func ConvertAndValidateID(str string) (uint64, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("'%s' is not an integer\n", str)
	} else if id <= 0 {
		return 0, errors.New("integer can only be greater than 0")
	}

	return id, nil
}

// CheckHTTPStatusCode is a helper function to check the status code and drains
// the body of the response if the status code is not OK. An error is
// returned if so.
func CheckHTTPStatusCode(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error reading body of response: %v", err)
		}

		return fmt.Errorf("Server returned %d: %v", resp.StatusCode, string(body))
	}

	return nil
}
