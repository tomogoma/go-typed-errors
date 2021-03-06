package errors

import (
	"fmt"
	"net/http"
)

type IsAuthErrChecker interface {
	IsAuthError(error) bool
	IsUnauthorizedError(error) bool
	IsForbiddenError(error) bool
}

type IsNotFoundErrChecker interface {
	IsNotFoundError(error) bool
}

type IsNotImplErrChecker interface {
	IsNotImplementedError(error) bool
}

type IsClErrChecker interface {
	IsClientError(error) bool
}

type IsRetryableErrChecker interface {
	IsRetryableError(error) bool
}

type IsConflictErrChecker interface {
	IsConflictError(error) bool
}

type IsPreconditionFailedErrChecker interface {
	IsPreconditionFailedError(error) bool
}

type AllErrChecker interface {
	IsAuthErrChecker
	IsNotFoundErrChecker
	IsNotImplErrChecker
	IsClErrChecker
	IsRetryableErrChecker
	IsConflictErrChecker
	IsPreconditionFailedErrChecker
}

type ToHTTPResponser interface {
	ToHTTPResponse(err error, w http.ResponseWriter) (int, bool)
}

// ErrToHTTP implements ToHTTPResponser interface. It can be embedded in a struct
// to give said custom struct the ToHTTPResponse method. e.g:
//  type Custom struct {
//      ...
//      errors.ErrToHTTP
//  }
type ErrToHTTP struct {
}

// ToHTTPResponse attempts to run Error.ToHTTPResponse(w) returning
// the result if the call was successful, -1 and false otherwise.
func (e ErrToHTTP) ToHTTPResponse(err error, w http.ResponseWriter) (int, bool) {
	if err, ok := err.(Error); ok {
		return err.ToHTTPResponse(w)
	}
	return -1, false
}

// Error implements the Error interface and helps distinguish whether an error
// is a client error or an auth error.
type Error struct {
	IsAuthErr               bool
	IsUnauthorizedErr       bool
	IsForbiddenErr          bool
	IsClErr                 bool
	IsNotFoundErr           bool
	IsNotImplementedErr     bool
	IsRetryableErr          bool
	IsConflictErr           bool
	IsPreconditionFailedErr bool
	Data                    interface{}
	HttpMsg                 string
}

// Error returns the error message of the error (without the distinguishing flags
// such as client error).
func (e Error) Error() string {
	return fmt.Sprint(e.Data)
}

// Client returns true if this is a client error.
func (e Error) Client() bool {
	return e.IsClErr
}

// ToHTTPResp writes the content of the error to w while setting the HTTP status
// code to match the type of error received. Returns the HTTP status code
// assigned and true if error was written, -1 and false otherwise.
func (e Error) ToHTTPResponse(w http.ResponseWriter) (int, bool) {

	msg := e.HttpMsg
	if msg == "" {
		msg = e.Error()
	}

	if e.IsAuthErr || e.IsForbiddenErr || e.IsUnauthorizedErr {
		if e.IsForbiddenErr {
			http.Error(w, msg, http.StatusForbidden)
			return http.StatusForbidden, true
		}
		http.Error(w, msg, http.StatusUnauthorized)
		return http.StatusUnauthorized, true
	}

	if e.IsClErr {
		http.Error(w, msg, http.StatusBadRequest)
		return http.StatusBadRequest, true
	}

	if e.IsNotFoundErr {
		http.Error(w, msg, http.StatusNotFound)
		return http.StatusNotFound, true
	}

	if e.IsNotImplementedErr {
		http.Error(w, msg, http.StatusNotImplemented)
		return http.StatusNotImplemented, true
	}

	if e.IsRetryableErr {
		http.Error(w, msg, http.StatusServiceUnavailable)
		return http.StatusServiceUnavailable, true
	}

	if e.IsConflictErr {
		http.Error(w, msg, http.StatusConflict)
		return http.StatusConflict, true
	}

	if e.IsPreconditionFailedErr {
		http.Error(w, msg, http.StatusPreconditionFailed)
		return http.StatusPreconditionFailed, true
	}

	return -1, false
}

// NotImplemented returns true if the functionality requested is not implemented.
func (e Error) NotImplemented() bool {
	return e.IsNotImplementedErr
}

// Auth returns true if this is an auth error.
func (e Error) Auth() bool {
	return e.IsAuthErr || e.IsUnauthorizedErr || e.IsForbiddenErr
}

