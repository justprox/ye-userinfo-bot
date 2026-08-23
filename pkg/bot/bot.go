package bot

import (
	"context"

	tgbot "github.com/go-telegram/bot"
)

// New creates and configures a new Telegram bot instance with zero-logging and safe defaults.
func New(token string) (*tgbot.Bot, error) {
	opts := []tgbot.Option{
		tgbot.WithDefaultHandler(HandleUpdate),
		tgbot.WithErrorsHandler(func(err error) {
			// Silent error handler adhering strictly to zero-logging requirements
		}),
	}
	return tgbot.New(token, opts...)
}

// SetupCommands registers multi-language commands in Telegram menu.
func SetupCommands(ctx context.Context, b *tgbot.Bot) {
	setupBotCommands(ctx, b)
}
