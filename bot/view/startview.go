package view

import (
	"context"
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"log"
	"strconv"
)

func Start(c *Cmd) CmdFunc {
	const text = "Здравствуйте! Данный бот позволяет посмотреть расписание группы СПБПУ. \n" +
		"Воспользуйтесь командой /help, чтобы узнать список возможных команд."

	return func(ctx context.Context, bot *telego.Bot, update telego.Update) error {
		cmdText, _, _ := telegoutil.ParseCommand(update.Message.Text)

		if cmdText != c.CmdString {
			return nil
		}

		userIdValue := strconv.Itoa(int(update.Message.From.ID))
		if !c.ViewContext.Cache.Set(UserId, userIdValue, 0) {
			errStr := fmt.Sprintf("Start command executed but userId wasn't saved %s",
				userIdValue)
			log.Println(errStr)
			return errors.New(errStr)
		}

		log.Printf("Start command executed with userId %s", userIdValue)

		_, err := bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   text,
		})

		if err != nil {
			return err
		}
		return nil
	}
}
