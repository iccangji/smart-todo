package todo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenRouterResponse struct {
	Message string `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func AskAI(prompt string) (string, error) {

	body := OpenRouterRequest{
		Model: os.Getenv("OPENROUTER_MODEL"),

		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"POST",
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+os.Getenv("OPENROUTER_API_KEY"),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	content := result["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	return content, nil
}
