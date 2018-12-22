package errors_test

import (
	"fmt"
	"testing"

	"github.com/tomogoma/go-typed-errors"
)

type testCase struct {
	name          string
	message       string
	messageParams []interface{}
}

var errWithAllFlagsTrue = errors.Error{
	IsAuthErr: true, IsUnauthorizedErr: true, IsForbiddenErr: true, IsClErr: true,
	IsNotFoundErr: true, IsNotImplementedErr: true, IsRetryableErr: true,
	IsConflictErr: true, IsPreconditionFailedErr: true, Data: "",
}

func Example() {

	// embed relevant 'Checkers' in struct
	ms := struct {
		errors.AllErrCheck

		// do something returns an error which can be checked for type
		doSomething func() error
	}{
		doSomething: func() error {
			// return a typed error
			return errors.NewNotFoundf("something went wrong %s", "here")
		},
	}

	err := ms.doSomething()
	if err != nil {

		if ms.IsNotFoundError(err) {
			fmt.Println("resource not found")
			return
		}

		if ms.IsNotImplementedError(err) {
			fmt.Println("logic not implemented")
			return
		}

		// Act on generic error
		fmt.Printf("got a generic error: %v\n", err)
	}

	// Output: resource not found
}

func TestNew(t *testing.T) {
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.New(tc.message)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != tc.message {
				t.Fatalf("expected error message '%s', got '%s'",
					tc.message, err.Error())
			}
		})
	}
}

func TestNewf(t *testing.T) {
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.Newf(tc.message, tc.messageParams...)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != fmt.Sprintf(tc.message, tc.messageParams...) {
				t.Fatalf("expected error message '%s', got '%s'",
					fmt.Sprintf(tc.message, tc.messageParams...), err.Error())
			}
		})
	}
}

func TestNewAuth(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewAuth(tc.message)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != tc.message {
				t.Errorf("expected error message '%s', got '%s'",
					tc.message, err.Error())
			}
			if !checker.IsAuthError(err) {
				t.Errorf("Expected IsAuthError but got %v", err)
			}
		})
	}
}

func TestNewAuthf(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewAuthf(tc.message, tc.messageParams...)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != fmt.Sprintf(tc.message, tc.messageParams...) {
				t.Fatalf("expected error message '%s', got '%s'",
					fmt.Sprintf(tc.message, tc.messageParams...), err.Error())
			}
			if !checker.IsAuthError(err) {
				t.Errorf("Expected IsAuthError but got %v", err)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewClient(tc.message)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != tc.message {
				t.Errorf("expected error message '%s', got '%s'",
					tc.message, err.Error())
			}
			if !checker.IsClientError(err) {
				t.Errorf("Expected IsClientError but got %v", err)
			}
		})
	}
}

func TestNewClientf(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewClientf(tc.message, tc.messageParams...)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != fmt.Sprintf(tc.message, tc.messageParams...) {
				t.Fatalf("expected error message '%s', got '%s'",
					fmt.Sprintf(tc.message, tc.messageParams...), err.Error())
			}
			if !checker.IsClientError(err) {
				t.Errorf("Expected IsClientError but got %v", err)
			}
		})
	}
}

func TestNewNotImplemented(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewNotImplemented()
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != "not implemented" {
				t.Errorf("expected error message '%s', got '%s'",
					tc.message, err.Error())
			}
			if !checker.IsNotImplementedError(err) {
				t.Errorf("Expected IsNotImplementedError() true but got %t",
					checker.IsNotImplementedError(err))
			}
		})
	}
}

func TestNewNotImplementedf(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewNotImplementedf(tc.message, tc.messageParams...)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != fmt.Sprintf(tc.message, tc.messageParams...) {
				t.Fatalf("expected error message '%s', got '%s'",
					fmt.Sprintf(tc.message, tc.messageParams...), err.Error())
			}
			if !checker.IsNotImplementedError(err) {
				t.Errorf("Expected IsNotImplementedError() true but got %t",
					checker.IsNotImplementedError(err))
			}
		})
	}
}

func TestNewForbidden(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewForbidden(tc.message)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != tc.message {
				t.Errorf("expected error message '%s', got '%s'",
					tc.message, err.Error())
			}
			if !checker.IsAuthError(err) {
				t.Errorf("Expected IsAuthError() true but got %t",
					checker.IsAuthError(err))
			}
			if !checker.IsForbiddenError(err) {
				t.Errorf("Expected IsForbiddenError() true but got %t",
					checker.IsForbiddenError(err))
			}
		})
	}
}

func TestNewForbiddenf(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewForbiddenf(tc.message, tc.messageParams...)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != fmt.Sprintf(tc.message, tc.messageParams...) {
				t.Fatalf("expected error message '%s', got '%s'",
					fmt.Sprintf(tc.message, tc.messageParams...), err.Error())
			}
			if !checker.IsAuthError(err) {
				t.Errorf("Expected IsAuthError() true but got %t",
					checker.IsAuthError(err))
			}
			if !checker.IsForbiddenError(err) {
				t.Errorf("Expected IsForbiddenError() true but got %t",
					checker.IsForbiddenError(err))
			}
		})
	}
}

func TestNewUnauthorized(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewUnauthorized(tc.message)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != tc.message {
				t.Errorf("expected error message '%s', got '%s'",
					tc.message, err.Error())
			}
			if !checker.IsAuthError(err) {
				t.Errorf("Expected IsAuthError() true but got %t",
					checker.IsAuthError(err))
			}
			if !checker.IsUnauthorizedError(err) {
				t.Errorf("Expected IsUnauthorizedError() true but got %t",
					checker.IsUnauthorizedError(err))
			}
		})
	}
}

