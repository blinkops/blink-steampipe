package main

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/blinkops/blink-steampipe/scripts/consts"
	"github.com/pkg/errors"
)

func cloneMod(repo string) error {
	modName := extractModName(repo)
	modLocation := filepath.Join(consts.SteampipeBasePath, modName)
	cmd := exec.Command("git", "clone", repo, modLocation)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, cmd.String())
	}
	cmd = exec.Command("cd", modLocation)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, cmd.String())
	}
	cmd = exec.Command("steampipe", "mod", "install")
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, cmd.String())
	}
	return nil
}

func extractModName(repo string) string {
	splitPath := strings.Split(repo, "/")
	modname := strings.TrimSuffix(splitPath[len(splitPath)-1], ".git")
	return modname
}
