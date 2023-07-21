package validation;

import(
	"github.com/pkg/errors"
	proto "github.com/KhaledHosseini/play-microservices/scheduler/scheduler-service/proto"
)

type CreateJobRequestValidator struct {*proto.CreateJobRequest}
func (v *CreateJobRequestValidator) Validate() error {
	if v.GetName() == "" || v.GetJobData() == "" {
		return errors.New("Empty fields are not allowed")
	}
	return nil
}

type GetJobRequestValidator struct {*proto.GetJobRequest}
func (v *GetJobRequestValidator) Validate() error {
	if v.GetId() == "" {
		return errors.New("Empty fields are not allowed")
	}
	return nil
}

type ListJobsRequestValidator struct {*proto.ListJobsRequest}
func (v *ListJobsRequestValidator) Validate() error {
	if v.GetSize() <= 0 || v.GetPage() <= 0 {
		return errors.New("Zero value fields are not allowed")
	}
	return nil
}
type UpdateJobRequestValidator struct {*proto.UpdateJobRequest}
func (v *UpdateJobRequestValidator) Validate() error {
	if v.GetName() == "" || v.GetJobData() == "" || v.GetId() == "" {
		return errors.New("Empty fields are not allowed")
	}
	return nil
}

type DeleteJobRequestValidator struct {*proto.DeleteJobRequest}
func (v *DeleteJobRequestValidator) Validate() error {
	if v.GetId() == "" {
		return errors.New("Empty fields are not allowed")
	}
	return nil
}