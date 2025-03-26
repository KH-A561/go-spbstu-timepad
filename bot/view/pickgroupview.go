package view

import (
	"context"
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"log"
	"strconv"
	"strings"
	"universityTimepad/model"
	"universityTimepad/repo"
)

func SetGroup(c *Cmd, groupRepo *repo.Repository[model.Group]) CmdFunc {
	const successText = "Группа установлена: %s."
	const emptyGroupNameArg = "Укажите, пожалуйста, название группы после команды /setGroup"
	const groupNotFound = "Группа с именем %s не найдена. " +
		"Укажите, пожалуйста, полное название группы."

	return func(ctx context.Context, bot *telego.Bot, update telego.Update) error {
		cmdText, _, args := telegoutil.ParseCommand(update.Message.Text)

		if cmdText != c.CmdString {
			return nil
		}

		userIdValue, err := getIdAndCache(c, update)
		if err != nil {
			return err
		}

		var groupName string
		if len(args) < 1 || args[0] == "" {
			return sendErrorMessage(ctx, bot, update, emptyGroupNameArg)
		}
		groupName = strings.Join(args[0:], " ")
		log.Printf("SetGroup command executed with userId %s, groupName written: %s", userIdValue, groupName)

		group, err := (*groupRepo).GetByName(ctx, groupName)
		if err != nil {
			return err
		}
		if group == nil {
			return sendErrorMessage(ctx, bot, update, fmt.Sprintf(groupNotFound, groupName))
		}

		if !c.ViewContext.Cache.Set(GroupId, strconv.Itoa(group.Id), 0) {
			errStr := fmt.Sprintf("Error during GroupId save %s", strconv.Itoa(group.Id))
			log.Println(errStr)
			return errors.New(errStr)
		}

		_, err = bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   fmt.Sprintf(successText, group.Name),
		})

		if err != nil {
			return err
		}
		return nil
	}
}
