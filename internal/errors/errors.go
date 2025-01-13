package errors

import (
	err "errors"
	"fmt"

	"github.com/liqmix/slaptrax/internal/logger"
)

type AppError string

const (
	UNKNOWN_LOCALE AppError = "Unknown locale %s"
)

func Raise(e AppError, args ...interface{}) error {
	m := fmt.Sprintf(string(e), args...)
	logger.Error(m)
	return err.New(m)
}
