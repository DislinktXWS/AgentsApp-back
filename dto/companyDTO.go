package dto

type RequestCompany struct {
	ID                  int    `json:"ID"`
	Name                string `json:"name"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	Address             string `json:"address"`
	City                string `json:"city"`
	Country             string `json:"country"`
	YearOfEstablishment int    `json:"yearOfEstablishment"`
	Industry            string `json:"industry"`
	CompanyCulture      string `json:"companyCulture"`
	Description         string `json:"description"`
	Website             string `json:"website"`
	OwnerID             string
}

type ResponseId struct {
	Id int `json:"id"`
}

type RequestJobPosition struct {
	Position  string `json:"position"`
	Salary    int    `json:"salary"`
	CompanyID int
}
