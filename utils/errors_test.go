package utils

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "error with code and message",
			err: &Error{
				Code:    "TEST_ERROR",
				Message: "Test error message",
			},
			expected: "[TEST_ERROR] Test error message",
		},
		{
			name: "error with details",
			err: &Error{
				Code:    "TEST_ERROR",
				Message: "Test error message",
				Details: "Additional details",
			},
			expected: "[TEST_ERROR] Test error message: Additional details",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_WithDetails(t *testing.T) {
	err := NewError("TEST_ERROR", "Test message")
	err = err.WithDetails("Additional details")

	if err.Details != "Additional details" {
		t.Errorf("WithDetails() failed, got %v", err.Details)
	}
}

func TestError_WithLocation(t *testing.T) {
	err := NewError("TEST_ERROR", "Test message")
	err = err.WithLocation()

	if err.File == "" || err.Line == 0 {
		t.Errorf("WithLocation() failed, got file=%v, line=%v", err.File, err.Line)
	}
}

func TestNewError(t *testing.T) {
	err := NewError("TEST_ERROR", "Test message")

	if err.Code != "TEST_ERROR" {
		t.Errorf("NewError() code = %v, want TEST_ERROR", err.Code)
	}

	if err.Message != "Test message" {
		t.Errorf("NewError() message = %v, want Test message", err.Message)
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapError(originalErr, "WRAP_ERROR", "Wrapped message")

	if wrappedErr.Code != "WRAP_ERROR" {
		t.Errorf("WrapError() code = %v, want WRAP_ERROR", wrappedErr.Code)
	}

	if wrappedErr.Details != "original error" {
		t.Errorf("WrapError() details = %v, want original error", wrappedErr.Details)
	}
}

func TestWrapError_Nil(t *testing.T) {
	wrappedErr := WrapError(nil, "WRAP_ERROR", "Wrapped message")

	if wrappedErr != nil {
		t.Errorf("WrapError(nil) should return nil, got %v", wrappedErr)
	}
}

func TestIsError(t *testing.T) {
	err := NewError("TEST_ERROR", "Test message")

	if !IsError(err, "TEST_ERROR") {
		t.Errorf("IsError() should return true for matching code")
	}

	if IsError(err, "OTHER_ERROR") {
		t.Errorf("IsError() should return false for non-matching code")
	}

	regularErr := errors.New("regular error")
	if IsError(regularErr, "TEST_ERROR") {
		t.Errorf("IsError() should return false for regular error")
	}
}

func TestGetErrorCode(t *testing.T) {
	err := NewError("TEST_ERROR", "Test message")

	if GetErrorCode(err) != "TEST_ERROR" {
		t.Errorf("GetErrorCode() = %v, want TEST_ERROR", GetErrorCode(err))
	}

	regularErr := errors.New("regular error")
	if GetErrorCode(regularErr) != "UNKNOWN" {
		t.Errorf("GetErrorCode() for regular error = %v, want UNKNOWN", GetErrorCode(regularErr))
	}
}

func TestErrorCollector(t *testing.T) {
	ec := NewErrorCollector()

	if ec.HasErrors() {
		t.Errorf("New ErrorCollector should not have errors")
	}

	ec.Add(NewError("ERROR1", "First error"))
	ec.Add(NewError("ERROR2", "Second error"))

	if !ec.HasErrors() {
		t.Errorf("ErrorCollector should have errors after adding")
	}

	errors := ec.Errors()
	if len(errors) != 2 {
		t.Errorf("ErrorCollector should have 2 errors, got %d", len(errors))
	}

	errorMsg := ec.Error()
	if errorMsg == "" {
		t.Errorf("ErrorCollector.Error() should return error message")
	}

	combinedErr := ec.ToError()
	if combinedErr == nil {
		t.Errorf("ErrorCollector.ToError() should return error")
	}
}

func TestErrorCollector_AddNil(t *testing.T) {
	ec := NewErrorCollector()
	ec.Add(nil)

	if ec.HasErrors() {
		t.Errorf("ErrorCollector should not have errors after adding nil")
	}
}

func TestWrapError_WithCustomError(t *testing.T) {
	// 测试 WrapError 函数处理 *Error 类型错误的情况
	customErr := NewError("ORIGINAL_ERROR", "Original message").WithDetails("Original details")
	wrappedErr := WrapError(customErr, "WRAP_ERROR", "Wrapped message")

	if wrappedErr.Code != "WRAP_ERROR" {
		t.Errorf("WrapError() code = %v, want WRAP_ERROR", wrappedErr.Code)
	}

	if wrappedErr.Details != "Original details" {
		t.Errorf("WrapError() should use custom error's details, got %v", wrappedErr.Details)
	}

	// 测试没有 Details 的情况
	customErrWithoutDetails := NewError("ORIGINAL_ERROR", "Original message without details")
	wrappedErrWithoutDetails := WrapError(customErrWithoutDetails, "WRAP_ERROR", "Wrapped message")

	if wrappedErrWithoutDetails.Details != "Original message without details" {
		t.Errorf("WrapError() should use custom error's message when no details, got %v", wrappedErrWithoutDetails.Details)
	}
}

func TestErrorCollector_Error_Empty(t *testing.T) {
	// 测试空的 ErrorCollector.Error() 方法返回空字符串
	ec := NewErrorCollector()
	errorMsg := ec.Error()

	if errorMsg != "" {
		t.Errorf("Empty ErrorCollector.Error() should return empty string, got %v", errorMsg)
	}
}

func TestErrorCollector_ToError_Empty(t *testing.T) {
	// 测试空的 ErrorCollector.ToError() 方法返回 nil
	ec := NewErrorCollector()
	err := ec.ToError()

	if err != nil {
		t.Errorf("Empty ErrorCollector.ToError() should return nil, got %v", err)
	}
}
