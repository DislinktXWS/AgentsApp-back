package company_store

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
	JobSalary           []JobSalary    `gorm:"foreignKey:CompanyID"`
	JobInterview        []JobInterview `gorm:"foreignKey:CompanyID"`
	JobPosition         []JobPosition  `gorm:"foreignKey:CompanyID"`
	OwnerID             string
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
