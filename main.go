package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	converter "github.com/S-A-RB05/TestManager/converters"
	kubernetes "github.com/S-A-RB05/TestManager/kubernetes"
	"github.com/S-A-RB05/TestManager/messaging"
	"github.com/S-A-RB05/TestManager/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func CreateNewContainer(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func StartContainer(w http.ResponseWriter, r *http.Request) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	if err := cli.ContainerStart(ctx, "0be914e5052a28459884fc535507751c57337e7a09a413c08c32b578b984b000", types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func ExecuteCmd(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: cmd")

	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		panic(err)
	}

	respIdExecCreate, err := cli.ContainerExecCreate(context.Background(), "e07a8c609ace20685d357345e40df43868ce5ee628744c79076ec58d42e1eed2", types.ExecConfig{
		User:       "root",
		Privileged: true,
		Cmd: []string{
			"sh", "-c", "xvfb-run wine terminal.exe  /config:'Report\\confighoi.ini'",
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	response, err := cli.ContainerExecAttach(context.Background(), respIdExecCreate.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}
	defer response.Close()

	data, _ := ioutil.ReadAll(response.Reader)
	fmt.Println(string(data))

	fmt.Println("Executed")
}

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

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/create", CreateNewContainer)
	myRouter.HandleFunc("/start", StartContainer)
	myRouter.HandleFunc("/cmd", ExecuteCmd)
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
