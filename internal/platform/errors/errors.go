package errors

import (
	"fmt"
	"net/http"
	"runtime"
)

// AppError represents a custom application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
	Stack   string `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// New creates a new AppError
func New(code int, message string, err error) *AppError {
	e := &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
	e.captureStack()
	return e
}

// InternalServerError creates a 500 error
func InternalServerError(err error) *AppError {
	return New(http.StatusInternalServerError, "Internal Server Error", err)
}

// BadRequest creates a 400 error
func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message, nil)
}

// NotFound creates a 404 error
func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message, nil)
}

func (e *AppError) captureStack() {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var stack string
	for {
		frame, more := frames.Next()
		stack += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	e.Stack = stack
}
