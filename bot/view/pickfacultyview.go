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

func SetFaculty(c *Cmd, facRepo *repo.Repository[model.Faculty]) CmdFunc {
	const successText = "Факультет установлен: %s (%s)."
	const emptyFacultyNameArg = "Укажите, пожалуйста, название факультета после команды /setFaculty"
	const facultyNotFound = "Факультет с именем %s не найден. " +
		"Укажите, пожалуйста, полное название или аббревиатуру факультета."

	return func(ctx context.Context, bot *telego.Bot, update telego.Update) error {
		cmdText, _, args := telegoutil.ParseCommand(update.Message.Text)

		if cmdText != c.CmdString {
			return nil
		}

		userIdValue, err := getIdAndCache(c, update)
		if err != nil {
			return err
		}

		var facultyName string
		if len(args) < 1 || args[0] == "" {
			return sendErrorMessage(ctx, bot, update, emptyFacultyNameArg)
		}
		facultyName = strings.Join(args[0:], " ")
		log.Printf("SetFaculty command executed with userId %s, facultyName written: %s", userIdValue, facultyName)

		faculty, err := (*facRepo).GetByName(ctx, facultyName)
		if err != nil {
			return err
		}
		if faculty == nil {
			return sendErrorMessage(ctx, bot, update, fmt.Sprintf(facultyNotFound, facultyName))
		}

		if !c.ViewContext.Cache.Set(FacultyId, strconv.Itoa(faculty.Id), 0) {
			errStr := fmt.Sprintf("Error during FacultyId save %s", strconv.Itoa(faculty.Id))
			log.Println(errStr)
			return errors.New(errStr)
		}

		_, err = bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   fmt.Sprintf(successText, faculty.Name, faculty.Abbr),
		})

		if err != nil {
			return err
		}
		return nil
	}
}
