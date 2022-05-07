package company_store

import "modules/dto"

func jobPositionMapper(jobPositionReq *dto.RequestJobPosition) JobPosition {
	jobPosition := JobPosition{
		Position:  jobPositionReq.Position,
		Salary:    jobPositionReq.Salary,
		CompanyID: jobPositionReq.CompanyID,
	}
	return jobPosition
}
