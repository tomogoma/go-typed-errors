package typederrs

import "fmt"

type IsAuthErrChecker interface {
	IsAuthError(error)bool
	IsUnauthorizedError(error)bool
	IsForbiddenError(error)bool
}

type IsNotFoundErrChecker interface {
	IsNotFoundError(error)bool
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

type AllErrChecker interface {
	IsAuthErrChecker
	IsNotFoundErrChecker
	IsNotImplErrChecker
	IsClErrChecker
	IsRetryableErrChecker
}

// Error implements the Error interface and helps distinguish whether an error
// is a client error or an auth error.
type Error struct {
	IsAuthErr           bool
	IsUnauthorizedErr   bool
	IsForbiddenErr      bool
	IsClErr             bool
	IsNotFoundErr       bool
	IsNotImplementedErr bool
	IsRetryableErr bool
	Data                interface{}
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

// Retryable returns true if this is error is not permanent and should
// be retried
func (e Error) Retryable() bool {
	return e.IsRetryableErr
}

// New creates a new error.
func New(data interface{}) Error {
	return Error{Data: data, IsClErr: false}
}

// Newf creates a new error with fmt.Printf formatting.
func Newf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return New(data)
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

// NewAuth creates a new auth error. It is not specific to the type of auth
// error. Use NewForbidden(string) or NewUnauthorized(string) to establish
// a more specific Auth error.
func NewAuth(data interface{}) Error {
	return Error{Data: data, IsAuthErr: true}
}

// NewForbidden creates a new forbidden auth error a la 403 (http.StatusForbidden) error.
// This will also resolve as an Auth error.
func NewForbidden(data interface{}) Error {
	return Error{Data: data, IsAuthErr: true, IsForbiddenErr: true}
}

// NewUnauthorized creates a new unauthorized auth error a la 401 (http.StatusUnauthorized) error.
// This will also resolve as an Auth error.
func NewUnauthorized(data interface{}) Error {
	return Error{Data: data, IsAuthErr: true, IsUnauthorizedErr: true}
}

// NewAuthf creates a new auth error with fmt.Printf style formatting.
func NewAuthf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewAuth(data)
}

// NewForbiddenf creates a new forbidden auth error with fmt.Printf style formatting.
// This will also resolve as an Auth error.
func NewForbiddenf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewForbidden(data)
}

// NewUnauthorizedf creates a new unauthorized auth error with fmt.Printf style formatting.
// This will also resolve as an Auth error.
func NewUnauthorizedf(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewUnauthorized(data)
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

// NewRetryable creates a new retryable error.
func NewRetryable(data interface{}) Error {
	return Error{Data: data, IsRetryableErr: true}
}

// NewRetryablef creates a new retryable error with fmt.Printf style formatting.
func NewRetryablef(format string, a ...interface{}) Error {
	data := fmt.Sprintf(format, a...)
	return NewRetryable(data)
}

// ClErrCheck is a helper struct that can be embedded in a custom struct to
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

// NotImplErrCheck is a helper struct that can be embedded in a custom struct to
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

// AuthErrCheck is a helper struct that can be embedded in a custom struct to
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

// NotFoundErrCheck is a helper struct that can be embedded in a custom struct to
// give the custom struct the extra method IsNotFoundError(err error). e.g:
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

// RetryableErrCheck is a helper struct that can be embedded in a custom struct to
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

// AllErrCheck is a helper struct that can be embedded in a custom struct to
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
}
