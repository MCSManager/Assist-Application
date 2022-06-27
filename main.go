package main

import (
	"fmt"
	"os"

	"github.com/MCSManager/Assist-Application/utils"
)

func main() {
	parameters := make([]string, 0)
	for idx, args := range os.Args {
		if idx == 0 || args == "" {
			continue
		}
		if args == "-BIG5" {
			utils.BIG5 = true
			continue
		}
		parameters = append(parameters, args)
	}
	fmt.Printf("ARGS: %v %d\n", parameters, len(parameters))

	// *.exe unzip /root/mcsm.zip /www
	if parameters[0] == "unzip" {
		if len(parameters) != 3 {
			fmt.Printf("Error: Incorrect parameter")
			os.Exit(-1)
		}
		fmt.Printf("Unzip: %s -> %s\n", parameters[1], parameters[2])
		err := utils.Unzip(parameters[1], parameters[2])
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(-1)
		}
	}

	// *.exe zip /www /root/a.zip
	if parameters[0] == "zip" {
		if len(parameters) != 3 {
			fmt.Printf("Error: Incorrect parameter")
			os.Exit(-1)
		}
		fmt.Printf("Zip: %s -> %s\n", parameters[1], parameters[2])
		err := utils.Zip(parameters[1], parameters[2])
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(-1)
		}
	}

	//os.Exit(0)
}
