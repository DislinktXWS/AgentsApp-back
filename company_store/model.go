package company_store

import "time"

type Company struct {
	ID                  int            `json:"id"`
	Name                string         `json:"name"`
	Phone               string         `json:"phone"`
	Email               string         `json:"email"`
	Address             string         `json:"address"`
	YearOfEstablishment string         `json:"yearOfEstablishment"`
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
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Position    string   `json:"position"`
	Industry    string   `json:"industry"`
	Description string   `json:"description"`
	Skills      []Skills `gorm:"foreignKey:JobPositionID"`
	CompanyID   int      `json:"companyID"`
	IsShared    bool     `json:"isShared"`
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
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	Proficiency   SkillProficiency `json:"proficiency"`
	JobPositionID int              `json:"jobPositionID"`
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
	ID                int       `json:"id"`
	Username          string    `json:"username"`
	Name              string    `json:"name"`
	Surname           string    `json:"surname"`
	DateOfBirth       string    `json:"dateOfBirth"`
	Gender            Gender    `json:"gender"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	Password          string    `json:"password"`
	Role              Role      `json:"role"`
	Company           []Company `gorm:"foreignKey:OwnerID"`
	IsConnected       bool      `json:"isConnected"`
	ApiKey            string    `json:"apiKey"`
	Token             string    `json:"token"`
	TokenCreationDate time.Time `json:"tokenCreationDate"`
	IsVerified        bool      `json:"isVerified"`
	TwoAuth           bool      `json:"twoAuth"`
}

type TwoFactorAuth struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Totp     []byte `bson:"totp"`
}
