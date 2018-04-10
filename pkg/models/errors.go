package models

import (
	"errors"
)

// These are errors that are returned by various model functions

// FetchError is a base type that can typecast for specific APIs,
// It indicicates why an attempted db fetch failed.
type FetchError string

const (
	// FetchErrorNotFound means that the requested record does not exist
	FetchErrorNotFound FetchError = "FETCH_NOT_FOUND"
	// FetchErrorForbidden means that the record exists but that the user does not have access to it
	FetchErrorForbidden FetchError = "FETCH_FORBIDDEN"
)

// ErrFetchNotFound means that the requested record does not exist
var ErrFetchNotFound = errors.New("FETCH_NOT_FOUND")

// ErrFetchForbidden means that the record exists but that the user does not have access to it
var ErrFetchForbidden = errors.New("FETCH_FORBIDDEN")

// RecordNotFoundErrorString is the error string returned when no matching rows exist in the database
// This is ugly, but the best we can do with go's sql adapter
const RecordNotFoundErrorString = "sql: no rows in result set"
