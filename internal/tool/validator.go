package tool

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var jobPathRe = regexp.MustCompile(`^[^/]`)

var jobPathRules = []validation.Rule{
	validation.Required,
	validation.Match(jobPathRe).Error("не должен начинаться с /"),
}

var projectNameRules = []validation.Rule{
	validation.Required,
}

func (i ListJobsInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.ProjectName, projectNameRules...),
	)
}

func (i GetJobInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, jobPathRules...),
	)
}

func (i ListBuildsInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, jobPathRules...),
	)
}

func (i GetLastBuildInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, jobPathRules...),
	)
}

func (i BuildRefInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, jobPathRules...),
		validation.Field(&i.BuildNumber, validation.Required, validation.Min(int64(1))),
	)
}

func (i GetBuildLogInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, jobPathRules...),
		validation.Field(&i.BuildNumber, validation.Required, validation.Min(int64(1))),
	)
}

func (i TriggerBuildInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, jobPathRules...),
	)
}

func (i SearchLogInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.JobPath, jobPathRules...),
		validation.Field(&i.BuildNumber, validation.Required, validation.Min(int64(1))),
		validation.Field(&i.Pattern, validation.Required),
	)
}

func (i ListNodesInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.ProjectName, projectNameRules...),
	)
}

func (i GetQueueInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.ProjectName, projectNameRules...),
	)
}
