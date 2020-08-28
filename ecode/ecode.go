package ecode

import (
	"strconv"
)

// Codes ecode error interface which has a code & message.
type Codes interface {
	// sometimes Error return Code in string form
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
}

// Add init you self service error code
func Add(code int) Error {
	return Error(code)
}

func New(code int, msg string) Error {
	errorMap.Store(code, msg)
	return Error(code)
}

// A Code is an int error code spec.
type Error int

func (e Error) Error() string {
	if msg, ok := errorMap.Load(e.Code()); ok {
		return msg.(string)
	}
	return strconv.Itoa(int(e))
}

// Code return error code
func (e Error) Code() int { return int(e) }

// Message return error message
func (e Error) Message() string {
	if msg, ok := errorMap.Load(e.Code()); ok {
		return msg.(string)
	}
	return e.Error()
}

// analyse error info
func AnalyseError(err error) Codes {
	if err == nil {
		return OK
	}
	if codes, ok := err.(Error); ok {
		return codes
	}
	return errStringToError(err.Error())
}

func errStringToError(e string) Error {
	if e == "" {
		return OK
	}
	i, err := strconv.Atoi(e)
	if err != nil {
		return New(-1, e)
	}
	return Error(i)
}
