package httperr

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestRestError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    RestError
		want string
	}{
		{
			name: "Test Error 1",
			e:    RestError{ErrError: "Error 1"},
			want: "Error 1",
		},
		{
			name: "Test Error 2",
			e:    RestError{ErrError: "Error 2"},
			want: "Error 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("RestError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestError_Status(t *testing.T) {
	tests := []struct {
		name string
		e    RestError
		want int
	}{
		{
			name: "Test Status 400",
			e:    RestError{ErrStatus: 400},
			want: 400,
		},
		{
			name: "Test Status 404",
			e:    RestError{ErrStatus: 404},
			want: 404,
		},
		{
			name: "Test Status 500",
			e:    RestError{ErrStatus: 500},
			want: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Status(); got != tt.want {
				t.Errorf("RestError.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestError_Causes(t *testing.T) {
	tests := []struct {
		name string
		e    RestError
		want interface{}
	}{
		{
			name: "Test Cause 1",
			e:    RestError{ErrCauses: "Cause 1"},
			want: "Cause 1",
		},
		{
			name: "Test Cause 2",
			e:    RestError{ErrCauses: "Cause 2"},
			want: "Cause 2",
		},
		{
			name: "Test Cause 3",
			e:    RestError{ErrCauses: errors.New("Cause 3")},
			want: errors.New("Cause 3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Causes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RestError.Causes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRestError(t *testing.T) {
	tests := []struct {
		name   string
		status int
		err    string
		causes interface{}
		want   RestErr
	}{
		{
			name:   "Test NewRestError 1",
			status: 400,
			err:    "Bad Request",
			causes: "Cause 1",
			want:   RestError{ErrStatus: 400, ErrError: "Bad Request", ErrCauses: "Cause 1"},
		},
		{
			name:   "Test NewRestError 2",
			status: 404,
			err:    "Not Found",
			causes: "Cause 2",
			want:   RestError{ErrStatus: 404, ErrError: "Not Found", ErrCauses: "Cause 2"},
		},
		{
			name:   "Test NewRestError 3",
			status: 500,
			err:    "Internal Server Error",
			causes: "Cause 3",
			want:   RestError{ErrStatus: 500, ErrError: "Internal Server Error", ErrCauses: "Cause 3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRestError(tt.status, tt.err, tt.causes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRestError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRestErrorWithMessage(t *testing.T) {
	tests := []struct {
		name   string
		status int
		err    string
		causes interface{}
		want   RestErr
	}{
		{
			name:   "Test NewRestErrorWithMessage 1",
			status: 400,
			err:    "Bad Request",
			causes: "Cause 1",
			want:   RestError{ErrStatus: 400, ErrError: "Bad Request", ErrCauses: "Cause 1"},
		},
		{
			name:   "Test NewRestErrorWithMessage 2",
			status: 404,
			err:    "Not Found",
			causes: "Cause 2",
			want:   RestError{ErrStatus: 404, ErrError: "Not Found", ErrCauses: "Cause 2"},
		},
		{
			name:   "Test NewRestErrorWithMessage 3",
			status: 500,
			err:    "Internal Server Error",
			causes: "Cause 3",
			want:   RestError{ErrStatus: 500, ErrError: "Internal Server Error", ErrCauses: "Cause 3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRestErrorWithMessage(tt.status, tt.err, tt.causes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRestErrorWithMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want RestErr
	}{
		{
			name: "Test Deadline Exceeded Error",
			err:  context.DeadlineExceeded,
			want: NewRestError(http.StatusRequestTimeout, ErrRequestTimeoutError.Error(), context.DeadlineExceeded),
		},
		{
			name: "Test Unmarshal Error",
			err:  errors.New("Unmarshal error"),
			want: NewRestError(http.StatusBadRequest, ErrBadRequest.Error(), errors.New("Unmarshal error")),
		},
		{
			name: "Test RestError",
			err:  NewRestError(http.StatusForbidden, ErrForbidden.Error(), "Forbidden error"),
			want: NewRestError(http.StatusForbidden, ErrForbidden.Error(), "Forbidden error"),
		},
		{
			name: "Test Other Error",
			err:  errors.New("Other error"),
			want: NewRestError(http.StatusInternalServerError, ErrInternalServerError.Error(), errors.New("Other error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromError(tt.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRestErrorFromBytes(t *testing.T) {
	tests := []struct {
		name    string
		bytes   []byte
		wantErr bool
		want    RestErr
	}{
		{
			name:    "Valid JSON",
			bytes:   []byte(`{"status":400,"error":"Bad Request"}`),
			wantErr: false,
			want:    NewRestError(400, "Bad Request", nil),
		},
		{
			name:    "Invalid JSON",
			bytes:   []byte(`invalid`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRestErrorFromBytes(tt.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRestErrorFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRestErrorFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBadRequestError(t *testing.T) {
	tests := []struct {
		name   string
		causes interface{}
		want   RestErr
	}{
		{
			name:   "Test NewBadRequestError with string cause",
			causes: "string cause",
			want:   RestError{ErrStatus: http.StatusBadRequest, ErrError: ErrBadRequest.Error(), ErrCauses: "string cause"},
		},
		{
			name:   "Test NewBadRequestError with error cause",
			causes: errors.New("error cause"),
			want:   RestError{ErrStatus: http.StatusBadRequest, ErrError: ErrBadRequest.Error(), ErrCauses: errors.New("error cause")},
		},
		{
			name:   "Test NewBadRequestError with nil cause",
			causes: nil,
			want:   RestError{ErrStatus: http.StatusBadRequest, ErrError: ErrBadRequest.Error(), ErrCauses: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBadRequestError(tt.causes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBadRequestError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNotFoundError(t *testing.T) {
	tests := []struct {
		name   string
		causes interface{}
		want   RestErr
	}{
		{
			name:   "Test NewNotFoundError with string cause",
			causes: "string cause",
			want:   RestError{ErrStatus: http.StatusNotFound, ErrError: ErrNotFound.Error(), ErrCauses: "string cause"},
		},
		{
			name:   "Test NewNotFoundError with error cause",
			causes: errors.New("error cause"),
			want:   RestError{ErrStatus: http.StatusNotFound, ErrError: ErrNotFound.Error(), ErrCauses: errors.New("error cause")},
		},
		{
			name:   "Test NewNotFoundError with nil cause",
			causes: nil,
			want:   RestError{ErrStatus: http.StatusNotFound, ErrError: ErrNotFound.Error(), ErrCauses: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotFoundError(tt.causes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	tests := []struct {
		name   string
		causes interface{}
		want   RestErr
	}{
		{
			name:   "Test NewUnauthorizedError with string cause",
			causes: "string cause",
			want:   RestError{ErrStatus: http.StatusUnauthorized, ErrError: ErrUnauthorized.Error(), ErrCauses: "string cause"},
		},
		{
			name:   "Test NewUnauthorizedError with error cause",
			causes: errors.New("error cause"),
			want:   RestError{ErrStatus: http.StatusUnauthorized, ErrError: ErrUnauthorized.Error(), ErrCauses: errors.New("error cause")},
		},
		{
			name:   "Test NewUnauthorizedError with nil cause",
			causes: nil,
			want:   RestError{ErrStatus: http.StatusUnauthorized, ErrError: ErrUnauthorized.Error(), ErrCauses: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUnauthorizedError(tt.causes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnauthorizedError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewForbiddenError(t *testing.T) {
	tests := []struct {
		name   string
		causes interface{}
		want   RestErr
	}{
		{
			name:   "Test NewForbiddenError with string cause",
			causes: "string cause",
			want:   RestError{ErrStatus: http.StatusForbidden, ErrError: ErrForbidden.Error(), ErrCauses: "string cause"},
		},
		{
			name:   "Test NewForbiddenError with error cause",
			causes: errors.New("error cause"),
			want:   RestError{ErrStatus: http.StatusForbidden, ErrError: ErrForbidden.Error(), ErrCauses: errors.New("error cause")},
		},
		{
			name:   "Test NewForbiddenError with nil cause",
			causes: nil,
			want:   RestError{ErrStatus: http.StatusForbidden, ErrError: ErrForbidden.Error(), ErrCauses: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewForbiddenError(tt.causes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewForbiddenError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInternalServerError(t *testing.T) {
	tests := []struct {
		name   string
		causes interface{}
		want   RestErr
	}{
		{
			name:   "Test NewInternalServerError with string cause",
			causes: "string cause",
			want:   RestError{ErrStatus: http.StatusInternalServerError, ErrError: ErrInternalServerError.Error(), ErrCauses: "string cause"},
		},
		{
			name:   "Test NewInternalServerError with error cause",
			causes: errors.New("error cause"),
			want:   RestError{ErrStatus: http.StatusInternalServerError, ErrError: ErrInternalServerError.Error(), ErrCauses: errors.New("error cause")},
		},
		{
			name:   "Test NewInternalServerError with nil cause",
			causes: nil,
			want:   RestError{ErrStatus: http.StatusInternalServerError, ErrError: ErrInternalServerError.Error(), ErrCauses: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInternalServerError(tt.causes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInternalServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}