func TestNewUnauthorizedf(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewUnauthorizedf(tc.message, tc.messageParams...)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != fmt.Sprintf(tc.message, tc.messageParams...) {
				t.Fatalf("expected error message '%s', got '%s'",
					fmt.Sprintf(tc.message, tc.messageParams...), err.Error())
			}
			if !checker.IsAuthError(err) {
				t.Errorf("Expected IsAuthError() true but got %t",
					checker.IsAuthError(err))
			}
			if !checker.IsUnauthorizedError(err) {
				t.Errorf("Expected IsUnauthorizedError() true but got %t",
					checker.IsUnauthorizedError(err))
			}
		})
	}
}

func TestNewNotFound(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewNotFound(tc.message)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != tc.message {
				t.Errorf("expected error message '%s', got '%s'",
					tc.message, err.Error())
			}
			if !checker.IsNotFoundError(err) {
				t.Errorf("Expected IsNotFoundError() true but got %t",
					checker.IsNotFoundError(err))
			}
		})
	}
}

func TestNewNotFoundf(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewNotFoundf(tc.message, tc.messageParams...)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != fmt.Sprintf(tc.message, tc.messageParams...) {
				t.Fatalf("expected error message '%s', got '%s'",
					fmt.Sprintf(tc.message, tc.messageParams...), err.Error())
			}
			if !checker.IsNotFoundError(err) {
				t.Errorf("Expected IsNotFoundError() true but got %t",
					checker.IsNotFoundError(err))
			}
		})
	}
}

func TestNewRetryable(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewRetryable(tc.message)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != tc.message {
				t.Errorf("expected error message '%s', got '%s'",
					tc.message, err.Error())
			}
			if !checker.IsRetryableError(err) {
				t.Errorf("Expected IsRetryableError() true but got %t",
					checker.IsRetryableError(err))
			}
		})
	}
}

func TestNewRetryablef(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewRetryablef(tc.message, tc.messageParams...)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != fmt.Sprintf(tc.message, tc.messageParams...) {
				t.Fatalf("expected error message '%s', got '%s'",
					fmt.Sprintf(tc.message, tc.messageParams...), err.Error())
			}
			if !checker.IsRetryableError(err) {
				t.Errorf("Expected IsRetryableError() true but got %t",
					checker.IsRetryableError(err))
			}
		})
	}
}

func TestNewConflict(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewConflict(tc.message)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != tc.message {
				t.Errorf("expected error message '%s', got '%s'",
					tc.message, err.Error())
			}
			if !checker.IsConflictError(err) {
				t.Errorf("Expected IsConflictError() true but got %t",
					checker.IsRetryableError(err))
			}
		})
	}
}

func TestNewConflictf(t *testing.T) {
	var checker errors.AllErrChecker
	checker = &errors.AllErrCheck{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = errors.NewConflictf(tc.message, tc.messageParams...)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if err.Error() != fmt.Sprintf(tc.message, tc.messageParams...) {
				t.Fatalf("expected error message '%s', got '%s'",
					fmt.Sprintf(tc.message, tc.messageParams...), err.Error())
			}
			if !checker.IsConflictError(err) {
				t.Errorf("Expected IsConflictError() true but got %t",
					checker.IsConflictError(err))
			}
		})
	}
}

func TestNewPreconditionFailed(t *testing.T) {
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			err := errors.NewPreconditionFailed(tc.message)
			if err.Error() != tc.message {
				t.Errorf("expected error message '%s', got '%s'",
					tc.message, err.Error())
			}
			if !err.IsPreconditionFailedErr {
				t.Errorf("Expected IsPreconditionFailedErr to be true but got %t", err.IsPreconditionFailedErr)
			}
		})
	}
}

func TestNewPreconditionFailedf(t *testing.T) {
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			err := errors.NewPreconditionFailedf(tc.message, tc.messageParams...)

			if err.Error() != fmt.Sprintf(tc.message, tc.messageParams...) {
				t.Fatalf("expected error message '%s', got '%s'",
					fmt.Sprintf(tc.message, tc.messageParams...), err.Error())
			}
			if !err.IsPreconditionFailedErr {
				t.Errorf("Expected IsPreconditionFailedErr to be true but got %t", err.IsPreconditionFailedErr)
			}
		})
	}
}

func TestPreconditionFailedErrCheck_IsPreconditionFailedError(t *testing.T) {
	preconditionErr := errWithAllFlagsTrue
	nonPreconditionErr := errWithAllFlagsTrue
	nonPreconditionErr.IsPreconditionFailedErr = false
	checker := errors.AllErrCheck{}
	if ok := checker.IsPreconditionFailedError(preconditionErr); !ok {
		t.Fatalf("expected IsPreconditionFailedError true but got false on %+v", preconditionErr)
	}
	if ok := checker.IsPreconditionFailedError(nonPreconditionErr); ok {
		t.Fatalf("expected IsPreconditionFailedError false but got true on %+v", nonPreconditionErr)
	}
}

func messageTestCases() []testCase {
	return []testCase{
		{name: "has-message", message: "this error message"},
		{name: "has-message", message: ""},
	}
}

func fmtdMessageTestCases() []testCase {
	return []testCase{
		{name: "has-message", message: "this %s error message", messageParams: []interface{}{"parametized"}},
		{name: "has-message", message: ""},
	}
}
