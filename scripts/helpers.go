package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/blinkops/blink-steampipe/scripts/consts"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func cloneAndInstallModFromPublicRepo(repo string) ([]byte, error) {
	decodedRepo, err := base64.StdEncoding.DecodeString(repo)
	if err != nil {
		return []byte("Failed to decode repo link"), errors.Wrap(err, "decode repo string")
	}

	repoStr := string(decodedRepo)
	modName := extractModName(repoStr)

	modLocation := filepath.Join(consts.SteampipeBasePath, "custom", modName)
	if _, err = os.Stat(modLocation); !os.IsNotExist(err) {
		log.Info("found existing mod, deleting")
		// remove repo if it exists so we can clone it
		if err = os.RemoveAll(modLocation); err != nil {
			return []byte(fmt.Sprintf("Failed to delete exisitng mod %s", modName)),
				errors.Wrap(err, "remove existing mod")
		}
	}

	cmd := exec.Command("git", "clone", repoStr, modLocation)
	if output, err := cmd.CombinedOutput(); err != nil {
		return output, errors.Wrap(err, "git clone mod")
	}

	cmd = exec.Command("steampipe", "mod", "install", "--force")
	cmd.Dir = modLocation
	if output, err := cmd.CombinedOutput(); err != nil {
		return output, errors.Wrap(err, "steampipe mod dependencies install")
	}

	return nil, nil
}

func extractModName(repo string) string {
	splitPath := strings.Split(repo, "/")
	modname := strings.TrimSuffix(splitPath[len(splitPath)-1], ".git")
	return modname
}
