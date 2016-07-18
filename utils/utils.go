package utils

import (
	"strconv"
)

func ParseIdInt64FromString(s string) (int64, error) {
	result, err := strconv.ParseInt(s, 10, 64)

	return result, err
}
