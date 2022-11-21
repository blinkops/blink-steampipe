package main

import (
	"fmt"
	"github.com/blinkops/blink-steampipe/internal/logger"
	"github.com/blinkops/blink-steampipe/internal/response_wrapper"
	"github.com/blinkops/blink-steampipe/scripts/generators"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func main() {
	if err := logger.SetUpLogger(); err != nil {
		fmt.Println("main - SetUpLogger error " + err.Error())
		response_wrapper.HandleResponse("", err.Error())
		os.Exit(1)
	}

	for idx, credentialGenerator := range generators.Generators {
		fmt.Printf("main - credentialGenerator: idx = %d", idx)
		if err := credentialGenerator.Generate(); err != nil {
			fmt.Println("main - credential Generate error" + err.Error())
			logrus.Error(err)
			response_wrapper.HandleResponse("", logger.GetLogs())
			os.Exit(1)
		}
	}

	cmdArg1 := os.Args[0]
	fmt.Println("main - args values: arg1 = " + cmdArg1)

	cmdArg2 := os.Args[1:]
	fmt.Printf("main - args values: arg2 = %s", cmdArg2)

	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		logrus.Error(err)
		response_wrapper.HandleResponse("", logger.GetLogs())
		os.Exit(1)
	}

	response_wrapper.HandleResponse(string(output), logger.GetLogs())
	os.Exit(0)
}
