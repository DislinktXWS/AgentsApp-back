package company_store

import "modules/dto"

func JobPositionMapper(jobPositionReq *dto.RequestJobPosition) JobPosition {
	jobPosition := JobPosition{
		Description: jobPositionReq.Description,
		Skills:      SkillsMapper(jobPositionReq),
		Name:        jobPositionReq.Name,
		Industry:    jobPositionReq.Industry,
		Position:    jobPositionReq.Position,
		CompanyID:   jobPositionReq.CompanyID,
	}
	return jobPosition
}

func SkillsMapper(jobPosition *dto.RequestJobPosition) []Skills {
	var skills []Skills
	for _, skill := range jobPosition.Skills {
		newSkill := Skills{
			Name:          skill.Name,
			Proficiency:   SkillProficiency(skill.Proficiency),
			JobPositionID: skill.JobPositionID,
		}
		skills = append(skills, newSkill)
	}
	return skills
}
