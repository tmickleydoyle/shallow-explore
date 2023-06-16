package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ChatGPTRequest struct {
	Message string `json:"message"`
}

type ChatGPTResponse struct {
	Reply string `json:"reply"`
}

func CallChatGPTAPI(message string) (string, error) {
	payload := ChatGPTRequest{
		Message: message,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := http.DefaultClient

	resp, err := client.Post("https://api.chatgpt.com/v1/chat/completions", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiResponse ChatGPTResponse
	err = json.Unmarshal(respBody, &apiResponse)
	if err != nil {
		return "", err
	}

	return apiResponse.Reply, nil
}
