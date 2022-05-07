package company_store

type Company struct {
	ID                  int    `json:"id"`
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
