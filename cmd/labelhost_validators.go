package cmd

import (
	"errors"

	"github.com/abice/gencheck"
)

// Validate is an automatically generated validation method provided by
// gencheck.
// See https://github.com/abice/gencheck for more details.
func (s labelRequest) Validate() error {

	vErrors := make(gencheck.ValidationErrors, 0, 3)

	// BEGIN Host Validations
	// required
	if s.Host == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("labelRequest", "Host", "required", errors.New("is required")))
	}
	// END Host Validations

	// BEGIN Key Validations
	// required
	if s.Key == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("labelRequest", "Key", "required", errors.New("is required")))
	}
	// END Key Validations

	// BEGIN Value Validations
	// required
	if s.Value == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("labelRequest", "Value", "required", errors.New("is required")))
	}
	// END Value Validations

	if len(vErrors) > 0 {
		return vErrors
	}

	return nil
}
