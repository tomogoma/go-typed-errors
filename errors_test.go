package typederrs_test

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

type errorChecker struct {
	typederrs.ClErrCheck
	typederrs.AuthErrCheck
	typederrs.NotFoundErrCheck
	typederrs.NotImplErrCheck
}

func Example() {

	// embed relevant 'Checkers' in struct
	ms := struct {
		typederrs.NotFoundErrCheck
		typederrs.NotImplErrCheck

		// do something returns an error which can be checked for type
		doSomething func() error
	}{
		doSomething: func() error {
			// return a typed error
			return typederrs.NewNotFoundf("something went wrong %s", "here")
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
			err = typederrs.New(tc.message)
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
			err = typederrs.Newf(tc.message, tc.messageParams...)
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
	checker := errorChecker{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewAuth(tc.message)
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
	checker := errorChecker{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewAuthf(tc.message, tc.messageParams...)
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
	checker := errorChecker{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewClient(tc.message)
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
	checker := errorChecker{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewClientf(tc.message, tc.messageParams...)
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
	checker := errorChecker{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewNotImplemented()
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
	checker := errorChecker{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewNotImplementedf(tc.message, tc.messageParams...)
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
	checker := errorChecker{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewForbidden(tc.message)
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
	checker := errorChecker{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewForbiddenf(tc.message, tc.messageParams...)
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
	checker := errorChecker{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewUnauthorized(tc.message)
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
	checker := errorChecker{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewUnauthorizedf(tc.message, tc.messageParams...)
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
	checker := errorChecker{}
	for _, tc := range messageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewNotFound(tc.message)
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
	checker := errorChecker{}
	for _, tc := range fmtdMessageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			err = typederrs.NewNotFoundf(tc.message, tc.messageParams...)
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
