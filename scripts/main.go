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
