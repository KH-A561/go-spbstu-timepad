package main

import (
	"context"
	"fmt"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	"log"
	"os"
	"universityTimepad/bot"
	"universityTimepad/bot/view"
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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	ctx := context.Background()
	token, _ := os.LookupEnv(tokenEnvName)

	botRef, err := telego.NewBot(token, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		return
	}

	botKit := bot.New(botRef)
	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 10000,
		MaxCost:     100000000,
		BufferItems: 64,
	})

	if err != nil {
		panic(err)
	}

	defer cache.Close()
	cmdContext := &view.CmdContext{Cache: cache}

	startCmd := &view.Cmd{CmdString: "start", ViewContext: cmdContext}
	startCmd.ViewFunc = view.Start(startCmd)
	botKit.RegisterCmdView(startCmd)

	helpCmd := &view.Cmd{CmdString: "help", ViewContext: cmdContext}
	helpCmd.ViewFunc = view.Help(helpCmd)
	botKit.RegisterCmdView(helpCmd)

	setFacultyCmd := &view.Cmd{CmdString: "setFaculty", ViewContext: cmdContext}
	setFacultyCmd.ViewFunc = view.SetFaculty(setFacultyCmd, &FacRepo)
	botKit.RegisterCmdView(setFacultyCmd)

	setGroupCmd := &view.Cmd{CmdString: "setGroup", ViewContext: cmdContext}
	setGroupCmd.ViewFunc = view.SetGroup(setGroupCmd, &GroupRepo)
	botKit.RegisterCmdView(setGroupCmd)

	getTimepadCmd := &view.Cmd{CmdString: "timepad", ViewContext: cmdContext}
	getTimepadCmd.ViewFunc = view.GetTimepadView(getTimepadCmd, &FacRepo, &GroupRepo)
	botKit.RegisterCmdView(getTimepadCmd)

	if err := botKit.Run(ctx); err != nil {
		log.Printf("[ERROR] failed to run botkit: %v", err)
	}
}
