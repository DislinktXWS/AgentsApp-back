package company_store

import "time"

type Company struct {
	ID                  int            `json:"id"`
	Name                string         `json:"name"`
	Phone               string         `json:"phone"`
	Email               string         `json:"email"`
	Address             string         `json:"address"`
	City                string         `json:"city"`
	Country             string         `json:"country"`
	YearOfEstablishment int            `json:"yearOfEstablishment"`
	Industry            string         `json:"industry"`
	CompanyCulture      string         `json:"companyCulture"`
	Description         string         `json:"description"`
	Website             string         `json:"website"`
	Accepted            bool           `json:"accepted"`
	Checked             bool           `json:"checked"`
	JobSalary           []JobSalary    `gorm:"foreignKey:CompanyID"`
	JobInterview        []JobInterview `gorm:"foreignKey:CompanyID"`
	JobPosition         []JobPosition  `gorm:"foreignKey:CompanyID"`
	OwnerID             int            `json:"ownerID"`
}

type JobSalary struct {
	ID        int    `json:"id"`
	Position  string `json:"position"`
	Salary    int    `json:"salary"`
	CompanyID int    `json:"companyID"`
}

type JobInterview struct {
	ID         int    `json:"id"`
	Position   string `json:"position"`
	Impression string `json:"impression"`
	CompanyID  int    `json:"companyID"`
}

type JobPosition struct {
	ID           int    `json:"id"`
	Position     string `json:"position"`
	WorkingHours string `json:"workingHours"`
	Description  string `json:"description"`
	Skills       string `json:"skills"`
	CompanyID    int    `json:"companyID"`
}

type Comment struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	CompanyID int    `json:"companyID"`
}

type Gender int64

const (
	Male Gender = iota
	Female
)

type Role int64

const (
	Regular Role = iota
	Admin
	Owner
)

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Birthday    time.Time `json:"birthday"`
	Gender      Gender    `json:"gender"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Password    string    `json:"password"`
	Role        Role      `json:"role"`
	Company     []Company `gorm:"foreignKey:OwnerID"`
}
