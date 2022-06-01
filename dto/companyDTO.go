package dto

import (
	"time"
)

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
	OwnerID             int    `json:"ownerID"`
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
}

type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseLogin struct {
	Token string `json:"token"`
}

type ResponseValidate struct {
	Username string `json:"username"`
	Role     int    `json:"role"`
}
