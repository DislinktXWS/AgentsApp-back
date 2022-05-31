package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"modules/dto"
	"net/http"
	"strconv"
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

func (ts *CompanyServer) getCompanyHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.GetCompany(id)

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

func (ts *CompanyServer) getJobPositionHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.GetJobPosition(id)

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

func (ts *CompanyServer) getCommentHandler(w http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.GetComment(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
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

func decodeComment(r io.Reader) (*dto.RequestComment, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rc dto.RequestComment
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeJobPosition(r io.Reader) (*dto.RequestJobPosition, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rc dto.RequestJobPosition
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeJobInterview(r io.Reader) (*dto.RequestJobInterview, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rc dto.RequestJobInterview
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeJobSalary(r io.Reader) (*dto.RequestJobSalary, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rc dto.RequestJobSalary
	if err := dec.Decode(&rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func decodeCompany(r io.Reader) (*dto.RequestCompany, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
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
