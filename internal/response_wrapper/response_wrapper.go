package response_wrapper

import (
	"encoding/json"
	"fmt"
)

const marshalError = "failed to marshal response"

type ResponseWrapper struct {
	Log    string `json:"log"`
	Output string `json:"output"`
}

func HandleResponse(output, log string) {
	resp := ResponseWrapper{
		Log:    log,
		Output: output,
	}

	marshaledResponse, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf(`{"output":"%s", "log": "%s"}`, output, marshalError)
		return
	}

	fmt.Println(string(marshaledResponse))
}
