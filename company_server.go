package main

import (
	ps "modules/company_store"
)

type CompanyServer struct {
	store *ps.CompanyStore
}

func NewCompanyServer() (*CompanyServer, error) {
	store, err := ps.New()
	if err != nil {
		return nil, err
	}
	return &CompanyServer{
		store: store,
	}, nil
}

func (s *CompanyServer) CloseDB() error {
	return s.store.Close()
}
