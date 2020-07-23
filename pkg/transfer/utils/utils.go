package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func Min(ints []int64) (int64, error) {
	if len(ints) == 0 {
		return -1, fmt.Errorf("Empty input array")
	}

	var result = ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] < result {
			result = ints[i]
		}
	}
	return result, nil
}

func HandleError(w http.ResponseWriter, err error, statusCode int) {
	if err == nil {
		err = fmt.Errorf("Unknown error")
	}
	http.Error(w, `{"error": "`+strings.ToLower(err.Error())+`"}`, statusCode)
}

func WrapError(text string, err error) error {
	return fmt.Errorf("%s: %s", text, err.Error())
}
