package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	"log"
	"time"
	"universityTimepad/calendar"
	"universityTimepad/model"
	"universityTimepad/repo"
	"universityTimepad/repo/inmem"
)

const (
	tokenEnvName = "TELEGRAM_API_TOKEN"
)

var (
	InmemConfig = inmem.NewDefault()
	Config      = repo.Config{
		StorageType: repo.StorageTypeInmemory,
		Inmemory:    InmemConfig,
	}
	FacRepo, _   = repo.New[model.Faculty](&Config, InmemConfig.FacInitFunc.GetInitFunc())
	GroupRepo, _ = repo.New[model.Group](&Config, InmemConfig.GroupInitFunc.GetInitFunc())
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {

	//ctx := context.Background()
	//token, _ := os.LookupEnv(tokenEnvName)
	//
	//bot, err := telego.NewBot(token, telego.WithDefaultDebugLogger())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//_, err = bot.GetMe(ctx)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}

	//updateLoop(bot, &ctx)
}

func updateLoop(bot *telego.Bot, ctx *context.Context) {
	//updChan, err := bot.UpdatesViaLongPolling(*ctx, nil)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//for upd := range updChan {
	//	fmt.Println(upd)
	//
	//	if upd.Message != nil {
	//		_, err := bot.CopyMessage(*ctx, &telego.CopyMessageParams{
	//			ChatID:     telego.ChatID{ID: upd.Message.Chat.ID},
	//			FromChatID: telego.ChatID{ID: upd.Message.Chat.ID},
	//			MessageID:  upd.Message.MessageID,
	//		})
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//	}
	//}

	var fileSnatcher = calendar.DefaultFileSnatcher
	var date, _ = time.Parse(calendar.DateFormat, "2025-3-25")
	fileSnatcher.SnatchCalendar(
		*ctx,
		124,
		40291,
		date,
	)
}
