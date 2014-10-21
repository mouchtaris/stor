package util

import (
    "fmt"
    "os"
)

//
//
type ErrorHandler interface {
    ErrorsChannel () chan<- error
}

//
//
type errorHandler struct {
    name        string
    errors      chan error
}

//
//
func handleErrors (name string, errors <-chan error) {
    for err := range errors {
        fmt.Fprintf(os.Stderr, "[%s]: %s\n", name, err)
    }
}

//
//
func NewErrorHandler (name string) ErrorHandler {
    return &errorHandler {
        name:       name,
        errors:     make(chan error, 20),
    }
}

//
//
func (eh *errorHandler) ErrorsChannel () chan<- error {
    return eh.errors
}
