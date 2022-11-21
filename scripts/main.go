package main

import (
	"github.com/blinkops/blink-steampipe/internal/logger"
	"github.com/blinkops/blink-steampipe/internal/response_wrapper"
	"github.com/blinkops/blink-steampipe/scripts/generators"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func main() {
	if err := logger.SetUpLogger(); err != nil {
		response_wrapper.HandleResponse("", err.Error())
		os.Exit(1)
	}

	for _, credentialGenerator := range generators.Generators {
		if err := credentialGenerator.Generate(); err != nil {
			logrus.Error(err)
			response_wrapper.HandleResponse("", logger.GetLogs())
			os.Exit(1)
		}
	}

	cmdName := os.Args[1]
	if cmdName == "" {
		logrus.Error("empty command supplied")
		response_wrapper.HandleResponse("", logger.GetLogs())
		os.Exit(1)
	}
	cmdArgs := os.Args[2:]

	cmd := exec.Command(cmdName, cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Error(err)
		response_wrapper.HandleResponse(string(output), logger.GetLogs())
		os.Exit(1)
	}

	response_wrapper.HandleResponse(string(output), logger.GetLogs())
	os.Exit(0)
}
