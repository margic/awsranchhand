package cmd

import (
	"errors"

	"github.com/abice/gencheck"
)

// Validate is an automatically generated validation method provided by
// gencheck.
// See https://github.com/abice/gencheck for more details.
func (s labelRequest) Validate() error {

	vErrors := make(gencheck.ValidationErrors, 0, 1)

	// BEGIN Host Validations
	// required
	if s.Host == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("labelRequest", "Host", "required", errors.New("is required")))
	}
	// END Host Validations

	if len(vErrors) > 0 {
		return vErrors
	}

	return nil
}
