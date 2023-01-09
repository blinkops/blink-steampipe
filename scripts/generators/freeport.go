package generators

import (
	"fmt"
	"github.com/blinkops/blink-steampipe/scripts/consts"
	"github.com/phayes/freeport"
	"os"
	"strconv"
	"strings"
)

const steampipeDBConfigurationFile = consts.SteampipeSpcConfigurationPath + "db.spc"

type FreePortGenerator struct{}

func (gen FreePortGenerator) Generate() error {
	port, err := freeport.GetFreePort()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(steampipeDBConfigurationFile)
	if err != nil {
		return fmt.Errorf("unable to prepare db port configuration: %w", err)
	}

	dataAsString := string(data)
	dataAsString = strings.ReplaceAll(dataAsString, "{{FREE_PORT}}", strconv.Itoa(port))

	if err = os.WriteFile(steampipeDBConfigurationFile, []byte(dataAsString), 0o600); err != nil {
		return fmt.Errorf("unable to prepare db port config file: %w", err)
	}

	return nil
}
