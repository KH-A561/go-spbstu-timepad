package bot

import (
	"context"
	"github.com/mymmrac/telego"
)

type ViewFunc func(ctx context.Context, bot *telego.Bot, update telego.Update) error

type Bot struct {
	botRef   *telego.Bot
	cmdViews map[string]ViewFunc
}

func New(botRef *telego.Bot) *Bot {
	return &Bot{botRef: botRef}
}

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}

	b.cmdViews[cmd] = view
}

func (b *Bot) Run(ctx context.Context) error {
	
}
