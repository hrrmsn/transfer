package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Min(ints []int64) (int64, error) {
	if len(ints) == 0 {
		return -1, fmt.Errorf("empty input array")
	}

	var result = ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] < result {
			result = ints[i]
		}
	}
	return result, nil
}

func HandleError(textError string, statusCode int, err error, w http.ResponseWriter) {
	log.Fatalf("%s: %s\n", textError, err.Error())

	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error": "` + strings.ToLower(textError) + `"}`))
}
