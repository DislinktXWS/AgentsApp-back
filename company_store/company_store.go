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
	err = ts.db.AutoMigrate(&Company{})
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (ts *CompanyStore) CreateCompany(companyReq dto.RequestCompany) int {
	company := companyMapper(&companyReq)
	ts.db.Create(&company)
	return company.ID
}

func (ts *CompanyStore) UpdateCompany(companyReq dto.RequestCompany) int {
	//ts.db.Model(&Company{}).Updates(&company)
	company := &Company{}
	ts.db.First(&company, companyReq.ID)
	*company = companyMapper(&companyReq)
	ts.db.Save(&company)
	return company.ID
}

func (ts *CompanyStore) GetAllCompanies() []Company {
	var companies []Company
	ts.db.Find(&companies)
	return companies
}

func (ts *CompanyStore) GetCompany(id int) (Company, error) {
	var company Company
	ownerID := strconv.Itoa(id)
	result := ts.db.Find(&company, "owner_id = "+ownerID)

	if result.RowsAffected > 0 {
		return company, nil
	}

	return Company{}, fmt.Errorf("company with ownerId=%d not found", id)
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
