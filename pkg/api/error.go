package api

import (
	"errors"
	"fmt"
	"strings"
)

// utility function

// for Go 1.13 error system
func ExtractErrorMessageChain(err error) string {
	var b strings.Builder
	for e := err; e != nil; e = errors.Unwrap(e) {
		b.WriteString(e.Error())
		b.WriteString("\n")
	}
	return b.String()
}

// for Go 1.13 error system & golang.org/x/xerrors
func WrapError(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// for all error frameworks
func GetErrorMessageTitle(err error) string {
	return strings.SplitN(err.Error(), ":", 2)[0]
}

// compare error using first part of message before colon, for example: "<first>: <rest>"
func IsErrorEqual(e1 error, e2 error) bool {
	if e1 == nil || e2 == nil {
		return e1 == e2
	} else {
		// CAUTION: error type is useless (in most cases)
		// return reflect.TypeOf(e1) == reflect.TypeOf(e2) && GetErrorMessageTitle(e1) == GetErrorMessageTitle(e2)
		return GetErrorMessageTitle(e1) == GetErrorMessageTitle(e2)
	}
}
