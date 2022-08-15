package main

import (
	"os"

	"github.com/blinkops/blink-steampipe/build/scripts/generators"
	"github.com/sirupsen/logrus"
)

func main() {
	for _, credentialGenerator := range generators.Generators {
		if err := credentialGenerator.Generate(); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
