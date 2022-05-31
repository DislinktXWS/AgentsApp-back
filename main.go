package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)
	server, err := NewCompanyServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer server.CloseDB()

	//COMPANY HANDLERS
	router.HandleFunc("/company", server.createCompanyHandler).Methods("POST")
	router.HandleFunc("/company", server.updateCompanyHandler).Methods("PUT")
	router.HandleFunc("/company", server.getAllCompaniesHandler).Methods("GET")
	router.HandleFunc("/company/{id:[0-9a-zA-Z]+}/", server.getCompanyHandler).Methods("GET")
	router.HandleFunc("/company/accept", server.acceptCompanyHandler).Methods("PUT")

	//JOB SALARY HANDLERS
	router.HandleFunc("/jobSalary", server.createJobSalaryHandler).Methods("POST")
	router.HandleFunc("/jobSalary/{id:[0-9a-zA-Z]+}", server.getJobSalaryHandler).Methods("GET")
	router.HandleFunc("/jobSalary/{id:[0-9a-zA-Z]+}", server.deleteJobSalaryHandler).Methods("DELETE")

	//JOB INTERVIEW HANDLERS
	router.HandleFunc("/jobInterview", server.createJobInterviewHandler).Methods("POST")
	router.HandleFunc("/jobInterview/{id:[0-9a-zA-Z]+}", server.getJobInterviewHandler).Methods("GET")

	//JOB POSITION HANDLERS
	router.HandleFunc("/jobPosition", server.createJobPositionHandler).Methods("POST")
	router.HandleFunc("/jobPosition/{id:[0-9a-zA-Z]+}", server.getJobPositionHandler).Methods("GET")

	//COMMENT HANDLER
	router.HandleFunc("/comment", server.createCommentHandler).Methods("POST")
	router.HandleFunc("/comment/{id:[0-9a-zA-Z]+}", server.getCommentHandler).Methods("GET")

	srv := &http.Server{Addr: "0.0.0.0:9000", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}
