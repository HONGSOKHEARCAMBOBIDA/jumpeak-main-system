package helper

import (
	"context"
	"errors"
	"fmt"
	"mysql/constant/apperror"

	"github.com/go-sql-driver/mysql"
)

func MapError(err error, action string) error {
	// err → error ដែលកើតឡើង
	// action → សកម្មភាពដែលកំពុងធ្វើ
	var mysqlErr *mysql.MySQLError
	// បង្កើត variable មួយដែលមាន type *mysql.MySQLError ប្រើសម្រាប់រកថា err ដែលយើងទទួលបាន គឺជា MySQL error ឬអត់
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		// errors.As(err, &mysqlErr) តើ err នេះអាចជា *mysql.MySQLError ដែរឬទេ
		// 1062 = Duplicate entry
		return apperror.New(apperror.CodeConflict, "already exists", err)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return apperror.New(apperror.CodeUnavailable, "request timed out, please try again", err)
	}
	// តើ error នេះកើតឡើងដោយសារ request timeout មែនទេ?
	// context.DeadlineExceeded មានន័យថា context បានដល់ deadline ដែលបានកំណត់
	return apperror.Internal(fmt.Sprintf("failed to %s academic", action), err)
}
