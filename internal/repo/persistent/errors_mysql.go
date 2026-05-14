// Package persistent provides MySQL-specific error handlers.
package persistent

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// MySQL error codes.
const (
	ErrMySqlDuplicateEntry = 1062
	ErrMySqlForeignKey     = 1452
)

// IsDuplicateEntry returns true if the error is a duplicate entry violation.
func IsDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == ErrMySqlDuplicateEntry
	}
	return false
}
