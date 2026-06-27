package qp

import (
	"errors"
	"fmt"
)

// Sentinel errors reported through Result.Error. Match them with errors.Is:
//
//	r := q.ParseInt("age", qp.Between(0, 120))
//	if errors.Is(r.Error, qp.ErrNotAllowed) {
//	    // the value parsed but fell outside the allowed range/set
//	}
var (
	// ErrInvalid means a present, non-empty value could not be parsed into the
	// target type (for example "abc" as an int).
	ErrInvalid = errors.New("qp: invalid value")

	// ErrNotAllowed means the value parsed but failed validation — it fell
	// outside the configured range and was not in the allowed set.
	ErrNotAllowed = errors.New("qp: value not allowed")

	// ErrEmptyElement means a slice element was empty under the ErrorOnEmpty
	// policy.
	ErrEmptyElement = errors.New("qp: empty element")
)

// ParseError carries the key and offending raw value alongside one of the
// sentinel errors. Retrieve it with errors.As:
//
//	var pe *qp.ParseError
//	if errors.As(r.Error, &pe) {
//	    log.Printf("bad %q value %q", pe.Key, pe.Value)
//	}
type ParseError struct {
	Key   string // the query parameter name
	Value string // the offending raw value (empty for ErrEmptyElement)
	Err   error  // the wrapped sentinel: ErrInvalid, ErrNotAllowed, ErrEmptyElement
}

// Error implements error.
func (e *ParseError) Error() string {
	switch {
	case errors.Is(e.Err, ErrEmptyElement):
		return fmt.Sprintf("qp: empty element for key %q", e.Key)
	case errors.Is(e.Err, ErrNotAllowed):
		return fmt.Sprintf("qp: value not allowed for key %q: %q", e.Key, e.Value)
	default:
		return fmt.Sprintf("qp: invalid value for key %q: %q", e.Key, e.Value)
	}
}

// Unwrap exposes the wrapped sentinel for errors.Is.
func (e *ParseError) Unwrap() error {
	return e.Err
}

// invalidErr builds an ErrInvalid ParseError.
func invalidErr(key, value string) error {
	return &ParseError{Key: key, Value: value, Err: ErrInvalid}
}

// notAllowedErr builds an ErrNotAllowed ParseError.
func notAllowedErr(key, value string) error {
	return &ParseError{Key: key, Value: value, Err: ErrNotAllowed}
}

// emptyElementErr builds an ErrEmptyElement ParseError.
func emptyElementErr(key string) error {
	return &ParseError{Key: key, Err: ErrEmptyElement}
}
