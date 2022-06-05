package company_store

import (
	"fmt"
	"modules/dto"
	"modules/utils"
	"net/http"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CompanyStore struct {
	db *gorm.DB
}

func New() (*CompanyStore, error) {
	ts := &CompanyStore{}

	host := "localhost"
	user := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := "AgentDB"
	dbport := "5432"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, dbport)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	ts.db = db
	err = ts.db.AutoMigrate(&User{}, &Company{}, &JobSalary{}, &JobInterview{}, &JobPosition{}, &Comment{}, &Skills{})
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (ts *CompanyStore) CreateCompany(companyReq dto.RequestCompany) int {
	company := CompanyMapper(&companyReq)
	ts.db.Create(&company)
	return company.ID
}

func (ts *CompanyStore) UpdateCompany(companyReq dto.RequestCompany) int {
	//ts.db.Model(&Company{}).Updates(&company)
	company := &Company{}
	ts.db.First(&company, companyReq.ID)
	*company = CompanyMapper(&companyReq)
	ts.db.Save(&company)
	return company.ID
}

func (ts *CompanyStore) AcceptCompany(companyReq dto.RequestAcceptCompany) int {
	ts.db.Model(&Company{}).Where("id = ?", companyReq.ID).Update("accepted", companyReq.Accept)
	ts.db.Model(&Company{}).Where("id = ?", companyReq.ID).Update("checked", true)
	if companyReq.Accept {
		var company Company
		result := ts.db.Find(&company, Company{ID: companyReq.ID})
		if result.RowsAffected > 0 {
			ts.db.Model(&User{}).Where("id = ?", company.OwnerID).Update("role", 2)
		}
	}
	return companyReq.ID
}

func (ts *CompanyStore) GetAllCompanies() []Company {
	var companies []Company
	ts.db.Find(&companies, "accepted = true")
	return companies
}

func (ts *CompanyStore) GetCompanyRequests() []Company {
	var companies []Company
	ts.db.Find(&companies, "checked = false")
	return companies
}

func (ts *CompanyStore) GetAllJobPositions() []JobPosition {
	var jobPositions []JobPosition
	var jobPositionsWithSkills []JobPosition
	ts.db.Find(&jobPositions)
	for _, jobPosition := range jobPositions {
		for _, skill := range ts.GetSkillsByJobPosition(jobPosition.ID) {
			jobPosition.Skills = append(jobPosition.Skills, skill)
		}
		jobPositionsWithSkills = append(jobPositionsWithSkills, jobPosition)
	}
	return jobPositionsWithSkills
}

func (ts *CompanyStore) GetSkillsByJobPosition(id int) []Skills {
	var skills []Skills
	jobPositionID := strconv.Itoa(id)
	result := ts.db.Find(&skills, "job_position_id = "+jobPositionID)
	if result.RowsAffected > 0 {
		return skills
	}
	return nil
}

func (ts *CompanyStore) GetOwnersCompanies(id int) ([]Company, error) {
	var company []Company
	ownerID := strconv.Itoa(id)
	result := ts.db.Find(&company, "owner_id = "+ownerID)

	if result.RowsAffected > 0 {
		return company, nil
	}

	return company, fmt.Errorf("company with ownerId=%d not found", id)
}

func (ts *CompanyStore) GetCompanyById(id int) (Company, error) {
	var company Company
	companyID := strconv.Itoa(id)
	result := ts.db.Find(&company, "id = "+companyID)

	if result.RowsAffected > 0 {
		return company, nil
	}

	return company, fmt.Errorf("company with ownerId=%d not found", id)
}

func (ts *CompanyStore) GetJobSalary(id int) ([]JobSalary, error) {
	var jobSalaries []JobSalary
	companyID := strconv.Itoa(id)
	result := ts.db.Find(&jobSalaries, "company_id = "+companyID)

	if result.RowsAffected > 0 {
		return jobSalaries, nil
	}

	return jobSalaries, fmt.Errorf("company with id = %d does not have any public salaries", id)
}

func (ts *CompanyStore) DeleteJobSalary(id int) error {
	result := ts.db.Delete(&JobSalary{}, id)
	if result.RowsAffected > 0 {
		return nil
	}

	return fmt.Errorf("job position with id = %d not found", id)
}

func (ts *CompanyStore) CreateJobSalary(jobSalaryReq dto.RequestJobSalary) int {
	jobSalary := JobSalaryMapper(&jobSalaryReq)
	ts.db.Create(&jobSalary)
	return jobSalary.ID
}

func (ts *CompanyStore) CreateJobInterview(jobInterviewReq dto.RequestJobInterview) int {
	jobInterview := JobInterviewMapper(&jobInterviewReq)
	ts.db.Create(&jobInterview)
	return jobInterview.ID
}

func (ts *CompanyStore) GetJobInterview(id int) ([]JobInterview, error) {
	var interviews []JobInterview
	ownerID := strconv.Itoa(id)
	result := ts.db.Find(&interviews, "company_id = "+ownerID)

	if result.RowsAffected > 0 {
		return interviews, nil
	}

	return interviews, fmt.Errorf("company with id=%d does not have any interviews feed", id)
}

func (ts *CompanyStore) CreateJobPosition(jobPositionReq dto.RequestJobPosition) int {
	jobPosition := JobPositionMapper(&jobPositionReq)
	for _, skill := range jobPosition.Skills {
		skill.ID = jobPosition.ID
	}
	ts.db.Create(&jobPosition)
	return jobPosition.ID
}

func (ts *CompanyStore) GetJobPosition(id int) ([]JobPosition, error) {
	var positions []JobPosition
	ownerID := strconv.Itoa(id)
	result := ts.db.Find(&positions, "company_id = "+ownerID)

	if result.RowsAffected > 0 {
		return positions, nil
	}

	return positions, fmt.Errorf("company with id=%d does not have any new positions", id)
}

func (ts *CompanyStore) CreateComment(commentReq dto.RequestComment) int {
	jobInterview := CommentMapper(&commentReq)
	ts.db.Create(&jobInterview)
	return jobInterview.ID
}

func (ts *CompanyStore) RegisterUser(userReq dto.RequestUser) User {
	user := UserMapper(&userReq)
	ts.db.Create(&user)
	return user
}

func (ts *CompanyStore) LoginUser(loginReq dto.RequestLogin) (string, int) {
	var user User
	result := ts.db.Find(&user, User{Username: loginReq.Username})
	if result.RowsAffected == 0 {
		return "", http.StatusNotFound
	}
	match := utils.CheckPasswordHash(loginReq.Password, user.Password)
	if !match {
		return "", http.StatusNotFound
	}
	secretKey := os.Getenv("JWT_SECRET_KEY")
	wrapper := JwtWrapper{SecretKey: secretKey, ExpirationHours: 1}
	token, _ := wrapper.GenerateToken(&user)
	return token, http.StatusOK
}

func (ts *CompanyStore) GetComment(id int) ([]Comment, error) {
	var comments []Comment
	companyID := strconv.Itoa(id)
	result := ts.db.Find(&comments, "company_id = "+companyID)

	if result.RowsAffected > 0 {
		return comments, nil
	}

	return comments, fmt.Errorf("company with id=%d does not have any new comments", id)
}

func (ts *CompanyStore) Validate(token string) (int, int, string, int) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	wrapper := JwtWrapper{SecretKey: secretKey, ExpirationHours: 1}
	claims, err := wrapper.ValidateToken(token)
	if err != nil {
		return http.StatusBadRequest, -1, "", -1
	}
	var user User
	result := ts.db.Find(&user, User{Username: claims.Username})
	if result.RowsAffected == 0 {
		return http.StatusBadRequest, -1, "", -1
	}
	return http.StatusOK, claims.Id, claims.Username, claims.Role
}

func (ts *CompanyStore) Close() error {
	db, err := ts.db.DB()
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}
