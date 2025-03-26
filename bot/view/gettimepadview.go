package view

import (
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"time"
	"universityTimepad/model"
	"universityTimepad/parser"
	"universityTimepad/repo"
)

func GetTimepadView(cmd *Cmd, faculties *repo.Repository[model.Faculty], groups *repo.Repository[model.Group]) CmdFunc {
	const spbstuGetIcalAddress = "https://ruz.spbstu.ru/faculty/%d/groups/%d"

	const facultyNotSet = "Факультет не установлен. Укажите, пожалуйста, название факультета после команды /setFaculty."
	const groupNotSet = "Группа не установлена. Укажите, пожалуйста, название группы после команды /setGroup."
	const weekFormat = "Факультет: %s, группа: %s\n%s "
	const dayFormat = "%s, %s\n%s"
	const lessonFormat = "%s - %s\n%s\n"

	return func(ctx context.Context, bot *telego.Bot, update telego.Update) error {
		cmdText, _, _ := telegoutil.ParseCommand(update.Message.Text)

		if cmdText != cmd.CmdString {
			return nil
		}

		var cache = cmd.ViewContext.Cache
		_, err := getIdAndCache(cmd, update)
		if err != nil {
			return err
		}

		groupId, isFound := cache.Get(GroupId)
		if !isFound {
			return sendErrorMessage(ctx, bot, update, groupNotSet)
		}
		groupIdInt, err := strconv.Atoi(groupId)
		if err != nil {
			return sendErrorMessage(ctx, bot, update, groupNotSet)
		}
		group, err := (*groups).GetByID(ctx, groupIdInt)
		if err != nil {
			return sendErrorMessage(ctx, bot, update, groupNotSet)
		}
		if group == nil {
			return sendErrorMessage(ctx, bot, update, groupNotSet)
		}

		fId, isFound := cache.Get(FacultyId)
		if !isFound {
			return sendErrorMessage(ctx, bot, update, facultyNotSet)
		}
		fIdI, err := strconv.Atoi(fId)
		if err != nil {
			return sendErrorMessage(ctx, bot, update, facultyNotSet)
		}
		faculty, err := (*faculties).GetByID(ctx, fIdI)
		if err != nil {
			return sendErrorMessage(ctx, bot, update, facultyNotSet)
		}
		if faculty == nil {
			return sendErrorMessage(ctx, bot, update, facultyNotSet)
		}

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		var url = fmt.Sprintf(spbstuGetIcalAddress, faculty.GetId(), group.GetId())

		req.SetRequestURI(url)
		req.Header.SetMethod(fasthttp.MethodGet)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		var client = &fasthttp.Client{}
		err = client.Do(req, resp)
		if err != nil {
			return err
		}

		if statusCode := resp.StatusCode(); statusCode >= fasthttp.StatusInternalServerError {
			err := sendErrorMessage(ctx, bot, update, "internal server error")
			if err != nil {
				return err
			}
			return fmt.Errorf("internal server error: %d", statusCode)
		}

		timepads, err := parser.ReadTimepad(resp.Body(), group.GetId())
		if err != nil {
			return sendErrorMessage(ctx, bot, update, facultyNotSet)
		}

		days := make([]string, 6)
		for i, timepad := range timepads {
			lessons := make([]string, len(timepad.Lessons))
			for j, lesson := range timepad.Lessons {
				fields := make([]string, 6)
				fields[0] = lesson.Subject
				fields[1] = lesson.TypeObj.Name
				if len(lesson.Teachers) > 0 {
					fields[2] = lesson.Teachers[0].FullName
				}
				if len(lesson.Auditories) > 0 {
					fields[3] =
						fmt.Sprintf("%s, %s", lesson.Auditories[0].Name, lesson.Auditories[0].Building.Name)
				}
				lessonString := fmt.Sprintf(lessonFormat,
					lesson.TimeStart,
					lesson.TimeEnd,
					strings.Join(fields, "\n"),
				)
				lessons[j] = lessonString
			}
			days[i] = fmt.Sprintf(dayFormat, timepad.Date, time.Weekday(timepad.Weekday), strings.Join(lessons, "\n"))
		}

		_, err = bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text:   fmt.Sprintf(weekFormat, faculty.Name, group.Name, strings.Join(days, "\n")),
		})

		if err != nil {
			return err
		}
		return nil
	}
}
