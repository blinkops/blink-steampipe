package response_wrapper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blinkops/blink-steampipe/scripts/consts"
	uuid "github.com/satori/go.uuid"
)

const (
	queryWarningMessage = "Warning: executeQueries: "
	queryErrorMessage   = "Error: "
	rpcErrorCode        = "rpc error: code ="
)

const (
	generalErrorMessage = "failed to initialize plugin"
)

const responseWrapperFormat = `{"output":"%s", "log": "%s", "is_error": "%v"}`

type ResponseWrapper struct {
	Log     string `json:"log"`
	Output  string `json:"output"`
	IsError bool   `json:"is_error"`
}

func DebugModeEnabled() bool {
	debugEnv := os.Getenv("BLINK_STEAMPIPE_DEBUG")
	if strings.ToLower(debugEnv) == `"true"` {
		return true
	}

	return false
}

func HandleResponse(output, log string, exitWithError bool) {
	resp := &ResponseWrapper{
		Log: log,
	}

	result := strings.TrimSpace(output)
	isError := false

	if !DebugModeEnabled() {
		// only show friendly errors in operational mode
		// so dev/cs can investigate issues if needed
		result, isError = formatErrorMessage(result)
	}

	if result == "" {
		result = generalErrorMessage
		isError = true
	}

	resp.Output = result
	resp.IsError = isError || exitWithError
	handleFileIdentifierResponseIfRequired(resp)

	marshaledResponse, err := json.Marshal(resp)
	if err != nil {
		updatedLog := fmt.Sprintf("%s\nfailed to marshal response: %v", log, err.Error())
		fmt.Printf(responseWrapperFormat, resp.Output, updatedLog, resp.IsError)
		return
	}

	fmt.Println(string(marshaledResponse))
}

func handleFileIdentifierResponseIfRequired(resp *ResponseWrapper) {
	fileIdentifierParentDir := os.Getenv(consts.FileIdentifierParentDirEnvVar)
	if fileIdentifierParentDir == "" {
		return
	}

	fileIdentifier := os.Getenv(consts.FileIdentifierEnvVar)
	if fileIdentifier == "" {
		return
	}

	fileIdentifierPath := filepath.Join(fileIdentifierParentDir, fileIdentifier)
	resp.Log = fmt.Sprintf("%s\nwrite to output to file '%s', outout:\n%s", resp.Log, fileIdentifierPath, resp.Output)
	if err := os.WriteFile(fileIdentifierPath, []byte(resp.Output), 0o644); err != nil {
		resp.IsError = true
		resp.Output = fmt.Sprintf("failed to write file: %v\n", err.Error())
		return
	}

	overrideFormat := consts.FileOutputOverrideFormat
	if resp.IsError {
		overrideFormat = consts.FileOutputOverrideOnErrorFormat
	}
	resp.Output = fmt.Sprintf(overrideFormat, fileIdentifier)
}

// formatErrorMessage format the error message to be cleaner
func formatErrorMessage(result string) (msg string, isError bool) {
	if strings.Contains(result, queryWarningMessage) {
		result = strings.TrimSpace(strings.ReplaceAll(result, queryWarningMessage, ""))
		// older versions of the runner do not correctly report an error when steampipe returns a nonzero error,
		// so we have to parse the error message to determine if it is an error anyway
		isError = true
	}
	if strings.HasPrefix(result, queryErrorMessage) {
		result = strings.TrimSpace(strings.ReplaceAll(result, queryErrorMessage, ""))
		isError = true
	}
	if strings.Contains(result, rpcErrorCode) {
		result = strings.TrimSpace(strings.ReplaceAll(result, rpcErrorCode, ""))
		if firstIndex := strings.Index(result, "="); firstIndex != -1 {
			result = strings.TrimSpace(result[firstIndex+1:])
		}

		isError = true

		lowerMessage := strings.ToLower(result)
		if strings.Contains(lowerMessage, "hydrate function") {
			return "Connection doesn't have permissions to run specified query", isError
		}

		if strings.Contains(lowerMessage, "invalid regions") {
			return "Connection has unsupported regions", isError
		}

		if strings.Contains(lowerMessage, "connection") || strings.Contains(lowerMessage, "credential") || strings.Contains(lowerMessage, "no such host") {
			return "Invalid connection was provided", isError
		}

		if strings.Contains(lowerMessage, "sqlstate") {
			sqlStateMatch := regexp.MustCompile(`\(SQLSTATE .*\)`)
			result = sqlStateMatch.ReplaceAllString(result, "")
			return result, isError
		}

		identifier := uuid.NewV4().String()
		return fmt.Sprintf("Invalid error was received, Identifier: %s", identifier), isError
	}

	return result, isError
}
