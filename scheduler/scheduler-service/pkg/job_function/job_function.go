package job_function

import (
	"context"
	"fmt"

	q "github.com/reugn/go-quartz/quartz"
)

// Function represents an argument-less function which returns a generic type R and a possible error.
type Function[R any] func(context.Context) (R, error)

// FunctionJob represents a Job that invokes the passed Function, implements the quartz.Job interface.
type FunctionJob[R any] struct {
	function  *Function[R]
	desc      string
	Result    *R
	Error     error
	JobStatus q.JobStatus
	key       int
}

// NewFunctionJob returns a new FunctionJob without an explicit description.
func NewFunctionJobWithKey[R any](key int, function Function[R]) *FunctionJob[R] {
	return &FunctionJob[R]{
		function:  &function,
		desc:      fmt.Sprintf("FunctionJob:%p", &function),
		Result:    nil,
		Error:     nil,
		JobStatus: q.NA,
		key:       key,
	}
}

// Description returns the description of the FunctionJob.
func (f *FunctionJob[R]) Description() string {
	return f.desc
}

// Key returns the unique FunctionJob key.
func (f *FunctionJob[R]) Key() int {
	return f.key
}

// Execute is called by a Scheduler when the Trigger associated with this job fires.
// It invokes the held function, setting the results in Result and Error members.
func (f *FunctionJob[R]) Execute(ctx context.Context) {
	result, err := (*f.function)(ctx)
	if err != nil {
		f.JobStatus = q.FAILURE
		f.Result = nil
		f.Error = err
	} else {
		f.JobStatus = q.OK
		f.Error = nil
		f.Result = &result
	}
}
