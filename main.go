package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	ch := handlers.CORS(handlers.AllowedOrigins([]string{"http://localhost:4200", "http://localhost:4201"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Authorization", "Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"}),
	)
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
	router.HandleFunc("/companyRequests", server.getCompanyRequestsHandler).Methods("GET")
	router.HandleFunc("/company/{id:[0-9a-zA-Z]+}/", server.getCompanyByIDHandler).Methods("GET")
	router.HandleFunc("/companyByOwner/{id:[0-9a-zA-Z]+}/", server.getOwnersCompaniesHandler).Methods("GET")
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
	router.HandleFunc("/jobPosition", server.getAllJobPositionHandler).Methods("GET")
	router.HandleFunc("/jobPosition/{id:[0-9a-zA-Z]+}", server.getJobPositionHandler).Methods("GET")
	router.HandleFunc("/shareOnDislinkt/{apiKey}", server.shareJobPosition).Methods("POST")
	router.HandleFunc("/isJobPositionShared/{id}", server.isJobPositionShared).Methods("GET")

	//COMMENT HANDLER
	router.HandleFunc("/comment", server.createCommentHandler).Methods("POST")
	router.HandleFunc("/comment/{id:[0-9a-zA-Z]+}", server.getCommentHandler).Methods("GET")

	//CONNECTION WITH DISLINKT HANDLER
	router.HandleFunc("/connectWithDislinkt/{username}/{id}", server.connectWithDislinkt).Methods("PUT")
	router.HandleFunc("/isConnected/{id}", server.isConnectedHandler).Methods("GET")

	//AUTH HANDLER
	router.HandleFunc("/registration", server.registerHandler).Methods("POST")
	router.HandleFunc("/login", server.loginHandler).Methods("POST")
	router.HandleFunc("/validate/{token}", server.validateHandler).Methods("GET")
	router.HandleFunc("/verifyAccount/{token}", server.verifyAccountHandler).Methods("GET")

	srv := &http.Server{Addr: "0.0.0.0:9000", Handler: ch(router)}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
		/*if err := http.ListenAndServeTLS("0.0.0.0:9000", "D:\\localhost.crt", "D:\\localhost.key", ch(router)); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}*/
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
