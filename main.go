package main

import (
	"context"
	"log"
	"time"
)

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal("Load config:", err)
	}

	db, err := NewSQLiteDatabase(config.DBPath)
	if err != nil {
		log.Fatal("Init DB:", err)
	}

	aiClient := NewOpenAIClient(config.AIKey)
	notionClient := NewNotionClient(config.NotionToken)
	tgClient, err := NewTelegramClient(config.TelegramAppID, config.TelegramAppHash, config.TelegramPhone, config.SessionFile)
	if err != nil {
		log.Fatal("Init Telegram:", err)
	}

	// Setup Notion database if not exists
	if config.NotionDatabaseID == "" {
		databaseID, err := notionClient.CreateDatabase("")
		if err != nil {
			log.Fatal("Create Notion DB:", err)
		}
		config.NotionDatabaseID = databaseID
		if err := saveConfig("config.json", config); err != nil {
			log.Fatal("Save config:", err)
		}
	}

	// Main loop for recurrent scanning
	lastID, err := db.GetLastUpdateID()
	if err != nil {
		log.Fatal("Get last ID:", err)
	}

	ctx := context.Background()
	for {
		messages, newLastID, err := tgClient.FetchMessages(ctx, lastID)
		if err != nil {
			log.Printf("Fetch messages: %v", err)
			time.Sleep(5 * time.Second) // Backoff on error
			continue
		}

		for _, msg := range messages {
			if err := processMessage(notionClient, aiClient, config.NotionDatabaseID, msg); err != nil {
				log.Printf("Process message %d: %v", msg.ID, err)
			}
		}

		if newLastID > lastID {
			lastID = newLastID
			if err := db.SetLastUpdateID(lastID); err != nil {
				log.Printf("Set last ID %d: %v", lastID, err)
			}
		}

		log.Printf("Processed %d messages, last ID: %d", len(messages), lastID)
		time.Sleep(30 * time.Second) // Adjust polling interval
	}
}

func processMessage(notionClient NotionClient, aiClient AIClient, databaseID string, msg TelegramMessage) error {
	content := msg.Content
	if content == "" {
		return nil // Skip empty messages
	}

	// Analyze with AI
	category, description, err := aiClient.AnalyzeMessage(content)
	if err != nil {
		return err
	}

	// Generate title (first 50 chars)
	title := content
	if len(title) > 50 {
		title = title[:50] + "..."
	}

	// Insert into Notion
	err = notionClient.InsertRow(databaseID, title, category, msg.Date, msg.ID, description, msg.Link, content)
	if err != nil {
		return err
	}
	return nil
}
