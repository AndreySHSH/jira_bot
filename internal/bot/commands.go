package bot

import (
	"errors"
	"github.com/andygrunwald/go-jira"
	"strings"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	errorAuthJira = "неправильный логин или пароль"
)

type commandEntity struct {
	key  string
	desc string
}

func (b *Bot) initCommands() (*tgBotApi.APIResponse, error) {
	commands := []commandEntity{
		{
			key:  "start",
			desc: "Запустить бота",
		},
		{
			key:  "create_issue",
			desc: "Создать задачу",
		},
		{
			key:  "get_issues",
			desc: "Список задач",
		},
	}

	tgCommands := make([]tgBotApi.BotCommand, 0, len(commands))

	for _, cmd := range commands {
		tgCommands = append(tgCommands, tgBotApi.BotCommand{
			Command:     cmd.key,
			Description: cmd.desc,
		})
	}

	return b.BotAPI.Request(tgBotApi.NewSetMyCommands(tgCommands...))
}

func (b *Bot) start(update tgBotApi.Update) error {

	s := strings.Split(update.Message.Text, ":")

	if len(s) != 2 {
		return errors.New(errorAuthJira)
	}

	client, err := b.authJira(s[0], s[1])
	if err != nil {
		return err
	}

	if client.Authentication.Authenticated() {
		return errors.New(errorAuthJira)
	}

	b.services.User.CreateUser(update.Message.Chat.ID, s[0], s[1])
	return nil
}

func (b *Bot) authJira(login, password string) (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: login,
		Password: password,
	}

	client, err := jira.NewClient(tp.Client(), b.JiraURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (b *Bot) createIssue(id int64, text string) string {

	user := b.services.User.GetUserByID(id)
	if user == nil {
		return errorAuthJira
	}

	u := b.cacheCommands[id]

	if u.projectNameIssue == "" {
		b.cacheCommands[id] = userCommand{
			pastCommand:      u.pastCommand,
			projectNameIssue: text,
		}
		return "Введите тему задачи"
	}

	if u.titleIssue == "" {
		b.cacheCommands[id] = userCommand{
			pastCommand:      u.pastCommand,
			projectNameIssue: u.projectNameIssue,
			titleIssue:       text,
		}
		return "Введите что нужно сделать"
	}

	if u.textIssue == "" {
		b.cacheCommands[id] = userCommand{
			pastCommand:      u.pastCommand,
			projectNameIssue: u.projectNameIssue,
			titleIssue:       u.titleIssue,
			textIssue:        text,
		}
		return "Выберите приоритет (обычный, срочно, очень срочно)"
	}

	if u.priorityIssue == "" {
		b.cacheCommands[id] = userCommand{
			pastCommand:      u.pastCommand,
			projectNameIssue: u.projectNameIssue,
			titleIssue:       u.titleIssue,
			textIssue:        u.textIssue,
			priorityIssue:    text,
		}
		return "Введите имя исполнителя"
	}

	if u.executorIssue == "" {
		b.cacheCommands[id] = userCommand{
			pastCommand:      u.pastCommand,
			projectNameIssue: u.projectNameIssue,
			titleIssue:       u.titleIssue,
			textIssue:        u.textIssue,
			priorityIssue:    u.priorityIssue,
			executorIssue:    text,
		}
	}

	u = b.cacheCommands[id]

	jiraClient, err := b.authJira(user.Login, user.Password)
	if err != nil {
		return errorAuthJira
	}

	if jiraClient.Authentication.Authenticated() {
		return errorAuthJira
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Project: jira.Project{
				Key: u.projectNameIssue,
			},
			Priority: &jira.Priority{
				Name: u.priorityIssue,
			},
			Assignee: &jira.User{
				Name: u.executorIssue,
			},
			Reporter: &jira.User{
				Name: user.Login,
			},
			Summary: u.textIssue,
		},
	}
	_, _, err = jiraClient.Issue.Create(&i)
	if err != nil {
		return "Ошибка создания задачи"
	}

	b.cacheCommands[id] = userCommand{}

	return "Задача создана"
}

func (b *Bot) getTasks() string {
	return "список задач"
}

func (b *Bot) getAllIssues(id int64) ([]jira.Issue, error) {

	user := b.services.User.GetUserByID(id)
	if user == nil {
		return nil, errors.New(errorAuthJira)
	}

	client, err := b.authJira(user.Login, user.Password)
	if err != nil {
		return nil, err
	}

	if !client.Authentication.Authenticated() {
		return nil, errors.New(errorAuthJira)
	}

	req, _ := client.NewRequest("GET", "rest/api/2/issue", nil)

	issues := new([]jira.Issue)
	_, err = client.Do(req, issues)
	if err != nil {
		panic(err)
	}

	return *issues, nil
}
