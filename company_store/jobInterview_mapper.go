package company_store

import (
	"modules/dto"
)

func JobInterviewMapper(jobInterviewReq *dto.RequestJobInterview) JobInterview {
	jobInterview := JobInterview{
		Impression: jobInterviewReq.Impression,
		Position:   jobInterviewReq.Position,
		CompanyID:  jobInterviewReq.CompanyID,
	}
	return jobInterview
}
