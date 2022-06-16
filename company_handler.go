package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"modules/dto"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (ts *CompanyServer) createCompanyHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeCompany(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateCompany(*rt)
	renderJSON(w, dto.ResponseId{Id: id})
}

func (ts *CompanyServer) updateCompanyHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeCompany(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.UpdateCompany(*rt)
	renderJSON(w, dto.ResponseId{Id: id})
}

func (ts *CompanyServer) acceptCompanyHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeAcceptRequest(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.AcceptCompany(*rt)
	renderJSON(w, dto.ResponseId{Id: id})
}

func (ts *CompanyServer) getOwnersCompaniesHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.GetOwnersCompanies(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *CompanyServer) getCompanyByIDHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.GetCompanyById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *CompanyServer) getJobSalaryHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.GetJobSalary(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *CompanyServer) deleteJobSalaryHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := ts.store.DeleteJobSalary(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (ts *CompanyServer) getAllCompaniesHandler(w http.ResponseWriter, req *http.Request) {
	allCompanies := ts.store.GetAllCompanies()
	renderJSON(w, allCompanies)
}

func (ts *CompanyServer) getCompanyRequestsHandler(w http.ResponseWriter, req *http.Request) {
	allCompanies := ts.store.GetCompanyRequests()
	renderJSON(w, allCompanies)
}

func (ts *CompanyServer) getAllJobPositionHandler(w http.ResponseWriter, req *http.Request) {
	allJobPositions := ts.store.GetAllJobPositions()
	renderJSON(w, allJobPositions)
}

func (ts *CompanyServer) createJobSalaryHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeJobSalary(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateJobSalary(*rt)
	renderJSON(w, dto.ResponseId{Id: id})
}

func (ts *CompanyServer) createJobInterviewHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeJobInterview(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateJobInterview(*rt)
	renderJSON(w, dto.ResponseId{Id: id})
}

func (ts *CompanyServer) getJobInterviewHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.GetJobInterview(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *CompanyServer) createJobPositionHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeJobPosition(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateJobPosition(*rt)
	renderJSON(w, dto.ResponseId{Id: id})
}

func (ts *CompanyServer) shareJobPosition(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	apiKey, _ := mux.Vars(req)["apiKey"]
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeJobPosition(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(rt)
	result, id := ts.store.ShareJobPosition(*rt, apiKey)

	if result {
		ts.store.SetJobPositionIsShared(id)
	}

	renderJSON(w, "Successfully shared")
}

func (ts *CompanyServer) getJobPositionHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.GetJobPosition(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *CompanyServer) isConnectedHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.IsConnected(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *CompanyServer) isJobPositionShared(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.IsJobPositionShared(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *CompanyServer) createCommentHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeComment(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateComment(*rt)
	renderJSON(w, dto.ResponseId{Id: id})
}

func (ts *CompanyServer) connectWithDislinkt(w http.ResponseWriter, req *http.Request) {
	username, _ := mux.Vars(req)["username"]
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	apiKey := ts.store.ConnectWithDislinkt(username, id)
	response, err := json.Marshal(apiKey)
	if apiKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (ts *CompanyServer) registerHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeUser(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := ts.store.RegisterUser(*rt)
	renderJSON(w, dto.ResponseId{Id: user.ID})
}

func (ts *CompanyServer) verifyAccountHandler(w http.ResponseWriter, req *http.Request) {
	token, _ := mux.Vars(req)["token"]
	jwtToken, status := ts.store.VerifyAccount(token)

	if status != 200 {
		return
	}

	renderJSON(w, dto.ResponseLogin{Token: jwtToken})
}

func (ts *CompanyServer) loginHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeLogin(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, status := ts.store.LoginUser(*rt)
	if status != 200 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, dto.ResponseLogin{Token: token})
}

func (ts *CompanyServer) passwordlessLoginHandler(w http.ResponseWriter, req *http.Request) {
	email, _ := mux.Vars(req)["email"]

	user, status := ts.store.PasswordlessLogin(email)
	if status != 200 {
		return
	}
	renderJSON(w, dto.ResponseId{Id: user.ID})
}

func (ts *CompanyServer) accountRecoveryHandler(w http.ResponseWriter, req *http.Request) {
	email, _ := mux.Vars(req)["email"]

	user, status := ts.store.AccountRecovery(email)
	if status != 200 {
		return
	}
	renderJSON(w, dto.ResponseId{Id: user.ID})
}

func (ts *CompanyServer) changePasswordHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := mux.Vars(req)["id"]
	userId, _ := strconv.Atoi(id)

	log.Println(req.Body)
	password, err := decodePassword(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, status := ts.store.ChangePassword(userId, password.Password)
	if status != 200 {
		return
	}
	renderJSON(w, dto.ResponseId{Id: user.ID})
}

func (ts *CompanyServer) getCommentHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.GetComment(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

func (ts *CompanyServer) validateHandler(w http.ResponseWriter, req *http.Request) {
	token, _ := mux.Vars(req)["token"]
	status, id, username, role := ts.store.Validate(token)
	if status != 200 {
		http.Error(w, string(status), http.StatusNotFound)
		return
	}
	renderJSON(w, dto.ResponseValidate{ID: id, Username: username, Role: role})

}

func decodeAcceptRequest(r io.Reader) (*dto.RequestAcceptCompany, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rc dto.RequestAcceptCompany
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeUser(r io.Reader) (*dto.RequestUser, error) {
	dec := json.NewDecoder(r)
	var rc dto.RequestUser
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodePassword(r io.Reader) (*dto.Password, error) {
	dec := json.NewDecoder(r)
	log.Println(r)
	var rc dto.Password
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeLogin(r io.Reader) (*dto.RequestLogin, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rc dto.RequestLogin
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeComment(r io.Reader) (*dto.RequestComment, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rc dto.RequestComment
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeConnection(r io.Reader) (*dto.Connection, error) {
	dec := json.NewDecoder(r)
	var rc dto.Connection
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeJobPosition(r io.Reader) (*dto.RequestJobPosition, error) {
	dec := json.NewDecoder(r)
	var rc dto.RequestJobPosition
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeJobInterview(r io.Reader) (*dto.RequestJobInterview, error) {
	dec := json.NewDecoder(r)
	var rc dto.RequestJobInterview
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeJobSalary(r io.Reader) (*dto.RequestJobSalary, error) {
	dec := json.NewDecoder(r)
	var rc dto.RequestJobSalary
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeCompany(r io.Reader) (*dto.RequestCompany, error) {
	dec := json.NewDecoder(r)
	var rc dto.RequestCompany
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		return
	}
}
