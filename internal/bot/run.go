package bot

import (
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	start       = "/start"
	createIssue = "/create_issue"
	getIssues   = "/get_issues"

	successAuth = "Вы авторизованны! Теперь можете создавать задачи."
	errorAuth   = "Ошибка авторизации, попробуйте снова"

	authMessage        = "Введите логин:пароль JIRA"
	createIssueMessage = "Введите название проекта"
)

func (b *Bot) Run() {
	u := tgBotApi.NewUpdate(0)
	for update := range b.BotAPI.GetUpdatesChan(u) {
		if update.Message != nil {
			b.Logger.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			message := b.checkCommand(update)
			msg := tgBotApi.NewMessage(update.Message.Chat.ID, message)

			b.BotAPI.Send(msg)
		}
	}
}

func (b *Bot) checkCommand(update tgBotApi.Update) (message string) {
	switch update.Message.Text {
	case start:
		b.cacheCommands[update.Message.Chat.ID] = userCommand{
			pastCommand: start,
		}
		message = authMessage
	case createIssue:
		b.cacheCommands[update.Message.Chat.ID] = userCommand{
			pastCommand: createIssue,
		}
		message = createIssueMessage
	case getIssues:
		//b.Logger.Println(b.getAllIssues(update.Message.Chat.ID))

		message = b.getTasks()
	default:
		switch b.cacheCommands[update.Message.Chat.ID].pastCommand {
		case start:
			err := b.start(update)
			if err != nil {
				b.Logger.Println(err)
				message = errorAuth
				return
			}
			b.cacheCommands[update.Message.Chat.ID] = userCommand{}
			message = successAuth
		case createIssue:
			message = b.createIssue(update.Message.Chat.ID, update.Message.Text)
		}
	}
	return
}
