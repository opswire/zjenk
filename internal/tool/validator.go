package tool

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var jobPathRe = regexp.MustCompile(`^[^/]`)

func (i GetJobInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, validation.Required, validation.Match(jobPathRe).Error("must not start with /")),
	)
}

func (i ListBuildsInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, validation.Required, validation.Match(jobPathRe).Error("must not start with /")),
	)
}

func (i BuildRefInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, validation.Required, validation.Match(jobPathRe).Error("must not start with /")),
		validation.Field(&i.BuildNumber, validation.Required, validation.Min(int64(1))),
	)
}

func (i TriggerBuildInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, validation.Required, validation.Match(jobPathRe).Error("must not start with /")),
	)
}

func (i SearchLogInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, validation.Required, validation.Match(jobPathRe).Error("must not start with /")),
		validation.Field(&i.BuildNumber, validation.Required, validation.Min(int64(1))),
		validation.Field(&i.Pattern, validation.Required),
	)
}
