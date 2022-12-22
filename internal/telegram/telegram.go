package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func Bot(messages, response chan string) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGTOKEN"))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "Отправь мне название города и код. Например, Moscow RU"
		case "status":
			msg.Text = "Я в порядке :3"
		default:
			messages <- update.Message.Text

			msg.Text = <-response
			//if len(update.Message.Text) != 2 {
			//	msg.Text = "Некорректные данные"
			//} else {
			//	messages <- update.Message.Text
			//	time.Sleep(time.Second)
			//	msg.Text = <-messages
			//}
		}

		if _, err = bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
