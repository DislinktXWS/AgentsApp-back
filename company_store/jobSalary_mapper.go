package company_store

import (
	"modules/dto"
)

func JobSalaryMapper(jobPositionReq *dto.RequestJobSalary) JobSalary {
	jobSalary := JobSalary{
		Position:  jobPositionReq.Position,
		Salary:    jobPositionReq.Salary,
		CompanyID: jobPositionReq.CompanyID,
	}
	return jobSalary
}
