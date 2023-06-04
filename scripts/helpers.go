package main

import (
	"encoding/base64"
	"os"
	"os/exec"
	"strings"

	"github.com/blinkops/blink-steampipe/scripts/consts"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func cloneMod(repo string) error {
	log.Info("starting to handle remote repo")
	decodedRepo, err := base64.StdEncoding.DecodeString(repo)
	if err != nil {
		return errors.Wrap(err, "decode repo string")
	}
	modName := extractModName(string(decodedRepo))
	modLocation := consts.SteampipeBasePath + modName
	if _, err := os.Stat(modLocation); err == nil {
		log.Info("found existing mod, deleting")
		// remove repo if it exists so we can clone it
		if err = os.RemoveAll(modLocation); err != nil {
			return errors.Wrap(err, "remove existing mod")
		}
	}
	cmd := exec.Command("git", "clone", repo, modLocation)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "git clone mod")
	}
	log.Info("cloned repo")
	cmd = exec.Command("steampipe", "mod", "install")
	cmd.Dir = modLocation
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "steampipe mod dependencies install")
	}
	log.Info("installed mod")
	return nil
}

func extractModName(repo string) string {
	splitPath := strings.Split(repo, "/")
	modname := strings.TrimSuffix(splitPath[len(splitPath)-1], ".git")
	return modname
}
