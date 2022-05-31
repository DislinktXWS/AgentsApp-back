package company_store

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"modules/dto"
	"os"
	"strconv"
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
	err = ts.db.AutoMigrate(&Company{}, &JobSalary{}, &JobInterview{}, &JobPosition{}, &Comment{})
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
	return companyReq.ID
}

func (ts *CompanyStore) GetAllCompanies() []Company {
	var companies []Company
	ts.db.Find(&companies)
	return companies
}

func (ts *CompanyStore) GetCompany(id int) ([]Company, error) {
	var company []Company
	ownerID := strconv.Itoa(id)
	result := ts.db.Find(&company, "owner_id = "+ownerID)

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
	jobInterview := JobPositionMapper(&jobPositionReq)
	ts.db.Create(&jobInterview)
	return jobInterview.ID
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

func (ts *CompanyStore) GetComment(id int) ([]Comment, error) {
	var comments []Comment
	ownerID := strconv.Itoa(id)
	result := ts.db.Find(&comments, "company_id = "+ownerID)

	if result.RowsAffected > 0 {
		return comments, nil
	}

	return comments, fmt.Errorf("company with id=%d does not have any new comments", id)
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
