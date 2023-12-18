package authdb

import (
	"fmt"
	"time"
)

// Errors.

var ErrNilState = fmt.Errorf("DbState cannot be nil")
var ErrClosed = fmt.Errorf("Database connection is closed")

// OP_TIMEOUT is the maximum allowed duration allotted to a single
// database operation.
const OP_TIMEOUT = 300 * time.Millisecond
