package lib

import (
	"encoding/json"
	"strings"
)

type HttpValidationError map[string]map[string]string

func FormatError(err string) HttpValidationError {
	errMap := make(map[string]string)

	err = strings.TrimSuffix(err, ".")

	for e := range strings.SplitSeq(err, "; ") {
		key, value, _ := strings.Cut(e, ": ")

		errMap[key] = key + " " + value
	}

	return map[string]map[string]string{
		"errors": errMap,
	}
}

func (e HttpValidationError) Error() string {
	res, _ := json.Marshal(e)
	return string(res)
}
