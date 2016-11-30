package cmd

import (
	"errors"
	"time"

	"github.com/abice/gencheck"
)

// Validate is an automatically generated validation method provided by
// gencheck.
// See https://github.com/abice/gencheck for more details.
func (s waitrequest) Validate() error {

	vErrors := make(gencheck.ValidationErrors, 0, 4)

	// BEGIN Stack Validations
	// required
	if s.Stack == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("waitrequest", "Stack", "required", errors.New("is required")))
	}
	// END Stack Validations

	// BEGIN Service Validations
	// required
	if s.Service == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("waitrequest", "Service", "required", errors.New("is required")))
	}
	// END Service Validations

	// BEGIN State Validations
	// required
	if s.State == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("waitrequest", "State", "required", errors.New("is required")))
	}
	// END State Validations

	// BEGIN Timeout Validations
	// required
	var zeroTimeout time.Duration
	if s.Timeout == zeroTimeout {
		vErrors = append(vErrors, gencheck.NewFieldError("waitrequest", "Timeout", "required", errors.New("is required")))
	}
	// END Timeout Validations

	if len(vErrors) > 0 {
		return vErrors
	}

	return nil
}
