package product

import "errors"

var (
	ErrNoMatch          = errors.New("no matching product found")
	ErrLowConfidence    = errors.New("match confidence too low")
	ErrAlreadyMatched   = errors.New("product already matched")
	ErrInvalidCandidate = errors.New("invalid match candidate")
)
