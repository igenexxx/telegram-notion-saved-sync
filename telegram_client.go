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
	Run(ctx context.Context, fetchFunc func(ctx context.Context, fromID int) ([]TelegramMessage, int, error)) error
}

type telegramClient struct {
	client      *telegram.Client
	phone       string
	sessionFile string
}

func NewTelegramClient(appID int, appHash, phone, sessionFile string) (TelegramClient, error) {
	sessionStorage := &session.FileStorage{Path: sessionFile}
	client := telegram.NewClient(appID, appHash, telegram.Options{
		SessionStorage: sessionStorage,
	})
	return &telegramClient{client: client, phone: phone, sessionFile: sessionFile}, nil
}

func (t *telegramClient) Run(ctx context.Context, fetchFunc func(ctx context.Context, fromID int) ([]TelegramMessage, int, error)) error {
	return t.client.Run(ctx, func(ctx context.Context) error {
		// Authentication flow
		flow := auth.NewFlow(auth.Constant(t.phone, "", auth.CodeAuthenticatorFunc(func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
			fmt.Print("Enter Telegram auth code: ")
			var code string
			fmt.Scanln(&code)
			return code, nil
		})), auth.SendCodeOptions{})
		if err := flow.Run(ctx, t.client.Auth()); err != nil {
			return fmt.Errorf("auth: %v", err)
		}

		// Get self to verify connection
		_, err := t.client.Self(ctx)
		if err != nil {
			return fmt.Errorf("get self: %v", err)
		}

		// Keep the client running and delegate message fetching to the provided function
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Call the fetch function (passed from main.go)
				messages, latestID, err := fetchFunc(ctx, 0) // fromID will be managed in main.go
				if err != nil {
					return fmt.Errorf("fetch messages: %v", err)
				}
				// Process messages here or let main.go handle them
				fmt.Printf("Fetched %d messages, latest ID: %d\n", len(messages), latestID)
				time.Sleep(30 * time.Second) // Polling interval
			}
		}
	})
}

func (t *telegramClient) FetchMessages(ctx context.Context, fromID int) ([]TelegramMessage, int, error) {
	var messages []TelegramMessage
	var latestID int

	req := &tg.MessagesGetHistoryRequest{
		Peer:     &tg.InputPeerSelf{},
		OffsetID: fromID,
		Limit:    100, // Adjust batch size as needed
	}
	history, err := t.client.API().MessagesGetHistory(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("get history: %v", err)
	}

	switch h := history.(type) {
	case *tg.MessagesMessagesSlice:
		for _, msg := range h.Messages {
			m, ok := msg.(*tg.Message)
			if !ok || m.Message == "" {
				continue
			}
			latestID = m.ID
			messages = append(messages, TelegramMessage{
				ID:      m.ID,
				Date:    time.Unix(int64(m.Date), 0),
				Content: m.Message,
				Link:    fmt.Sprintf("https://t.me/me/%d", m.ID),
			})
		}
	case *tg.MessagesMessages:
		for _, msg := range h.Messages {
			m, ok := msg.(*tg.Message)
			if !ok || m.Message == "" {
				continue
			}
			latestID = m.ID
			messages = append(messages, TelegramMessage{
				ID:      m.ID,
				Date:    time.Unix(int64(m.Date), 0),
				Content: m.Message,
				Link:    fmt.Sprintf("https://t.me/me/%d", m.ID),
			})
		}
	case *tg.MessagesMessagesNotModified:
		return messages, latestID, nil // No new messages
	case *tg.MessagesChannelMessages:
		for _, msg := range h.Messages {
			m, ok := msg.(*tg.Message)
			if !ok || m.Message == "" {
				continue
			}
			latestID = m.ID
			messages = append(messages, TelegramMessage{
				ID:      m.ID,
				Date:    time.Unix(int64(m.Date), 0),
				Content: m.Message,
				Link:    fmt.Sprintf("https://t.me/me/%d", m.ID),
			})
		}
	default:
		return nil, 0, fmt.Errorf("unexpected history type: %T", history)
	}

	return messages, latestID, nil
}
