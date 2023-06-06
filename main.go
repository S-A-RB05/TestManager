package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	converter "github.com/S-A-RB05/TestManager/converters"
	kubernetes "github.com/S-A-RB05/TestManager/kubernetes"
	"github.com/S-A-RB05/TestManager/messaging"
	"github.com/S-A-RB05/TestManager/models"
	"github.com/gorilla/mux"
)

func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	fmt.Println("Updating config")
	// parse the request body into a Strategy struct
	var data converter.Data
	err := json.NewDecoder(body).Decode(&data)
	fmt.Println(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	var strat = readSingleStrat(data.ID)

	converter.GenerateConfig(data, strat)
}

// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("ok")

		// Next
		next.ServeHTTP(w, r)
		//return
	})

}

func runTest(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	fmt.Println("Updating config")
	// parse the request body into a Strategy struct
	var data converter.Data
	err := json.NewDecoder(body).Decode(&data)
	fmt.Println(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	var strat = readSingleStrat(data.ID)

	//Generate config file
	bytes := converter.GenerateConfig(data, strat)

	var test models.Test
	//TODO: extract userID
	test.StratId = data.ID

	//Create job (pod)
	jobId, jobError := kubernetes.CreateJob("development")
	if jobError != nil {
		fmt.Print(jobError)
	}

	// Send config to pod
	messaging.ProduceMessage(bytes, "mt5_test")

	// Insert test into DB
	test.Id = jobId
	var testId = insertTest(test)

	fmt.Fprintf(w, testId)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.Use(CORS)

	myRouter.HandleFunc("/updateconfig", runTest)

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go handleRequests()
	go messaging.ConsumeMessage("q.syncStrat", insertStrat)
	<-stop
}