// Unauthorized returns true if this is an Unauthorized error.
func (e Error) Unauthorized() bool {
	return e.IsUnauthorizedErr
}

// Forbidden returns true if this is a Forbidden error.
func (e Error) Forbidden() bool {
	return e.IsForbiddenErr
}

// NotFound returns true if this is error denotes that a resource
// being fetched was not found.
func (e Error) NotFound() bool {
	return e.IsNotFoundErr
}

// Retryable returns true if this error is not permanent and should
// be retried
func (e Error) Retryable() bool {
	return e.IsRetryableErr
}

// Conflict returns true if this error denotes a conflict in resources a la
// HTTPs 409 error
func (e Error) Conflict() bool {
	return e.IsConflictErr
}

// PreconditionFailed returns true if this error denotes a
// precondition failure in resources a la HTTPs 412 error
func (e Error) PreconditionFailed() bool {
	return e.IsPreconditionFailedErr
}

// New creates a new error.
func New(data interface{}) Error {
	return Error{Data: data}
}

// Newf creates a new error with fmt.Printf formatting.
func Newf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return New(data)
}

// NewWithHttp creates a new error containing a http specific
// error message.
func NewWithHttp(httpMsg string, data interface{}) Error {
	return Error{Data: data, HttpMsg: httpMsg}
}

// NewWithHttp creates a new error containing a http specific
// error message.
func NewWithHttpf(httpMsg string, format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewWithHttp(httpMsg, data)
}

// NewClient creates a new client error.
func NewClient(data interface{}) Error {
	return Error{Data: data, IsClErr: true}
}

// NewClientf creates a new client error with fmt.Printf style formatting.
func NewClientf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewClient(data)
}

// NewClientWithHttp creates a new error containing a http specific
// error message.
func NewClientWithHttp(httpMsg string, data interface{}) Error {
	return Error{Data: data, HttpMsg: httpMsg, IsClErr: true}
}

// NewClientWithHttp creates a new error containing a http specific
// error message.
func NewClientWithHttpf(httpMsg string, format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewClientWithHttp(httpMsg, data)
}

// NewNotImplemented creates a new not implemented error.
func NewNotImplemented() Error {
	return Error{IsNotImplementedErr: true, Data: "not implemented"}
}

// NewNotImplementedf creates a new not implemented error with fmt.Printf
// style formatting.
func NewNotImplementedf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return Error{Data: data, IsNotImplementedErr: true}
}

// NewNotImplementedWithHttp creates a new error containing a http specific
// error message.
func NewNotImplementedWithHttp(httpMsg string, data interface{}) Error {
	return Error{Data: data, HttpMsg: httpMsg, IsNotImplementedErr: true}
}

// NewNotImplementedWithHttp creates a new error containing a http specific
// error message.
func NewNotImplementedWithHttpf(httpMsg string, format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewNotImplementedWithHttp(httpMsg, data)
}

// NewAuth creates a new auth error. It is not specific to the type of auth
// error. Use NewForbidden(string) or NewUnauthorized(string) to establish
// a more specific Auth error.
func NewAuth(data interface{}) Error {
	return Error{Data: data, IsAuthErr: true}
}

// NewAuthf creates a new auth error with fmt.Printf style formatting.
func NewAuthf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewAuth(data)
}

// NewAuthWithHttp creates a new error containing a http specific
// error message.
func NewAuthWithHttp(httpMsg string, data interface{}) Error {
	return Error{Data: data, HttpMsg: httpMsg, IsAuthErr: true}
}

// NewWithHttp creates a new error containing a http specific
// error message.
func NewAuthWithHttpf(httpMsg string, format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewAuthWithHttp(httpMsg, data)
}

// NewForbidden creates a new forbidden auth error a la 403 (http.StatusForbidden) error.
// This will also resolve as an Auth error.
func NewForbidden(data interface{}) Error {
	return Error{Data: data, IsAuthErr: true, IsForbiddenErr: true}
}

// NewForbiddenf creates a new forbidden auth error with fmt.Printf style formatting.
// This will also resolve as an Auth error.
func NewForbiddenf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewForbidden(data)
}

// NewForbiddentWithHttp creates a new error containing a http specific
// error message.
func NewForbiddentWithHttp(httpMsg string, data interface{}) Error {
	return Error{Data: data, HttpMsg: httpMsg, IsForbiddenErr: true}
}

