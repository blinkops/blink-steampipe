package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/blinkops/blink-steampipe/scripts/consts"
	"github.com/pkg/errors"
)

func cloneMod(repo string) error {
	modName := extractModName(repo)
	modLocation := consts.SteampipeBasePath + modName
	if _, err := os.Stat(modLocation); err == nil {
		// remove repo if it exists so we can clone it
		if err = os.RemoveAll(modLocation); err != nil {
			return errors.Wrap(err, "remove existing mod")
		}
	}
	cmd := exec.Command("git", "clone", repo, modLocation)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "git clone mod")
	}
	cmd = exec.Command("steampipe", "mod", "install")
	cmd.Dir = modLocation
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "steampipe mod dependencies install")
	}
	return nil
}

func extractModName(repo string) string {
	splitPath := strings.Split(repo, "/")
	modname := strings.TrimSuffix(splitPath[len(splitPath)-1], ".git")
	return modname
}
