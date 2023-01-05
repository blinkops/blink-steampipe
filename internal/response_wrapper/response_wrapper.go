package response_wrapper

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"regexp"
	"strings"
)

const (
	queryWarningMessage = "Warning: executeQueries: "
	queryErrorMessage   = "Error: "
	rpcErrorCode        = "rpc error: code ="
)

const generalErrorMessage = "failed to initialize plugin"

const responseWrapperFormat = `{"output":"%s", "log": "%s", "is_error": "%v"}`

type ResponseWrapper struct {
	Log     string `json:"log"`
	Output  string `json:"output"`
	IsError bool   `json:"is_error"`
}

func HandleResponse(output, log string, exitWithError bool) {
	resp := ResponseWrapper{
		Log: log,
	}

	result, isError := formatErrorMessage(output)
	if result == "" {
		result = generalErrorMessage
	}

	resp.Output = result
	resp.IsError = isError || exitWithError

	marshaledResponse, err := json.Marshal(resp)
	if err != nil {
		updatedLog := fmt.Sprintf("%s\nfailed to marshal response: %v", log, err.Error())
		fmt.Printf(responseWrapperFormat, resp.Output, updatedLog, resp.IsError)
		return
	}

	fmt.Println(string(marshaledResponse))
}

// formatErrorMessage format the error message to be cleaner
func formatErrorMessage(result string) (msg string, isError bool) {
	if strings.Contains(result, queryWarningMessage) {
		result = strings.TrimSpace(strings.ReplaceAll(result, queryWarningMessage, ""))
		// older versions of the runner do not correctly report an error when steampipe returns a nonzero error,
		// so we have to parse the error message to determine if it is an error anyway
		isError = true
	}
	if strings.Contains(result, queryErrorMessage) {
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
			return "Connection have unsupported regions", isError
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