// NewForbiddentWithHttp creates a new error containing a http specific
// error message.
func NewForbiddentWithHttpf(httpMsg string, format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewForbiddentWithHttp(httpMsg, data)
}

// NewUnauthorized creates a new unauthorized auth error a la 401 (http.StatusUnauthorized) error.
// This will also resolve as an Auth error.
func NewUnauthorized(data interface{}) Error {
	return Error{Data: data, IsAuthErr: true, IsUnauthorizedErr: true}
}

// NewUnauthorizedf creates a new unauthorized auth error with fmt.Printf style formatting.
// This will also resolve as an Auth error.
func NewUnauthorizedf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewUnauthorized(data)
}

// NewUnauthorizedWithHttp creates a new error containing a http specific
// error message.
func NewUnauthorizedWithHttp(httpMsg string, data interface{}) Error {
	return Error{Data: data, HttpMsg: httpMsg, IsUnauthorizedErr: true}
}

// NewWithHttp creates a new error containing a http specific
// error message.
func NewUnauthorizedWithHttpf(httpMsg string, format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewUnauthorizedWithHttp(httpMsg, data)
}

// NewNotFound creates a new not found error.
func NewNotFound(data interface{}) Error {
	return Error{Data: data, IsNotFoundErr: true}
}

// NewNotFoundf creates a new not found error with fmt.Printf style formatting.
func NewNotFoundf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewNotFound(data)
}

// NewNotFoundWithHttp creates a new error containing a http specific
// error message.
func NewNotFoundWithHttp(httpMsg string, data interface{}) Error {
	return Error{Data: data, HttpMsg: httpMsg, IsNotFoundErr: true}
}

// NewNotFoundWithHttp creates a new error containing a http specific
// error message.
func NewNotFoundWithHttpf(httpMsg string, format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewNotFoundWithHttp(httpMsg, data)
}

// NewRetryable creates a new retryable error.
func NewRetryable(data interface{}) Error {
	return Error{Data: data, IsRetryableErr: true}
}

// NewRetryablef creates a new retryable error with fmt.Printf style formatting.
func NewRetryablef(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewRetryable(data)
}

// NewRetryableWithHttp creates a new error containing a http specific
// error message.
func NewRetryableWithHttp(httpMsg string, data interface{}) Error {
	return Error{Data: data, HttpMsg: httpMsg, IsRetryableErr: true}
}

// NewRetryableWithHttp creates a new error containing a http specific
// error message.
func NewRetryableWithHttpf(httpMsg string, format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewRetryableWithHttp(httpMsg, data)
}

// NewConflict creates a new Conflict error.
func NewConflict(data interface{}) Error {
	return Error{Data: data, IsConflictErr: true}
}

// NewConflictf creates a new Conflict error with fmt.Printf style formatting.
func NewConflictf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewConflict(data)
}

// NewConflictWithHttp creates a new error containing a http specific
// error message.
func NewConflictWithHttp(httpMsg string, data interface{}) Error {
	return Error{Data: data, HttpMsg: httpMsg, IsConflictErr: true}
}

// NewConflictWithHttp creates a new error containing a http specific
// error message.
func NewConflictWithHttpf(httpMsg string, format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewConflictWithHttp(httpMsg, data)
}

// NewPreconditionFailed creates a new PreconditionFailed error.
func NewPreconditionFailed(data interface{}) Error {
	return Error{Data: data, IsPreconditionFailedErr: true}
}

// NewPreconditionFailedf creates a new PreconditionFailed error with fmt.Printf style formatting.
func NewPreconditionFailedf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewPreconditionFailed(data)
}

// NewPreconditionFailedWithHttp creates a new error containing a http specific
// error message.
func NewPreconditionFailedWithHttp(httpMsg string, data interface{}) Error {
	return Error{Data: data, HttpMsg: httpMsg, IsPreconditionFailedErr: true}
}

// NewPreconditionFailedWithHttp creates a new error containing a http specific
// error message.
func NewPreconditionFailedWithHttpf(httpMsg string, format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewPreconditionFailedWithHttp(httpMsg, data)
}

// ClErrCheck implements the ClErrChecker interface. It can be embedded in a custom struct to
// give the custom struct the extra method IsClientError(err error). e.g:
//  type Custom struct {
//      ...
//      errors.ClErrCheck
//  }
type ClErrCheck struct {
}

