package view

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/mymmrac/telego"
	"log"
	"strconv"
)

const (
	UserId    = "user_id"
	FacultyId = "faculty_id"
	GroupId   = "group_id"
)

type CmdView interface {
	GetCmdFunc() CmdFunc
	GetCmdStrKey() string
}

type CmdFunc func(ctx context.Context, bot *telego.Bot, update telego.Update) error

type CmdContext struct {
	Cache *ristretto.Cache[string, string]
}

type Cmd struct {
	ViewContext *CmdContext
	ViewFunc    CmdFunc
	CmdString   string
}

func (c *Cmd) GetCmdStrKey() string {
	return c.CmdString
}

func (c *Cmd) GetCmdFunc() CmdFunc {
	return c.ViewFunc
}

func sendErrorMessage(ctx context.Context, bot *telego.Bot, update telego.Update, text string) error {
	_, err := bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID: update.Message.Chat.ChatID(),
		Text:   text,
	})
	if err != nil {
		return err
	}
	return nil
}

func getIdAndCache(c *Cmd, update telego.Update) (string, error) {
	userIdValue, isId := c.ViewContext.Cache.Get(UserId)
	if !isId {
		userIdValue := strconv.Itoa(int(update.Message.From.ID))

		errStr := fmt.Sprintf("%s command called but userId hasn't been found %s, saving it into a cache",
			c.CmdString, userIdValue)
		log.Println(errStr)

		if !c.ViewContext.Cache.Set(UserId, userIdValue, 0) {
			errStr := fmt.Sprintf("Error during userId save %s", userIdValue)
			log.Println(errStr)
			return "", errors.New(errStr)
		}
	}
	return userIdValue, nil
}
