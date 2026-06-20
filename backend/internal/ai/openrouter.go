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

// Get response with stream
func StreamResponse(
	feature AiFeatures,
	data []byte,
	writer io.Writer,
	flusher http.Flusher,
	onData func(string),
) error {
	prompt := GeneratePrompt(feature, string(data))
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
		fmt.Printf("[RAW] %s\n", line)

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

		if feature == Breakdown {
			buffer.WriteString(content)
			text := buffer.String()

			if strings.Contains(text, "\n") {

				fmt.Fprintf(writer, "data: %s\n\n", strings.TrimSpace(text))
				flusher.Flush()

				if onData != nil {
					onData(strings.TrimSpace(text))
				}

				buffer.Reset()
			}
		} else {
			fmt.Fprintf(writer, "data: %s\n\n", content)
			flusher.Flush()

			if onData != nil {
				onData(content)
			}

			continue
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

// Get response as single string
func FetchResponse(
	feature AiFeatures,
	data []byte,
	onData func(string),
) error {
	prompt := GeneratePrompt(feature, string(data))
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
		fmt.Printf("[RAW] %s\n", line)

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

		if feature == Breakdown {
			buffer.WriteString(content)
			text := buffer.String()

			if strings.Contains(text, "\n") {
				if onData != nil {
					onData(strings.TrimSpace(text))
				}

				buffer.Reset()
			}
		} else {
			if onData != nil {
				onData(content)
			}
			continue
		}
	}

	remaining := strings.TrimSpace(
		buffer.String(),
	)

	if remaining != "" {
		if onData != nil {
			onData(strings.TrimSpace(remaining))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
