package model

import "errors"

type err error

type Error struct {
	Code        string
	Description string
	err
}

func (e *Error) Unwrap() error { return e.err }

var ErrorNotFound = Error{
	"e404",
	"Not Found exception",
	errors.New("not found exception"),
}

var ErrorUnexpected = Error{
	"e500",
	"unexpected exception",
	errors.New("unexpected exception"),
}

var ErrorParsing = Error{
	"e500",
	"Parsing exception",
	errors.New("parsing exception"),
}

var ErrorMistmatchType = Error{
	"e500",
	"mismatch value type",
	errors.New("mismatch value type"),
}
