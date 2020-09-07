package errors

import (
	"errors"
)

func Catch(err error) {
	if err != nil {
		errors.New("strconv.Atoi() failed to convert int")
	}
}
