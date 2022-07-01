package company_store

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"modules/dto"
	"modules/utils"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sec51/twofactor"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

type CompanyStore struct {
	db *gorm.DB
}

type ApiKey struct {
	apiKey string `json:"apiKey"`
}

func New() (*CompanyStore, error) {
	ts := &CompanyStore{}

	//ZA LOCAL
	//connectionStr := "host=localhost user=postgres password=bbogi1219 dbname=AgentDB port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	//ZA DOCKER
	connectionStr := "host=database user=root password=root dbname=AgentsDB port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	ts.db = db
	err = ts.db.AutoMigrate(&User{}, &Company{}, &JobSalary{}, &JobInterview{}, &JobPosition{}, &Comment{}, &Skills{}, &TwoFactorAuth{})
	if err != nil {
		return nil, err
	}

	infoFile, err := os.OpenFile("logs/info.log", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(infoFile, "INFO: ", log.LstdFlags|log.Lshortfile)

	errFile, err1 := os.OpenFile("logs/error.log", os.O_APPEND|os.O_WRONLY, 0666)
	if err1 != nil {
		log.Fatal(err1)
	}
	ErrorLogger = log.New(errFile, "ERROR: ", log.LstdFlags|log.Lshortfile)

	return ts, nil
}

func (ts *CompanyStore) VerifyTwoFactorAuthToken(username string, twoFactorAuthToken string) (status int64, error string, JWTtoken string) {
	var two_factor_auth TwoFactorAuth
	ts.db.Find(&two_factor_auth, "username = ?", username)
	var user User
	ts.db.Find(&user, "username = ?", username)
	otp, _ := twofactor.TOTPFromBytes(two_factor_auth.Totp, username)
	err := otp.Validate(twoFactorAuthToken)
	if err != nil {
		return http.StatusNotFound, err.Error(), ""
	}
	secretKey := os.Getenv("JWT_SECRET_KEY")
	wrapper := JwtWrapper{SecretKey: secretKey, ExpirationHours: 1}
	token, _ := wrapper.GenerateToken(&user)
	return http.StatusOK, "", token
}

func (ts *CompanyStore) GetTwoFactorAuth(username string) bool {
	var user User
	ts.db.Find(&user, "username = ?", username)
	return user.TwoAuth
}

func (ts *CompanyStore) ChangeTwoFactorAuth(username string) (qrCode string, error string) {
	var user User
	ts.db.Find(&user, "username = ?", username)
	if user.TwoAuth == true {
		ts.db.Model(&User{}).Where("id = ?", user.ID).Update("two_auth", false)
		var two_factor TwoFactorAuth
		ts.db.Find(&two_factor, "username = ?", username)
		ts.db.Delete(&TwoFactorAuth{}, two_factor.Id)
	} else {
		ts.db.Model(&User{}).Where("id = ?", user.ID).Update("two_auth", true)
		otp, _ := twofactor.NewTOTP(user.Email, user.Username, crypto.SHA1, 6)
		qrBytes, _ := otp.QR()
		base64QR := base64.StdEncoding.EncodeToString(qrBytes)
		twoFact := TwoAuthMapper(username, otp)
		ts.db.Create(&twoFact)
		return base64QR, ""
	}
	return "", ""
}

func (ts *CompanyStore) CreateCompany(companyReq dto.RequestCompany) int {
	company := CompanyMapper(&companyReq)
	ts.db.Create(&company)
	InfoLogger.Println("Action: 1, Message: Company " + company.Name + " created")
	return company.ID
}

func (ts *CompanyStore) UpdateCompany(companyReq dto.RequestCompany) int {
	//ts.db.Model(&Company{}).Updates(&company)
	company := &Company{}
	ts.db.First(&company, companyReq.ID)
	company.Address = companyReq.Address
	company.CompanyCulture = companyReq.CompanyCulture
	company.Description = companyReq.Description
	company.Email = companyReq.Email
	company.Industry = companyReq.Industry
	company.Name = companyReq.Name
	company.Phone = companyReq.Phone
	company.Website = companyReq.Website
	company.YearOfEstablishment = companyReq.YearOfEstablishment
	ts.db.Save(&company)
	InfoLogger.Println("Action: 2, Message: Company " + company.Name + " updated")
	return company.ID
}

func (ts *CompanyStore) AcceptCompany(companyReq dto.RequestAcceptCompany) int {
	ts.db.Model(&Company{}).Where("id = ?", companyReq.ID).Update("accepted", companyReq.Accept)
	ts.db.Model(&Company{}).Where("id = ?", companyReq.ID).Update("checked", true)
	if companyReq.Accept {
		var company Company
		result := ts.db.Find(&company, Company{ID: companyReq.ID})
		if result.RowsAffected > 0 {
			ts.db.Model(&User{}).Where("id = ?", company.OwnerID).Update("role", 2)
			InfoLogger.Println("Action: 3, Message: Company " + company.Name + " accepted")
		}
	}
	if !companyReq.Accept {
		ErrorLogger.Println("Action: 3, Message: Company rejected")
	}
	return companyReq.ID
}

func (ts *CompanyStore) GetAllCompanies() []Company {
	var companies []Company
	ts.db.Find(&companies, "accepted = true")
	return companies
}

func (ts *CompanyStore) GetCompanyRequests() []Company {
	var companies []Company
	ts.db.Find(&companies, "checked = false")
	return companies
}

func (ts *CompanyStore) GetAllJobPositions() []JobPosition {
	var jobPositions []JobPosition
	var jobPositionsWithSkills []JobPosition
	ts.db.Find(&jobPositions)
	for _, jobPosition := range jobPositions {
		for _, skill := range ts.GetSkillsByJobPosition(jobPosition.ID) {
			jobPosition.Skills = append(jobPosition.Skills, skill)
		}
		jobPositionsWithSkills = append(jobPositionsWithSkills, jobPosition)
	}
	return jobPositionsWithSkills
}

func (ts *CompanyStore) GetSkillsByJobPosition(id int) []Skills {
	var skills []Skills
	jobPositionID := strconv.Itoa(id)
	result := ts.db.Find(&skills, "job_position_id = "+jobPositionID)
	if result.RowsAffected > 0 {
		return skills
	}
	return nil
}

func (ts *CompanyStore) GetOwnersCompanies(id int) ([]Company, error) {
	var company []Company
	ownerID := strconv.Itoa(id)
	result := ts.db.Find(&company, "owner_id = "+ownerID)

	if result.RowsAffected > 0 {
		return company, nil
	}
	ErrorLogger.Println("Action: 4, Message: Owner " + ownerID + " does not have any company.")

	return company, fmt.Errorf("company with ownerId=%d not found", id)
}

func (ts *CompanyStore) GetCompanyById(id int) (Company, error) {
	var company Company
	companyID := strconv.Itoa(id)
	result := ts.db.Find(&company, "id = "+companyID)

	if result.RowsAffected > 0 {
		return company, nil
	}

	ErrorLogger.Println("Action: 4, Message: Company by the id " + companyID + " does not exist.")

	return company, fmt.Errorf("company with ownerId=%d not found", id)
}

func (ts *CompanyStore) GetJobSalary(id int) ([]JobSalary, error) {
	var jobSalaries []JobSalary
	companyID := strconv.Itoa(id)
	result := ts.db.Find(&jobSalaries, "company_id = "+companyID)

	if result.RowsAffected > 0 {
		return jobSalaries, nil
	}
	ErrorLogger.Println("Action: 4, Message: Company by the id " + companyID + " does not have any salaries.")
	return jobSalaries, fmt.Errorf("company with id = %d does not have any public salaries", id)
}

func (ts *CompanyStore) IsConnected(id int) (bool, error) {
	userId := strconv.Itoa(id)
	var user User
	result := ts.db.Find(&user, "id = "+userId)
	if result.RowsAffected > 0 {
		return user.IsConnected, nil
	}
	ErrorLogger.Println("Action: 4, Message: User by the id " + userId + " does not exist.")

	return false, fmt.Errorf("user not found")
}

func (ts *CompanyStore) IsJobPositionShared(id int) (bool, error) {
	jobPositionId := strconv.Itoa(id)
	var jobPosition JobPosition
	result := ts.db.Find(&jobPosition, "id = "+jobPositionId)
	if result.RowsAffected > 0 {
		return jobPosition.IsShared, nil
	}
	return false, fmt.Errorf("user not found")
}

func (ts *CompanyStore) GetUserApiKey(id int) string {
	userId := strconv.Itoa(id)
	var user User
	result := ts.db.Find(&user, "id = "+userId)
	if result.RowsAffected > 0 {
		return user.ApiKey
	}

	return ""
}

func (ts *CompanyStore) DeleteJobSalary(id int) error {
	result := ts.db.Delete(&JobSalary{}, id)
	if result.RowsAffected > 0 {
		return nil
	}
	ErrorLogger.Println("Action: 4, Message: Job salary does not exist.")
	return fmt.Errorf("job position with id = %d not found", id)
}

func (ts *CompanyStore) CreateJobSalary(jobSalaryReq dto.RequestJobSalary) int {
	jobSalary := JobSalaryMapper(&jobSalaryReq)
	ts.db.Create(&jobSalary)
	InfoLogger.Println("Action: 5, Message: New job salary.")
	return jobSalary.ID
}

func (ts *CompanyStore) CreateJobInterview(jobInterviewReq dto.RequestJobInterview) int {
	jobInterview := JobInterviewMapper(&jobInterviewReq)
	ts.db.Create(&jobInterview)
	InfoLogger.Println("Action: 6, Message: New job interview.")
	return jobInterview.ID
}

func (ts *CompanyStore) GetJobInterview(id int) ([]JobInterview, error) {
	var interviews []JobInterview
	ownerID := strconv.Itoa(id)
	result := ts.db.Find(&interviews, "company_id = "+ownerID)

	if result.RowsAffected > 0 {
		return interviews, nil
	}
	ErrorLogger.Println("Action: 4, Message: Job interview for the company does not exist.")
	return interviews, fmt.Errorf("company with id=%d does not have any interviews feed", id)
}

func (ts *CompanyStore) CreateJobPosition(jobPositionReq dto.RequestJobPosition) int {
	jobPosition := JobPositionMapper(&jobPositionReq)
	for _, skill := range jobPosition.Skills {
		skill.ID = jobPosition.ID
	}
	ts.db.Create(&jobPosition)
	InfoLogger.Println("Action: 7, Message: New job position.")

	return jobPosition.ID
}

func (ts *CompanyStore) SetJobPositionIsShared(jobPositionId int) {
	var jobPosition JobPosition
	result := ts.db.Find(&jobPosition, JobPosition{ID: jobPositionId})
	if result.RowsAffected > 0 {
		jobPosition.IsShared = true
	}
	ts.db.Save(&jobPosition)
}

func (ts *CompanyStore) ShareJobPosition(jobPositionReq dto.RequestJobPosition, apiKey string) (bool, int) {
	jobPosition := JobPositionMapper(&jobPositionReq)
	client := &http.Client{}

	jsonVal, err := json.Marshal(jobPosition)
	if err != nil {
		ErrorLogger.Println("Action: 9, Message: Invalid json value.")
	}

	//http request
	req, err := http.NewRequest(http.MethodPost, "https://localhost:8000/shareBusinessOffer", bytes.NewBuffer(jsonVal))
	if err != nil {
		ErrorLogger.Println("Action: 9, Message: Dislinkt app is not available.")
	}

	req.Header.Set("Content-Type", "application/json; charset-utf-8")
	req.Header.Set("ApiKey", apiKey)
	resp, _ := client.Do(req)
	if err != nil {
		ErrorLogger.Println("Action: 9, Message: Response is not valid.")
	}

	if resp.StatusCode == 500 {
		ErrorLogger.Println("Action: 9, Message: Response is not valid.")
		return false, -1
	}
	resp.Body.Close()
	InfoLogger.Println("Action: 8, Message: Job position shared.")
	return true, jobPosition.ID
}

func (ts *CompanyStore) GetJobPosition(id int) ([]JobPosition, error) {
	var jobPositions []JobPosition
	ownerID := strconv.Itoa(id)
	result := ts.db.Find(&jobPositions, "company_id = "+ownerID)

	if result.RowsAffected == 0 {
		ErrorLogger.Println("Action: 4, Message: Company does not have any new positions.")
		return jobPositions, fmt.Errorf("company with id=%d does not have any new positions", id)
	}

	var jobPositionsWithSkills []JobPosition
	for _, jobPosition := range jobPositions {
		for _, skill := range ts.GetSkillsByJobPosition(jobPosition.ID) {
			jobPosition.Skills = append(jobPosition.Skills, skill)
		}
		jobPositionsWithSkills = append(jobPositionsWithSkills, jobPosition)
	}
	return jobPositionsWithSkills, nil
}

func (ts *CompanyStore) CreateComment(commentReq dto.RequestComment) int {
	jobInterview := CommentMapper(&commentReq)
	ts.db.Create(&jobInterview)
	InfoLogger.Println("Action: 10, Message: New comment created.")
	return jobInterview.ID
}

func (ts *CompanyStore) ConnectWithDislinkt(username string, id int) string {
	apiKey := ""
	if checkIfUserExists(username) {
		apiKey = changeApiKey(username)
		//sendEmail(apiKey)
		var user User
		result := ts.db.Find(&user, User{ID: id})
		if result.RowsAffected > 0 {
			InfoLogger.Println("Action: 11, Message: Connected with dislinkt.")
			user.IsConnected = true
			user.ApiKey = apiKey
		}
		ts.db.Save(&user)
		fmt.Println(apiKey)
	} else {
		ErrorLogger.Println("Action: 12, Message: cannot connect with dislinkt.")
		fmt.Println("Username does not exist")
	}

	return apiKey
}

func checkIfUserExists(username string) bool {

	resp, err := http.Get("https://localhost:8000/users/userByUsername/" + username)
	if err != nil {
		ErrorLogger.Println("Action: 13, Message: User not found.")
		return false
	}

	//We Read the response body on the line below
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorLogger.Println("Action: 13, Message: Invalid response from dislinkt.")
		return false
	}

	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 500 {
		return false
	}

	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
	InfoLogger.Println("Action: 14, Message: User found.")
	return true
}

