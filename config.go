package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	NotionToken      string `json:"notion_token"`
	TelegramAppID    int    `json:"telegram_app_id"`
	TelegramAppHash  string `json:"telegram_app_hash"`
	TelegramPhone    string `json:"telegram_phone"`
	AIKey            string `json:"ai_key"`
	NotionPageID     string `json:"notion_page_id"`
	NotionDatabaseID string `json:"notion_database_id"`
	DBPath           string `json:"db_path"`
	SessionFile      string `json:"session_file"` // New field for session storage
}

func loadConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	if config.DBPath == "" {
		config.DBPath = "local.db"
	}
	if config.SessionFile == "" {
		config.SessionFile = "telegram_session.json" // Default session file
	}
	return &config, nil
}

func saveConfig(file string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 0644)
}
