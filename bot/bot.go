package bot

import (
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"log"
	"runtime/debug"
	"slices"
	"time"
	"universityTimepad/bot/view"
)

type Bot struct {
	botRef   *telego.Bot
	cmdViews map[string]view.CmdFunc
}

func New(botRef *telego.Bot) *Bot {
	return &Bot{botRef: botRef}
}

func (b *Bot) RegisterCmdView(cmd view.CmdView) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]view.CmdFunc)
	}

	b.cmdViews[cmd.GetCmdStrKey()] = cmd.GetCmdFunc()
}

func (b *Bot) Run(ctx context.Context) error {
	updChan, err := b.botRef.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for {
		select {
		case update := <-updChan:
			fmt.Println(update)
			updateCtx, updateCancel := context.WithTimeout(context.Background(), 5*time.Minute)
			b.handleUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update telego.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	if update.Message == nil {
		return
	}

	cmdIndex := slices.IndexFunc(update.Message.Entities,
		func(e telego.MessageEntity) bool { return e.Type == telego.EntityTypeBotCommand })

	if cmdIndex == -1 && update.CallbackQuery == nil {
		return
	}

	cmd, _, _ := telegoutil.ParseCommand(update.Message.Text)

	cmdView, ok := b.cmdViews[cmd]
	if !ok {
		log.Printf("[ERROR] command %s not found", cmd)
		return
	}

	if err := cmdView(ctx, b.botRef, update); err != nil {
		log.Printf("[ERROR] failed to execute view: %v", err)

		if _, err := b.botRef.SendMessage(
			ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: update.Message.Chat.ID},
				Text:   fmt.Sprintf("%s error", err)},
		); err != nil {
			log.Printf("[ERROR] failed to send error message: %v", err)
		}
	}
}
