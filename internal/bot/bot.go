package bot

import (
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"jira_bot/internal/services"
	"log"
	"os"
)

type userCommand struct {
	pastCommand      string
	projectNameIssue string
	titleIssue       string
	textIssue        string
	priorityIssue    string
	executorIssue    string
}

type Bot struct {
	BotAPI        *tgBotApi.BotAPI
	Logger        *log.Logger
	JiraURL       string
	services      *services.Services
	cacheCommands map[int64]userCommand
}

func NewBot(token, jiraURL string, s *services.Services) (*Bot, error) {
	api, err := tgBotApi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		BotAPI:        api,
		JiraURL:       jiraURL,
		cacheCommands: make(map[int64]userCommand, 0),
		services:      s,
	}

	if bot.Logger == nil {
		bot.Logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	_, err = bot.initCommands()
	if err != nil {
		return nil, err
	}

	return bot, nil
}
