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

type RequestJobSalary struct {
	Position  string `json:"position"`
	Salary    int    `json:"salary"`
	CompanyID int
}

type RequestJobInterview struct {
	Position   string `json:"position"`
	Impression string `json:"impression"`
	CompanyID  int
}

type RequestJobPosition struct {
	Position     string `json:"position"`
	WorkingHours string `json:"workingHours"`
	Description  string `json:"description"`
	Skills       string `json:"skills"`
	CompanyID    int
}

type RequestComment struct {
	Content   string `json:"content"`
	CompanyID int
}

type RequestAcceptCompany struct {
	ID     int  `json:"ID"`
	Accept bool `json:"accept"`
}
