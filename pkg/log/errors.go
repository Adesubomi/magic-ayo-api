package log

import "errors"

type GameError error

var (
	RecordNotFoundError      = errors.New("record not found")
	DuplicateRecordError     = errors.New("duplicate record")
	InsufficientBalanceError = errors.New("balance is not sufficient")
	UnknownError             = errors.New("unknown error")
)

var (
	InvalidPotNumber GameError = errors.New("invalid pot number")
	EmptyPot         GameError = errors.New("pot is empty")
)

// Abort game
// Make move
// Generate lnurl
