package cachethq

import "errors"

// ErrComponentNotFound Error component not found
var ErrComponentNotFound = errors.New("component not found")

// ErrComponentGroupNotFound Error component group not found
var ErrComponentGroupNotFound = errors.New("component group not found")

// ErrStatusNotFound Error status not found
var ErrStatusNotFound = errors.New("status not found")
