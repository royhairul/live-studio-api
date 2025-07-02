package paramhandler

import (
	"fmt"
	"strconv"
)

func ParseUintParam(param string) (uint, error) {
	idUint, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid parameter %s: %w", param, err)
	}

	return uint(idUint), nil
}
