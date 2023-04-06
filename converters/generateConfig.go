package converter

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Variables struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefaultValue string `json:"defaultValue"`
	Start        int    `json:"start"`
	End          int    `json:"end"`
	Step         int    `json:"step"`
}

func Test() {
	jsonString := `[{
		"name": "testvar",
		"type": "int",
		"defaultValue": "1",
		"start": 1,
		"end": 1,
		"step": 1
	},{
		"name": "testvar",
		"type": "int",
		"defaultValue": "1",
		"start": 1,
		"end": 1,
		"step": 1
	},
	{
		"name": "testbool",
		"type": "bool",
		"defaultValue": "false",
		"boolValue": false
	}]`

	var vars []Variables
	err := json.Unmarshal([]byte(jsonString), &vars)
	if err != nil {
		panic(err)
	}

	for _, variables := range vars {
		fmt.Printf(
			"Name: %s, Type: %s, DefaultValue: %s, Start: %d, End: %d, Step: %d\n",
			variables.Name, variables.Type, variables.DefaultValue, variables.Start, variables.End, variables.Step,
		)
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

	// Write the tester inputs
	_, err = file.WriteString("[TesterInputs]\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonString := `[{
		"name": "testvar",
		"type": "int",
		"defaultValue": "1",
		"start": 1,
		"end": 1,
		"step": 1
	},{
		"name": "testvar",
		"type": "int",
		"defaultValue": "1",
		"start": 1,
		"end": 1,
		"step": 1
	},
	{
		"name": "testbool",
		"type": "bool",
		"defaultValue": "false",
		"boolValue": false
	}]`

	var vars []Variables
	KoenPoenTown := json.Unmarshal([]byte(jsonString), &vars)
	if KoenPoenTown != nil {
		panic(KoenPoenTown)
	}

	for _, variables := range vars {

		_, err = file.WriteString(variables.Name+"="+variables.DefaultValue+"||"+strconv.Itoa(variables.Start)+"||"+strconv.Itoa(variables.Step)+"||"+strconv.Itoa(variables.End)+"||N\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	
		fmt.Printf(
			"Name: %s, Type: %s, DefaultValue: %s, Start: %d, End: %d, Step: %d\n",
			variables.Name, variables.Type, variables.DefaultValue, variables.Start, variables.End, variables.Step,
		)
	}

	// Print a success message
	fmt.Println("config.ini file generated successfully")
}
