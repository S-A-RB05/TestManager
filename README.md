# TestManager Documentation

## API endpoints

### Endpoint: /updateconfig
Method: POST
- Description: Updates the configuration based on the request body.
- Request Body: Expects a JSON payload containing the configuration data.
- Response: Returns the ID of the generated test.
Example Usage:

```
curl -X POST -H "Content-Type: application/json" -d '{"ID": "exampleID", "data": "exampleData"}' http://localhost:8081/updateconfig
```


## Instructions for deployment

Start the server by running the Go application in the terminal.

```
go run .
``` 
Remember to adjust the host and port (http.ListenAndServe(":8081", myRouter)) based on your deployment environment or preferred configuration.

## Code docs

Please note that this documentation assumes that you are familiar with the Gorilla Mux router, JSON encoding/decoding, and other related Go concepts.

` Main() `
```
In the main function, a channel is created to receive OS signals. 
HTTP request handling and message consuming are started in separate goroutines. 
The program waits for an OS signal (such as interrupt or termination) to stop execution.
```

`HandleRequest() `
```
handleRequests sets up the router using the gorilla/mux package, 
defines the routes and their corresponding handlers, 
and starts the HTTP server to listen on port 8081.
```

`RunTest() ` 
```
runTest is an HTTP handler function that receives a request body containing data, 
generates a configuration file based on the data, 
creates a job (pod) on a Kubernetes cluster, 
sends the generated config to the pod using messaging, 
and inserts the test information into the database.
```

`Cors() ` 
```
CORS is a middleware function that adds CORS (Cross-Origin Resource Sharing) headers to the HTTP response. 
It allows cross-origin requests from any origin and sets allowed methods and headers.
```

`UpdateConfig()`  
```
UpdateConfig is an HTTP handler function that receives a request body containing data and 
updates the configuration based on the provided data.
```