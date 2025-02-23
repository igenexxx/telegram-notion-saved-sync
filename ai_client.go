package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"strings"
)

type AIClient interface {
	AnalyzeMessage(content string) (category, description string, err error)
}

type openAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) AIClient {
	return &openAIClient{
		client: openai.NewClient(apiKey),
	}
}

func (a *openAIClient) AnalyzeMessage(content string) (string, string, error) {
	prompt := fmt.Sprintf(`
Analyze the following message and provide the category and a short description of the message in a JSON format.

Message: %s

Response format:
{
	"category": "category name",
	"short_description": "short description"
}
`, content)

	resp, err := a.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})

	if err != nil {
		return "", "", err
	}

	var result struct {
		Category         string `json:"category"`
		ShortDescription string `json:"short_description"`
	}
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result)
	if err != nil {
		return "", "", err
	}

	return result.Category, strings.TrimSpace(result.ShortDescription), nil
}
