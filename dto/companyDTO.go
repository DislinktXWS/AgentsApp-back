package dto

type RequestCompany struct {
	ID                  int    `json:"ID"`
	Name                string `json:"name"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	Address             string `json:"address"`
	YearOfEstablishment string `json:"yearOfEstablishment"`
	Industry            string `json:"industry"`
	CompanyCulture      string `json:"companyCulture"`
	Description         string `json:"description"`
	Website             string `json:"website"`
	OwnerID             int    `json:"ownerID"`
}
type ResponseId struct {
	Id int `json:"id"`
}
type RequestJobSalary struct {
	Position  string `json:"position"`
	Salary    int    `json:"salary"`
	CompanyID int    `json:"companyID"`
}
type RequestJobInterview struct {
	Position   string `json:"position"`
	Impression string `json:"impression"`
	CompanyID  int
}
type RequestJobPosition struct {
	ID          int      `json:"ID"`
	Name        string   `json:"name"`
	Position    string   `json:"position"`
	Industry    string   `json:"industry"`
	Description string   `json:"description"`
	Skills      []Skills `json:"skills"`
	CompanyID   int      `json:"companyID"`
}
type SkillProficiency int64

const (
	Novice SkillProficiency = iota
	AdvancedBeginner
	Proficient
	Expert
	Master
)

type Skills struct {
	Name          string           `json:"name"`
	Proficiency   SkillProficiency `json:"proficiency"`
	JobPositionID int              `json:"jobPositionID"`
}
type RequestComment struct {
	Content   string `json:"content"`
	CompanyID int    `json:"companyID"`
}
type RequestAcceptCompany struct {
	ID     int  `json:"ID"`
	Accept bool `json:"accept"`
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

type RequestUser struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	DateOfBirth string `json:"dateOfBirth"`
	Gender      Gender `json:"gender"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	Role        Role   `json:"role"`
	IsConnected bool   `json:"isConnected"`
	ApiKey      string `json:"apiKey"`
}
type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type ResponseLogin struct {
	Token string `json:"token"`
}
type ResponseValidate struct {
	ID       int    `json:"ID"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}
type Connection struct {
	ID     int    `json:"ID"`
	ApiKey string `json:"apiKey"`
}

type Password struct {
	Password string `json:"password"`
}
