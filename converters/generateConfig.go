package converter

import (
	"fmt"
	"os"

	"github.com/Jeffail/gabs"
)

func Test() {
	data := []byte(`{
		"employees":{
		   "protected":false,
		   "address":{
			  "street":"22 Saint-Lazare",
			  "postalCode":"75003",
			  "city":"Paris",
			  "countryCode":"FRA",
			  "country":"France"
		   },
		   "employee":[
			  {
				 "id":1,
				 "first_name":"Jeanette",
				 "last_name":"Penddreth"
			  },
			  {
				 "id":2,
				 "firstName":"Giavani",
				 "lastName":"Frediani"
			  }
		   ]
		}
	 }`)

	jsonParsed, err := gabs.ParseJSON(data)
	if err != nil {
		panic(err)
	}

	// Search JSON
	fmt.Println("Get value of Protected:\t", jsonParsed.Path("employees.protected").Data())
	fmt.Println("Get value of Country:\t", jsonParsed.Search("employees", "address", "country").Data())
	fmt.Println("ID of first employee:\t", jsonParsed.Path("employees.employee.0.id").String())
	fmt.Println("Check Country Exists:\t", jsonParsed.Exists("employees", "address", "countryCode"))

	// Iterating employee array
	for _, child := range jsonParsed.Search("employees", "employee").Children() {
		fmt.Println(child.Data())
	}

	// Use index in your search
	for _, child := range jsonParsed.Search("employees", "employee", "0").Children() {
		fmt.Println(child.Data())
	}
}

func GenerateConfig() {
		// Create a new file to write the configuration settings to
		file, err := os.Create("config.ini")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
	
		// Write the configuration settings to the file
		_, err = file.WriteString("[Common]\n")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = file.WriteString("Login=123456\n")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = file.WriteString("Password=password\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		// Write the configuration settings to the file
		_, err = file.WriteString("[Tester]\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("Expert=Examples\\MACD\\MACD Sample\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("Symbol=EURUSD\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("Period=H1\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("Deposit=10000\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("Leverage=1:100\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("Model=0\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("ExecutionMode=1\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("Optimization=0\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("OptimizationCriterion=0\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("FromDate=2011.01.01\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("ToDate=2011.04.01\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("Report=test_macd\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("ReplaceReport=1\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		_, err = file.WriteString("ShutdownTerminal=0\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		if err := file.Sync(); err != nil {
			panic(err)
		}
	
		// Print a success message
		fmt.Println("config.ini file generated successfully")
}
