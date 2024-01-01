package myerr

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/oklog/ulid"

	"gharabakloo/search/pkg/trace"
)

const (
	debugKey  = "APP_DEBUG"
	debugMode = "true"
)

var errInternal = errors.New(http.StatusText(http.StatusInternalServerError))

var errorCodes = map[error]int{
	errInternal: http.StatusInternalServerError,
}

var localDetails = map[string]string{
	errInternal.Error(): "خطای سرور",
}

type Error struct {
	Code        int    `json:"code"`
	Detail      string `json:"detail"`
	LocalDetail string `json:"local_detail"`
	TraceID     string `json:"trace_id"`
	Path        string `json:"path"`
	Trace       string `json:"trace"`
	Err         error  `json:"-"`
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	if os.Getenv(debugKey) == debugMode {
		return fmt.Sprintf("code:%d\ndetail:%s\nlocal_detail:%s\n"+
			"trace_id:%s\npath:%s\nstacke trace:\n%s\nerror:\n%v\n",
			e.Code, e.Detail, e.LocalDetail, e.TraceID, e.Path, e.Trace, e.Err)
	}
	return fmt.Sprintf("code:%d\ndetail:%s\nlocal_detail:%s\n"+
		"trace_id:%s\npath:%s\nerror:\n%v\n",
		e.Code, e.Detail, e.LocalDetail, e.TraceID, e.Path, e.Err)
}

func SetErrorInternal(internalError error) {
	errInternal = internalError
}

func SetErrorCodes(codesMap map[error]int) {
	errorCodes = codesMap
}

func SetLocalDetails(localDetailsMap map[string]string) {
	localDetails = localDetailsMap
}

func Errorf(errs ...error) error {
	if len(errs) == 1 && errs[0] == nil {
		return nil
	}

	var e *Error
	if len(errs) > 0 && errors.As(errs[0], &e) {
		return e
	}

	stack := trace.StackTrace()
	e = &Error{
		TraceID: generate(),
		Path:    stack.Path,
		Trace:   stack.String(),
	}

	switch len(errs) {
	case 0:
		e.Err = errInternal
		e.Detail = errInternal.Error()
		e.LocalDetail = localDetails[e.Detail]
		e.Code = errorCodes[errInternal]
	case 1:
		e.Err = errs[0]
		e.Detail = errs[0].Error()
		e.LocalDetail = localDetails[e.Detail]
		e.Code = errorCodes[errs[0]]
		if e.Code == 0 {
			e.LocalDetail = localDetails[errInternal.Error()]
			e.Code = errorCodes[errInternal]
		}
	default:
		e.Err = errs[1]
		e.Detail = errs[0].Error()
		e.LocalDetail = localDetails[e.Detail]
		e.Code = errorCodes[errs[0]]
		if e.Code == 0 {
			e.LocalDetail = localDetails[errInternal.Error()]
			e.Code = errorCodes[errInternal]
		}
	}
	return e
}

func generate() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0) // nolint: gosec
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
