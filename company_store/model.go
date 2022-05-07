package company_store

type Company struct {
	ID                  int           `json:"id"`
	Name                string        `json:"name"`
	Phone               string        `json:"phone"`
	Email               string        `json:"email"`
	Address             string        `json:"address"`
	City                string        `json:"city"`
	Country             string        `json:"country"`
	YearOfEstablishment int           `json:"yearOfEstablishment"`
	Industry            string        `json:"industry"`
	CompanyCulture      string        `json:"companyCulture"`
	Description         string        `json:"description"`
	Website             string        `json:"website"`
	JobPosition         []JobPosition `gorm:"foreignKey:CompanyID"`
	OwnerID             string
}

type JobPosition struct {
	ID        int    `json:"id"`
	Position  string `json:"position"`
	Salary    int    `json:"salary"`
	CompanyID int    `json:"companyID"`
}
