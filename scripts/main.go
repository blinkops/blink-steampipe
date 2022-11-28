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
		response_wrapper.HandleResponse("", fmt.Sprintf("set up logger: %v", err.Error()), true)
		os.Exit(1)
	}

	for _, credentialGenerator := range generators.Generators {
		if err := credentialGenerator.Generate(); err != nil {
			logrus.Errorf("failed generate credentials: %v", err)
			response_wrapper.HandleResponse("", logger.GetLogs(), true)
			os.Exit(1)
		}
	}

	cmdName := os.Args[1]
	if cmdName == "" {
		logrus.Error("no command provided")
		response_wrapper.HandleResponse("", logger.GetLogs(), true)
		os.Exit(1)
	}
	cmdArgs := os.Args[2:]

	cmd := exec.Command(cmdName, cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("execute command: %v", err)
		response_wrapper.HandleResponse(string(output), logger.GetLogs(), true)
		os.Exit(1)
	}

	response_wrapper.HandleResponse(string(output), logger.GetLogs(), false)
	os.Exit(0)
}