// IsClientError returns true if the supplied error is a client error, false otherwise.
func (c *ClErrCheck) IsClientError(err error) bool {
	errC, ok := err.(Error)
	return ok && errC.Client()
}

// NotImplErrCheck implements the NotImplErrChecker interface. It can be embedded in a custom struct to
// give the custom struct the extra method IsNotFoundError(err error). e.g:
//  type Custom struct {
//      ...
//      errors.NotImplErrCheck
//  }
type NotImplErrCheck struct {
}

// IsNotImplementedError returns true if the supplied error is a client error, false otherwise.
func (c *NotImplErrCheck) IsNotImplementedError(err error) bool {
	errC, ok := err.(Error)
	return ok && errC.NotImplemented()
}

// AuthErrCheck implements the AuthErrChecker interface. It can be embedded in a custom struct to
// give the custom struct the extra methods IsAuthError(err error),
// IsForbiddenError(err error) and IsUnauthorizedErr(err error) . e.g:
//  type Custom struct {
//      ...
//      errors.AuthErrCheck
//  }
type AuthErrCheck struct {
}

// IsAuthError returns true if the supplied error is an
// authentication/authorization error, false otherwise.
func (c *AuthErrCheck) IsAuthError(err error) bool {
	errC, ok := err.(Error)
	return ok && errC.Auth()
}

// IsAuthError returns true if the supplied error is an
// authentication/authorization error, false otherwise.
func (c *AuthErrCheck) IsForbiddenError(err error) bool {
	errC, ok := err.(Error)
	return ok && errC.Forbidden()
}

// IsAuthError returns true if the supplied error is an
// authentication/authorization error, false otherwise.
func (c *AuthErrCheck) IsUnauthorizedError(err error) bool {
	errC, ok := err.(Error)
	return ok && errC.Unauthorized()
}

// NotFoundErrCheck implements the NotFoundErrChecker interface. It can be
// embedded in a custom struct to give the custom struct the extra
// method IsNotFoundError(err error). e.g:
//  type Custom struct {
//      ...
//      errors.NotFoundErrCheck
//  }
type NotFoundErrCheck struct {
}

// IsNotFoundError returns true if the supplied error is an not found error, false otherwise.
func (c *NotFoundErrCheck) IsNotFoundError(err error) bool {
	errC, ok := err.(Error)
	return ok && errC.NotFound()
}

// RetryableErrCheck implements the RetryableErrChecker interface. It can be embedded in a custom struct to
// give the custom struct the extra method IsRetryableError(err error). e.g:
//  type Custom struct {
//      ...
//      errors.NotFoundErrCheck
//  }
type RetryableErrCheck struct {
}

// IsRetryableError returns true if the supplied error retryable, false otherwise.
func (c *RetryableErrCheck) IsRetryableError(err error) bool {
	errC, ok := err.(Error)
	return ok && errC.Retryable()
}

// ConflictErrCheck implements the ConflictErrChecker interface. It can be embedded in a custom struct to
// give the custom struct the extra method IsConflictError(err error). e.g:
//  type Custom struct {
//      ...
//      errors.ConflictErrCheck
//  }
type ConflictErrCheck struct {
}

// IsConflictError returns true if the supplied error is a Conflict error, false otherwise.
func (c *ConflictErrCheck) IsConflictError(err error) bool {
	errC, ok := err.(Error)
	return ok && errC.Conflict()
}

// PreconditionFailedErrCheck implements the PreconditionFailedErrChecker interface.
// It can be embedded in a custom struct to give the custom struct the extra method
// IsPreconditionFailedError(err error). e.g:
//  type Custom struct {
//      ...
//      errors.PreconditionFailedErrCheck
//  }
type PreconditionFailedErrCheck struct {
}

// IsConflictError returns true if the supplied error is a Conflict error, false otherwise.
func (c *PreconditionFailedErrCheck) IsPreconditionFailedError(err error) bool {
	errC, ok := err.(Error)
	return ok && errC.PreconditionFailed()
}

// AllErrCheck implements the AllErrChecker interface. It can be embedded in a custom struct to
// give said custom struct the extra Is...Error(err error) methods. e.g:
//  type Custom struct {
//      ...
//      errors.AllErrCheck
//  }
type AllErrCheck struct {
	AuthErrCheck
	NotFoundErrCheck
	NotImplErrCheck
	ClErrCheck
	RetryableErrCheck
	ConflictErrCheck
	PreconditionFailedErrCheck
}
