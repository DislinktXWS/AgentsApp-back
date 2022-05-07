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

	rt, err := decodeBody(req.Body)
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

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.UpdateCompany(*rt)
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

func (ts *CompanyServer) getAllCompaniesHandler(w http.ResponseWriter, req *http.Request) {
	allCompanies := ts.store.GetAllCompanies()
	renderJSON(w, allCompanies)
}

func decodeBody(r io.Reader) (*dto.RequestCompany, error) {
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
