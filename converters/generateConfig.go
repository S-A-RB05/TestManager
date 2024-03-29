package converter

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/S-A-RB05/TestManager/models"
)

type Variable struct {
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	DefaultValue string  `json:"defaultValue"`
	Start        float32 `json:"start,omitempty"`
	End          float32 `json:"end,omitempty"`
	Step         float32 `json:"step,omitempty"`
	BoolValue    bool    `json:"boolValue,omitempty"`
}

type Data struct {
	ID        string     `json:"id"`
	Variables []Variable `json:"variables"`
}

func GenerateConfig(data Data, strat models.StrategyRequest) (config []byte) {

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
	_, err = file.WriteString("Login=64624\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString("Password=2oqipkri\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = file.WriteString("Server=AlbaBrokers-Demo\n")
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

	_, err = file.WriteString("Expert=/Advisors/" + strat.Name + "\n")
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

	_, err = file.WriteString("Report=./Report/tester_report.htm\n")
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

	for _, v := range data.Variables {
		fmt.Printf("Name: %s, Type: %s, DefaultValue: %s\n", v.Name, v.Type, v.DefaultValue)
		if v.Type == "int" {
			_, err = file.WriteString(v.Name + "=" + v.DefaultValue + "||" + strconv.FormatFloat(float64(v.Start), 'f', -1, 32) + "||" + strconv.FormatFloat(float64(v.Step), 'f', -1, 32) + "||" + strconv.FormatFloat(float64(v.End), 'f', -1, 32) + "||N\n")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Start: %d, End: %d, Step: %d\n", v.Start, v.End, v.Step)
		} else if v.Type == "bool" {
			_, err = file.WriteString(v.Name + "=" + v.DefaultValue + "||" + strconv.FormatBool(v.BoolValue) + "||0||" + strconv.FormatBool(!v.BoolValue) + "||N\n")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("BoolValue: %t\n", v.BoolValue)
		}
	}

	// Read the file contents into a byte slice
	fileBytes, err := ioutil.ReadFile("config.ini")
	if err != nil {
		log.Fatal(err)
	}

	// Print a success message
	fmt.Println("config.ini file generated successfully for strategy with id: " + data.ID)
	return fileBytes
}
