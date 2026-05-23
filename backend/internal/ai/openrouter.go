package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type StreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

func StreamResponse(
	prompt string,
	writer io.Writer,
	flusher http.Flusher,
	onData func(string),
) error {

	body := OpenRouterRequest{
		Model:  os.Getenv("OPENROUTER_MODEL"),
		Stream: true,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(
		"POST",
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return err
	}

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
		return err
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	var buffer strings.Builder

	for scanner.Scan() {

		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(
			line,
			"data: ",
		)

		if data == "[DONE]" {
			break
		}

		var chunk StreamChunk

		err := json.Unmarshal(
			[]byte(data),
			&chunk,
		)

		if err != nil {
			continue
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		content := chunk.
			Choices[0].
			Delta.
			Content

		if content == "" {
			continue
		}

		buffer.WriteString(content)

		text := buffer.String()

		if strings.Contains(text, "\n") {

			fmt.Fprintf(
				writer,
				"data: %s\n\n",
				strings.TrimSpace(text),
			)

			flusher.Flush()

			if onData != nil {
				onData(strings.TrimSpace(text))
			}

			buffer.Reset()
		}
	}

	remaining := strings.TrimSpace(
		buffer.String(),
	)

	if remaining != "" {
		fmt.Fprintf(
			writer,
			"data: %s\n\n",
			remaining,
		)
		if onData != nil {
			onData(strings.TrimSpace(remaining))
		}
		flusher.Flush()
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
