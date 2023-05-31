package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/blinkops/blink-steampipe/internal/logger"
	"github.com/blinkops/blink-steampipe/internal/response_wrapper"
	"github.com/blinkops/blink-steampipe/scripts/consts"
	"github.com/blinkops/blink-steampipe/scripts/generators"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := logger.SetUpLogger(); err != nil {
		response_wrapper.HandleResponse("", fmt.Sprintf("set up logger: %v", err.Error()), "", true)
		os.Exit(0)
	}

	for _, credentialGenerator := range generators.Generators {
		if err := credentialGenerator.Generate(); err != nil {
			log.Errorf("Failed to generate credentials: %v", err)
			response_wrapper.HandleResponse("", logger.GetLogs(), "", true)
			os.Exit(0)
		}
	}

	var cmdName, action string
	if len(os.Args) > 2 {
		cmdName = os.Args[1]
		action = os.Args[2]
	}

	if cmdName == "" {
		log.Error("No command provided")
		response_wrapper.HandleResponse("", logger.GetLogs(), action, true)
		os.Exit(0)
	}
	cmdArgs := os.Args[2:]

	// if a mod location was provided,
	// we clone the repo to the corresponding location mentioned in queryV2PreInvoke in controller
	if modLocation := os.Getenv(consts.SteampipeReportCustomModLocationEnvVar); modLocation != "" {
		if err := cloneMod(modLocation); err != nil {
			log.Errorf("Failed to load the custom Steampipe mod from %s. Please ensure the provided repository is publicly available.", modLocation)
			response_wrapper.HandleResponse("", logger.GetLogs(), action, true)
			os.Exit(0)
		}
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	output, err := cmd.CombinedOutput()

	// some steampipe benchmark ("check") may return an error code but return a result.
	// in such a case, we don't want the entire report to fail and display the result.
	if err != nil && (action != consts.CommandCheck || len(output) == 0) {
		log.Errorf("Failed to execute command: %v", err)
		response_wrapper.HandleResponse(string(output), logger.GetLogs(), action, true)
		os.Exit(0)
	}

	response_wrapper.HandleResponse(string(output), logger.GetLogs(), action, false)
	os.Exit(0)
}

func cloneMod(repo string) error {
	cmd := exec.Command("git", "clone", repo, filepath.Join(consts.SteampipeBasePath, consts.SteampipeCustomModAlias), "--quiet")
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, cmd.String())
	}
	return nil
}
