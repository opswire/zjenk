package tool

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (i GetJobInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobURL, validation.Required, is.URL),
	)
}

func (i ListBuildsInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobURL, validation.Required, is.URL),
	)
}

func (i BuildRefInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobURL, validation.Required, is.URL),
		validation.Field(&i.BuildNumber, validation.Required, validation.Min(int64(1))),
	)
}

func (i TriggerBuildInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobURL, validation.Required, is.URL),
	)
}

func (i SearchLogInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobURL, validation.Required, is.URL),
		validation.Field(&i.BuildNumber, validation.Required, validation.Min(int64(1))),
		validation.Field(&i.Pattern, validation.Required),
	)
}
