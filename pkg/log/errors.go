package log

import "errors"

var (
	RecordNotFoundError      = errors.New("record not found")
	DuplicateRecordError     = errors.New("duplicate record")
	InsufficientBalanceError = errors.New("balance is not sufficient")
)
