package helper

import (
	"context"
	"errors"
	"fmt"
	"mysql/constant/apperror"

	"github.com/go-sql-driver/mysql"
)

func MapError(err error, action string) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return apperror.New(apperror.CodeConflict, "already exists", err)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return apperror.New(apperror.CodeUnavailable, "request timed out, please try again", err)
	}

	return apperror.Internal(fmt.Sprintf("failed to %s academic", action), err)
}
