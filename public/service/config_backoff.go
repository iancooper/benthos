// Copyright 2025 Redpanda Data, Inc.

package service

import (
	"github.com/cenkalti/backoff/v4"
)

// NewBackOffField defines a new object type config field that describes an
// exponential back off policy, often used for timing retry attempts. It is then
// possible to extract a *backoff.ExponentialBackOff from the resulting parsed
// config with the method FieldBackOff.
//
// It is possible to configure a back off policy that has no upper bound (no
// maximum elapsed time set). In cases where this would be problematic the field
// allowUnbounded should be set `false` in order to add linting rules that
// ensure an upper bound is set.
//
// The defaults struct is optional, and if provided will be used to establish
// default values for time interval fields. Otherwise the chosen defaults result
// in one minute of retry attempts, starting at 500ms intervals.
func NewBackOffField(name string, allowUnbounded bool, defaults *backoff.ExponentialBackOff) *ConfigField {
	var (
		initDefault       = "500ms"
		maxDefault        = "10s"
		maxElapsedDefault = "1m"
	)
	if defaults != nil {
		initDefault = defaults.InitialInterval.String()
		maxDefault = defaults.MaxInterval.String()
		maxElapsedDefault = defaults.MaxElapsedTime.String()
	}

	maxElapsedTime := NewDurationField("max_elapsed_time").
		Description("The maximum overall period of time to spend on retry attempts before the request is aborted.").
		Default(maxElapsedDefault).Example("1m").Example("1h")
	if allowUnbounded {
		maxElapsedTime.field.Description += " Setting this value to a zeroed duration (such as `0s`) will result in unbounded retries."
	}

	// TODO: Add linting rule to ensure we aren't unbounded if necessary.
	return NewObjectField(name,
		NewDurationField("initial_interval").
			Description("The initial period to wait between retry attempts.").
			Default(initDefault).Example("50ms").Example("1s"),
		NewDurationField("max_interval").
			Description("The maximum period to wait between retry attempts").
			Default(maxDefault).Example("5s").Example("1m"),
		maxElapsedTime,
	).Description("Determine time intervals and cut offs for retry attempts.")
}

// FieldBackOff accesses a field from a parsed config that was defined with
// NewBackoffField and returns a *backoff.ExponentialBackOff, or an error if the
// configuration was invalid.
func (p *ParsedConfig) FieldBackOff(path ...string) (*backoff.ExponentialBackOff, error) {
	b := backoff.NewExponentialBackOff()

	var err error
	if b.InitialInterval, err = p.FieldDuration(append(path, "initial_interval")...); err != nil {
		return nil, err
	}
	if b.MaxInterval, err = p.FieldDuration(append(path, "max_interval")...); err != nil {
		return nil, err
	}
	if b.MaxElapsedTime, err = p.FieldDuration(append(path, "max_elapsed_time")...); err != nil {
		return nil, err
	}

	return b, nil
}

// NewBackOffToggledField defines a new object type config field that describes
// an exponential back off policy, often used for timing retry attempts. It is
// then possible to extract a *backoff.ExponentialBackOff from the resulting
// parsed config with the method FieldBackOff. This Toggled variant includes a
// field `enabled` that is `false` by default.
//
// It is possible to configure a back off policy that has no upper bound (no
// maximum elapsed time set). In cases where this would be problematic the field
// allowUnbounded should be set `false` in order to add linting rules that
// ensure an upper bound is set.
//
// The defaults struct is optional, and if provided will be used to establish
// default values for time interval fields. Otherwise the chosen defaults result
// in one minute of retry attempts, starting at 500ms intervals.
func NewBackOffToggledField(name string, allowUnbounded bool, defaults *backoff.ExponentialBackOff) *ConfigField {
	var (
		initDefault       = "500ms"
		maxDefault        = "10s"
		maxElapsedDefault = "1m"
	)
	if defaults != nil {
		initDefault = defaults.InitialInterval.String()
		maxDefault = defaults.MaxInterval.String()
		maxElapsedDefault = defaults.MaxElapsedTime.String()
	}

	maxElapsedTime := NewDurationField("max_elapsed_time").
		Description("The maximum overall period of time to spend on retry attempts before the request is aborted.").
		Default(maxElapsedDefault).Example("1m").Example("1h")
	if allowUnbounded {
		maxElapsedTime.field.Description += " Setting this value to a zeroed duration (such as `0s`) will result in unbounded retries."
	}

	// TODO: Add linting rule to ensure we aren't unbounded if necessary.
	return NewObjectField(name,
		NewBoolField("enabled").
			Description("Whether retries should be enabled.").
			Default(false),
		NewDurationField("initial_interval").
			Description("The initial period to wait between retry attempts.").
			Default(initDefault).Example("50ms").Example("1s"),
		NewDurationField("max_interval").
			Description("The maximum period to wait between retry attempts").
			Default(maxDefault).Example("5s").Example("1m"),
		maxElapsedTime,
	).Description("Determine time intervals and cut offs for retry attempts.")
}

// FieldBackOffToggled accesses a field from a parsed config that was defined
// with NewBackOffField and returns a *backoff.ExponentialBackOff and a boolean
// flag indicating whether retries are explicitly enabled, or an error if the
// configuration was invalid.
func (p *ParsedConfig) FieldBackOffToggled(path ...string) (boff *backoff.ExponentialBackOff, enabled bool, err error) {
	boff = backoff.NewExponentialBackOff()

	if enabled, err = p.FieldBool(append(path, "enabled")...); err != nil {
		return
	}
	if boff.InitialInterval, err = p.FieldDuration(append(path, "initial_interval")...); err != nil {
		return
	}
	if boff.MaxInterval, err = p.FieldDuration(append(path, "max_interval")...); err != nil {
		return
	}
	if boff.MaxElapsedTime, err = p.FieldDuration(append(path, "max_elapsed_time")...); err != nil {
		return
	}

	return
}
