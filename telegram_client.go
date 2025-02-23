package main

import (
	"context"
	"fmt"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"time"
)

type TelegramMessage struct {
	ID      int
	Date    time.Time
	Content string
	Link    string
}

type TelegramClient interface {
	FetchMessages(ctx context.Context, fromID int) ([]TelegramMessage, int, error)
}

type telegramClient struct {
	client      *telegram.Client
	phone       string
	sessionFile string
}

func NewTelegramClient(appID int, appHash, phone, sessionFile string) (TelegramClient, error) {
	// Load or create session storage
	sessionStorage := &session.FileStorage{Path: sessionFile}

	client := telegram.NewClient(appID, appHash, telegram.Options{
		SessionStorage: sessionStorage,
	})

	return &telegramClient{
		client:      client,
		phone:       phone,
		sessionFile: sessionFile,
	}, nil
}

func (t *telegramClient) FetchMessages(ctx context.Context, fromID int) ([]TelegramMessage, int, error) {
	var messages []TelegramMessage
	var latestID int

	err := t.client.Run(ctx, func(ctx context.Context) error {
		// Authentication flow
		flow := auth.NewFlow(auth.Constant(t.phone, "", auth.CodeAuthenticatorFunc(func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
			fmt.Print("Enter Telegram auth code: ")
			var code string
			fmt.Scanln(&code)
			return code, nil
		})), auth.SendCodeOptions{})

		// Attempt to authenticate; if session exists, this will succeed without code
		if err := flow.Run(ctx, t.client.Auth()); err != nil {
			return fmt.Errorf("auth: %v", err)
		}

		// Get self (Saved Messages is your own chat)
		_, err := t.client.Self(ctx)
		if err != nil {
			return fmt.Errorf("get self: %v", err)
		}

		// Fetch messages from Saved Messages
		req := &tg.MessagesGetHistoryRequest{
			Peer:     &tg.InputPeerSelf{},
			OffsetID: fromID,
			Limit:    100,
		}
		history, err := t.client.API().MessagesGetHistory(ctx, req)
		if err != nil {
			return fmt.Errorf("get history: %v", err)
		}

		msgList, ok := history.(*tg.MessagesMessagesSlice)
		if !ok {
			return fmt.Errorf("unexpected history type")
		}

		for _, msg := range msgList.Messages {
			m, ok := msg.(*tg.Message)
			if !ok || m.Message == "" {
				continue
			}
			latestID = m.ID // Track the latest ID
			messages = append(messages, TelegramMessage{
				ID:      m.ID,
				Date:    time.Unix(int64(m.Date), 0),
				Content: m.Message,
				Link:    fmt.Sprintf("https://t.me/me/%d", m.ID), // Approximate Saved Messages link
			})
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return messages, latestID, nil
}
