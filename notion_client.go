package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jomei/notionapi"
	"io"
	"net/http"
	"time"
)

type NotionClient interface {
	CreateDatabase(parentPageID string) (string, error)
	InsertRow(databaseID string, title, category string, date time.Time, id int, description, link, content string) error
}

type notionClient struct {
	client *notionapi.Client
}

func NewNotionClient(token string) NotionClient {
	return &notionClient{
		client: notionapi.NewClient(notionapi.Token(token)),
	}
}

func (n *notionClient) CreateDatabase(parentPageID string) (string, error) {
	// If no parentPageID is provided, create a new page first
	if parentPageID == "" {
		page, err := n.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
			Parent: notionapi.Parent{
				Type: notionapi.ParentTypeWorkspace,
			},
			Properties: notionapi.Properties{
				"title": notionapi.TitleProperty{
					Title: []notionapi.RichText{{Text: &notionapi.Text{Content: "Telegram Messages"}}},
				},
			},
		})
		if err != nil {
			return "", fmt.Errorf("create parent page: %v", err)
		}
		parentPageID = string(page.ID)
	}

	// Construct the raw JSON payload for the database creation
	payload := map[string]interface{}{
		"parent": map[string]string{
			"type":    "page_id",
			"page_id": parentPageID,
		},
		"title": []map[string]interface{}{
			{
				"type": "text",
				"text": map[string]string{
					"content": "Messages",
				},
			},
		},
		"properties": map[string]interface{}{
			"Title": map[string]interface{}{
				"type":  "title",
				"title": map[string]interface{}{}, // Empty title config
			},
			"Category": map[string]interface{}{
				"type":   "select",
				"select": map[string]interface{}{}, // Empty select config (options added dynamically)
			},
			"Date of Posting": map[string]interface{}{
				"type": "date",
				"date": map[string]interface{}{}, // Empty date config
			},
			"ID": map[string]interface{}{
				"type":   "number",
				"number": map[string]interface{}{}, // Empty number config
			},
			"Short Description": map[string]interface{}{
				"type":      "rich_text",
				"rich_text": map[string]interface{}{}, // Empty rich_text config
			},
			"Link": map[string]interface{}{
				"type": "url",
				"url":  map[string]interface{}{}, // Empty url config
			},
			"Content": map[string]interface{}{
				"type":      "rich_text",
				"rich_text": map[string]interface{}{}, // Empty rich_text config
			},
		},
	}

	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %v", err)
	}

	// Create and send the HTTP request
	req, err := http.NewRequest("POST", "https://api.notion.com/v1/databases", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+string(n.client.Token))
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	// Use a default HTTP client
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("notion API error: %s", string(body))
	}

	// Parse the response to get the database ID
	var db notionapi.Database
	if err := json.NewDecoder(resp.Body).Decode(&db); err != nil {
		return "", fmt.Errorf("decode response: %v", err)
	}

	return string(db.ID), nil
}

func (n *notionClient) InsertRow(databaseID, title, category string, date time.Time, id int, description, link, content string) error {
	_, err := n.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(databaseID),
		},
		Properties: notionapi.Properties{
			"Title": notionapi.TitleProperty{
				Title: []notionapi.RichText{{Text: &notionapi.Text{Content: title}}},
			},
			"Category": notionapi.SelectProperty{
				Select: notionapi.Option{Name: category},
			},
			"Date of Posting": notionapi.DateProperty{
				Date: &notionapi.DateObject{Start: (*notionapi.Date)(&date)},
			},
			"ID": notionapi.NumberProperty{
				Number: float64(id),
			},
			"Short Description": notionapi.RichTextProperty{
				RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: description}}},
			},
			"Link": notionapi.URLProperty{
				URL: link,
			},
			"Content": notionapi.RichTextProperty{
				RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: content}}},
			},
		},
	})
	return err
}