func changeApiKey(username string) string {
	client := &http.Client{}

	jsonVal, err := json.Marshal(username)
	if err != nil {
		ErrorLogger.Println("Action: 15, Message: Invalid json value.")
	}

	//http request
	req, err := http.NewRequest(http.MethodPut, "https://localhost:8000/users/user/apiKey/"+username, bytes.NewBuffer(jsonVal))
	if err != nil {
		ErrorLogger.Println("Action: 12, Message: Invalid request.")

	}

	req.Header.Set("Content-Type", "application/json; charset-utf-8")
	resp, err := client.Do(req)
	if err != nil {
		ErrorLogger.Println("Action: 12, Message: Invalid response.")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	bodyData := string(body)
	split := strings.Split(bodyData, "\"apiKey\":")
	split2 := strings.Split(split[1], "\"")
	resp.Body.Close()
	InfoLogger.Println("Action: 16, Message: Successfully changed api key.s")
	return split2[1]
}

func sendEmail(apiKey string) {
	email := "sarapoparic@gmail.com"

	// Sender data.
	from := "publickeyinfrastructuresomn@hotmail.com"
	password := "PkiSOMN123"

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.office365.com"
	smtpPort := "587"

	// Message.
	fromMessage := fmt.Sprintf("From: <%s>\r\n", "publickeyinfrastructuresomn@hotmail.com")
	toMessage := fmt.Sprintf("To: <%s>\r\n", "sarapoparic@gmail.com")
	subject := "You have connected your account with Dislinkt!\r\n" + apiKey
	body := "Api key to authentificate you are sharing posts is: " + apiKey
	msg := fromMessage + toMessage + subject + "\r\n" + body
	fmt.Println(msg)
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	if err != nil {
		ErrorLogger.Println("Action: 17, Message: Cannot send email.")
		return
	}
	InfoLogger.Println("Action: 18, Message: Email sent successfully")
}

func (ts *CompanyStore) RegisterUser(userReq dto.RequestUser) User {
	user := UserMapper(&userReq)
	ts.db.Create(&user)
	token := ts.generateVerificationToken(user.ID)
	sendRegistrationEmail(userReq.Email, userReq.Name, userReq.Surname, token)
	InfoLogger.Println("Action: 19, Message: User registered.")
	return user
}

func (ts *CompanyStore) generateVerificationToken(userId int) string {
	token := encodeToString(6)
	var user User
	result := ts.db.Find(&user, User{ID: userId})
	if result.RowsAffected > 0 {
		InfoLogger.Println("Action: 20, Message: Token created.")
		user.Token = utils.HashPassword(token)
		user.TokenCreationDate = time.Now()
	}
	ts.db.Save(&user)
	return token
}

func encodeToString(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func sendRegistrationEmail(email, name, surname, token string) {

	// Sender data.
	from := "bezbednostsomn@yahoo.com"
	password := "fcmhbptswmwtphum"

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.mail.yahoo.com"
	smtpPort := "587"

	// Message.
	fromMessage := fmt.Sprintf("From: <%s>\r\n", "bezbednostsomn@yahoo.com")
	toMessage := fmt.Sprintf("To: <%s>\r\n", email)
	subject := "Welcome to Agents App, please verify your registration!\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := "<p>Dear " + name + " " + surname + ",</p>"
	verifyURL := "http://localhost:4200/verification/" + token + "/" + email
	body = body + "<h3><a href=\"" + verifyURL + "\">VERIFY ACCOUNT</a></h3>"
	body = body + "<p>Thank you,<br>Agents App</p>"

	msg := fromMessage + toMessage + subject + mime + body
	fmt.Println(msg)
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	if err != nil {
		ErrorLogger.Println("Action: 17, Message: Cannot send email.")
		return
	}
	InfoLogger.Println("Action: 18, Message: Email sent successfully")
}

func (ts *CompanyStore) PasswordlessLogin(email string) (User, int) {
	var user User
	result := ts.db.Find(&user, User{Email: email})
	if result.RowsAffected == 0 {
		ErrorLogger.Println("Action: 13, Message: User not found.")
		return user, http.StatusNotFound
	}
	token := ts.generateVerificationToken(user.ID)
	sendPasswordlessLoginEmail(user.Email, user.Name, user.Surname, token)
	return user, http.StatusOK
}

func sendPasswordlessLoginEmail(email, name, surname, token string) {

	// Sender data.
	from := "bezbednostsomn@yahoo.com"
	password := "fcmhbptswmwtphum"

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.mail.yahoo.com"
	smtpPort := "587"

	// Message.
	fromMessage := fmt.Sprintf("From: <%s>\r\n", "bezbednostsomn@yahoo.com")
	toMessage := fmt.Sprintf("To: <%s>\r\n", email)
	subject := "Passwordless log in\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := "<p>Dear " + name + " " + surname + ",</p>"
	verifyURL := "http://localhost:4200/verification/" + token
	body = body + "<h3><a href=\"" + verifyURL + "\">LOG IN</a></h3>"
	body = body + "<p>Thank you,<br>Agents App</p>"

	msg := fromMessage + toMessage + subject + mime + body
	fmt.Println(msg)
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	if err != nil {
		ErrorLogger.Println("Action: 17, Message: Cannot send email.")
		return
	}
	InfoLogger.Println("Action: 18, Message: Email sent successfully")
}

func (ts *CompanyStore) AccountRecovery(email string) (User, int) {
	var user User
	result := ts.db.Find(&user, User{Email: email})
	if result.RowsAffected == 0 {
		ErrorLogger.Println("Action: 13, Message: User not found.")
		return user, http.StatusNotFound
	}
	secretKey := os.Getenv("JWT_SECRET_KEY")
	wrapper := JwtWrapper{SecretKey: secretKey, ExpirationHours: 1}
	token, _ := wrapper.GenerateToken(&user)
	sendAccountRecoveryEmail(user.Email, user.Name, user.Surname, token)
	return user, http.StatusOK
}

func sendAccountRecoveryEmail(email, name, surname, token string) {

	// Sender data.
	from := "bezbednostsomn@yahoo.com"
	password := "fcmhbptswmwtphum"

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.mail.yahoo.com"
	smtpPort := "587"

	// Message.
	fromMessage := fmt.Sprintf("From: <%s>\r\n", "bezbednostsomn@yahoo.com")
	toMessage := fmt.Sprintf("To: <%s>\r\n", email)
	subject := "Account recovery\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := "<p>Dear " + name + " " + surname + ",</p>"
	verifyURL := "http://localhost:4200/changePassword/" + token
	body = body + "<h3><a href=\"" + verifyURL + "\">RECOVER ACCOUNT</a></h3>"
	body = body + "<p>Thank you,<br>Agents App</p>"

	msg := fromMessage + toMessage + subject + mime + body
	fmt.Println(msg)
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	if err != nil {
		ErrorLogger.Println("Action: 17, Message: Cannot send email.")
		return
	}
	InfoLogger.Println("Action: 18, Message: Email sent successfully")

}

func (ts *CompanyStore) ChangePassword(id int, password string) (User, int) {
	var user User
	result := ts.db.Find(&user, User{ID: id})
	if result.RowsAffected == 0 {
		ErrorLogger.Println("Action: 13, Message: User not found.")
		return user, http.StatusNotFound
	}
	user.Password = utils.HashPassword(password)
	InfoLogger.Println("Action: 25, Message: User changed password.")
	ts.db.Save(&user)
	return user, http.StatusOK
}

func (ts *CompanyStore) DoesEmailExist(email string) bool {
	var user User
	result := ts.db.Find(&user, User{Email: email})
	if result.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

func (ts *CompanyStore) VerifyAccount(token, email string) (string, int) {
	var user User
	result := ts.db.Find(&user, User{Email: email})
	if result.RowsAffected == 0 {
		ErrorLogger.Println("Action: 13, Message: User not found.")
		return "", http.StatusNotFound
	}
	foundToken := utils.CheckPasswordHash(token, user.Token)
	if !foundToken {
		return "", http.StatusBadRequest
	}
	if time.Since(user.TokenCreationDate).Minutes() <= 20 {
		user.IsVerified = true
		ts.db.Save(&user)
		InfoLogger.Println("Action: 22, Message: User account verified.")
		secretKey := os.Getenv("JWT_SECRET_KEY")
		wrapper := JwtWrapper{SecretKey: secretKey, ExpirationHours: 1}
		token, _ := wrapper.GenerateToken(&user)
		return token, http.StatusOK
	}
	ErrorLogger.Println("Action: 21, Message: Token expired")
	return "", http.StatusOK
}

func (ts *CompanyStore) LoginUser(loginReq dto.RequestLogin) (string, int, bool) {
	var user User
	result := ts.db.Find(&user, User{Username: loginReq.Username})
	if result.RowsAffected == 0 {
		ErrorLogger.Println("Action: 13, Message: User not found.")
		return "", http.StatusNotFound, user.TwoAuth
	}
	match := utils.CheckPasswordHash(loginReq.Password, user.Password)
	if !match {
		ErrorLogger.Println("Action: 23, Message: Wrong email or password.")
		return "", http.StatusNotFound, user.TwoAuth
	}
	if user.IsVerified == false {
		ErrorLogger.Println("Action: 23, Message: User not verified.")
		return "", http.StatusUnauthorized, user.TwoAuth
	}
	secretKey := os.Getenv("JWT_SECRET_KEY")
	wrapper := JwtWrapper{SecretKey: secretKey, ExpirationHours: 1}
	token, _ := wrapper.GenerateToken(&user)
	InfoLogger.Println("Action: 24, Message: Login successful.")
	return token, http.StatusOK, user.TwoAuth
}

func (ts *CompanyStore) GetComment(id int) ([]Comment, error) {
	var comments []Comment
	companyID := strconv.Itoa(id)
	result := ts.db.Find(&comments, "company_id = "+companyID)

	if result.RowsAffected > 0 {
		return comments, nil
	}
	ErrorLogger.Println("Action: 4, Message: Company has no comments.")
	return comments, fmt.Errorf("company with id=%d does not have any new comments", id)
}

func (ts *CompanyStore) Validate(token string) (int, int, string, int) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	wrapper := JwtWrapper{SecretKey: secretKey, ExpirationHours: 1}
	claims, err := wrapper.ValidateToken(token)
	if err != nil {
		ErrorLogger.Println("Action: 21, Message: Token is not valid")
		return http.StatusBadRequest, -1, "", -1
	}
	var user User
	result := ts.db.Find(&user, User{Username: claims.Username})
	if result.RowsAffected == 0 {
		ErrorLogger.Println("Action: 13, Message: User not found.")
		return http.StatusBadRequest, -1, "", -1
	}
	InfoLogger.Println("Action: 22, Message: User validated.")
	return http.StatusOK, claims.Id, claims.Username, claims.Role
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
