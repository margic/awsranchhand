package cmd

import (
	"errors"

	"github.com/abice/gencheck"
)

// Validate is an automatically generated validation method provided by
// gencheck.
// See https://github.com/abice/gencheck for more details.
func (s finishRequest) Validate() error {

	vErrors := make(gencheck.ValidationErrors, 0, 2)

	// BEGIN Stack Validations
	// required
	if s.Stack == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("finishRequest", "Stack", "required", errors.New("is required")))
	}
	// END Stack Validations

	// BEGIN Service Validations
	// required
	if s.Service == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("finishRequest", "Service", "required", errors.New("is required")))
	}
	// END Service Validations

	if len(vErrors) > 0 {
		return vErrors
	}

	return nil
}
