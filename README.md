# TestManager Documentation

## API endpoints

### Endpoint: /updateconfig
Method: POST
- Description: Updates the configuration based on the request body.
- Request Body: Expects a JSON payload containing the configuration data.
- Response: Returns the ID of the generated test.
Example Usage:

```
curl -X POST -H "Content-Type: application/json" -d '
{
  "id": "645a46d8d9335b1665d85ef2",
  "variables": [
    {
      "name": "Inp_Signal_MACD_PeriodFast",
      "type": "int",
      "defaultValue": "12",
      "start": 12,
      "end": 12,
      "step": 12
    },
    {
      "name": "Inp_Signal_MACD_PeriodSlow",
      "type": "int",
      "defaultValue": "24",
      "start": 24,
      "end": 24,
      "step": 24
    },
    {
      "name": "Inp_Signal_MACD_PeriodSignal",
      "type": "int",
      "defaultValue": "9",
      "start": 9,
      "end": 9,
      "step": 9
    },
    {
      "name": "Inp_Signal_MACD_TakeProfit",
      "type": "int",
      "defaultValue": "50",
      "start": 50,
      "end": 50,
      "step": 50
    },
    {
      "name": "Inp_Signal_MACD_StopLoss",
      "type": "int",
      "defaultValue": "20",
      "start": 20,
      "end": 20,
      "step": 20
    }
  ]
}
' http://localhost:8081/updateconfig
```


## Running application

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