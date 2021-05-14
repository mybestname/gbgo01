package errors

import (
	"errors"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
const (
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

// Error is describes the cause of the error with structured details.
// For more details see https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto.
type Error struct {
	s *status.Status

	Domain   string            `json:"domain"`
	Reason   string            `json:"reason"`
	Metadata map[string]string `json:"metadata"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: domain = %s reason = %s metadata = %v", e.Domain, e.Reason, e.Metadata)
}

// New returns an error object for the code, message.
func New(code codes.Code, domain, reason, message string) *Error {
	return &Error{
		s:      status.New(code, message),
		Domain: domain,
		Reason: reason,
	}
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if se := FromError(err); err != nil {
		return se.Reason
	}
	return ""
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := *e
	err.Metadata = md
	return &err
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if target := new(Error); errors.As(err, &target) {
		return target
	}
	gs, ok := status.FromError(err)
	if ok {
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				return New(
					gs.Code(),
					d.Domain,
					d.Reason,
					gs.Message(),
				).WithMetadata(d.Metadata)
			}
		}
	}
	return New(gs.Code(), "", "", err.Error())
}