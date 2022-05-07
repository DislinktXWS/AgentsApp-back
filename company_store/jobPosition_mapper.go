package company_store

import "modules/dto"

func JobPositionMapper(jobInterviewReq *dto.RequestJobPosition) JobPosition {
	jobPosition := JobPosition{
		Description:  jobInterviewReq.Description,
		Skills:       jobInterviewReq.Skills,
		WorkingHours: jobInterviewReq.WorkingHours,
		Position:     jobInterviewReq.Position,
		CompanyID:    jobInterviewReq.CompanyID,
	}
	return jobPosition
}
