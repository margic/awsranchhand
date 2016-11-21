package cmd

import (
	"errors"

	"github.com/abice/gencheck"
)

// Validate is an automatically generated validation method provided by
// gencheck.
// See https://github.com/abice/gencheck for more details.
func (s rancherOpts) Validate() error {

	vErrors := make(gencheck.ValidationErrors, 0, 3)

	// BEGIN URL Validations
	// required
	if s.URL == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("rancherOpts", "URL", "required", errors.New("is required")))
	}
	// END URL Validations

	// BEGIN AccessKey Validations
	// required
	if s.AccessKey == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("rancherOpts", "AccessKey", "required", errors.New("is required")))
	}
	// END AccessKey Validations

	// BEGIN SecretKey Validations
	// required
	if s.SecretKey == "" {
		vErrors = append(vErrors, gencheck.NewFieldError("rancherOpts", "SecretKey", "required", errors.New("is required")))
	}
	// END SecretKey Validations

	if len(vErrors) > 0 {
		return vErrors
	}

	return nil
}
