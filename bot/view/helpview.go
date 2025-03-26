package view

import (
	"context"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"log"
	"strconv"
)

func Help(c *Cmd) CmdFunc {
	const text = "Список допустимых команд: \n" +
		"/setFaculty *название вашего факультета* \n" +
		"/setGroup *название вашей группы* \n" +
		"/timepad - получить расписание на ближайшую неделю для указанного факультета/группы"

	return func(ctx context.Context, bot *telego.Bot, update telego.Update) error {
		cmdText, _, _ := telegoutil.ParseCommand(update.Message.Text)

		if cmdText != c.CmdString {
			return nil
		}

		log.Printf("Help command executed with userId %s", strconv.Itoa(int(update.Message.From.ID)))

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
