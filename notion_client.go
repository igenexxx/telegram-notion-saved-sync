package main

import (
	"context"
	"github.com/jomei/notionapi"
	"time"
)

type NotionClient interface {
	CreateDatabase(parentPageID string) (string, error)
	InsertRow(databaseID, title, category string, date time.Time, id int, description, link, content string) error
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
	if parentPageID == "" {
		page, err := n.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
			Parent: notionapi.Parent{
				Type: notionapi.ParentTypeWorkspace,
			},
			Properties: notionapi.Properties{
				"title": notionapi.TitleProperty{
					Title: []notionapi.RichText{{
						Text: &notionapi.Text{
							Content: "Telegram Messages",
						},
					}},
				},
			},
		})

		if err != nil {
			return "", err
		}

		parentPageID = string(page.ID)
	}

	db, err := n.client.Database.Create(context.Background(), &notionapi.DatabaseCreateRequest{
		Parent: notionapi.Parent{
			Type:   notionapi.ParentTypePageID,
			PageID: notionapi.PageID(parentPageID),
		},
		Title: []notionapi.RichText{{Text: &notionapi.Text{Content: "Messages"}}},
		Properties: notionapi.PropertyConfigs{
			"Title":             notionapi.TitlePropertyConfig{Type: notionapi.PropertyConfigType(notionapi.PropertyTypeTitle)},
			"Category":          notionapi.SelectPropertyConfig{Type: notionapi.PropertyConfigType(notionapi.PropertyTypeSelect)},
			"Date of posting":   notionapi.DatePropertyConfig{Type: notionapi.PropertyConfigType(notionapi.PropertyTypeDate)},
			"ID":                notionapi.NumberPropertyConfig{Type: notionapi.PropertyConfigType(notionapi.PropertyTypeNumber)},
			"Short description": notionapi.RichTextPropertyConfig{Type: notionapi.PropertyConfigType(notionapi.PropertyTypeRichText)},
			"Link":              notionapi.URLPropertyConfig{Type: notionapi.PropertyConfigType(notionapi.PropertyTypeURL)},
			"Content":           notionapi.RichTextPropertyConfig{Type: notionapi.PropertyConfigType(notionapi.PropertyTypeRichText)},
		},
	})

	if err != nil {
		return "", err
	}

	return string(db.ID), nil
}

func (n *notionClient) InsertRow(databaseID, title, category string, date time.Time, id int, description, link,
	content string) error {
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
			"Date of posting": notionapi.DateProperty{
				Date: func() *notionapi.DateObject {
					d := notionapi.Date(date)
					return &notionapi.DateObject{Start: &d}
				}(),
			},
			"ID": notionapi.NumberProperty{
				Number: float64(id),
			},
			"Short description": notionapi.RichTextProperty{
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
