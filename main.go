package main

import (
	// "errors"
	"bytes"
	"encoding/json"
	"example/user/hello/controllers"
	"log"
	"net/http"
	"strconv"
	"strings"

	// "github.com/robfig/cron/v3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var moodKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("😡", "/mood 1"),
			tgbotapi.NewInlineKeyboardButtonData("😐", "/mood 2"),
			tgbotapi.NewInlineKeyboardButtonData("🙂", "/mood 3"),
	),
	tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("😃", "/mood 4"),
			tgbotapi.NewInlineKeyboardButtonData("🤩", "/mood 5"),
	),
)

func RequestMood(bot *tgbotapi.BotAPI){
	reqMood := tgbotapi.NewMessage(5383565084, "How do you feel right now?")
	reqMood.ReplyMarkup = moodKeyboard

	bot.Send(reqMood)
	log.Printf("mood requested")

}

func PostMood(mood int, bot *tgbotapi.BotAPI) {
	values := map[string]int{"mood":mood}

	json, _ := json.Marshal(values)

	log.Printf("posting mood: %d", mood)
	resp, err := http.Post("http://localhost:8080/mood", "application/json", bytes.NewBuffer(json))	

	if err != nil {
		log.Printf("error posting mood")
	}
	if resp != nil {
		log.Printf("Logged Mood Successfully")
		MoodRes := tgbotapi.NewMessage(5383565084, "How do you feel right now?")
		bot.Send(MoodRes)
	}
}



func main(){

	bot, err := tgbotapi.NewBotAPI("5366512490:AAFNjdosYKeQofgp4BdI0ehUissp7-sIGRM")
	bot.Debug = true
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	controllers.SetJobs(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil {
				// Construct a new message from the given chat ID and containing
				// the text that we received.
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

				// If the message was open, add a copy of our numeric keyboard.
				switch update.Message.Command() {
				case "mood": {

					mood, err := strconv.Atoi(update.Message.CommandArguments())
					log.Printf("mood:%d parsed!\n", mood)

					if err != nil || mood < 1 || mood > 5 {
						msg.Text = "Mood value must be between 1 and 5"
						// bot.Send(msg)
						// panic(err)
					} else {
						values := map[string]int{"mood":mood}

						json, _ := json.Marshal(values)

						resp, err := http.Post("http://localhost:8080/mood", "application/json", bytes.NewBuffer(json))

						if err != nil {
							panic(err)
						} else {
							log.Printf(resp.Status)
							log.Printf("mood:%d added!\n", mood)
						}
					}
				}
				default: log.Printf("command not found")
				}

		} else if update.CallbackQuery != nil {
				// Respond to the callback query, telling Telegram to show the user
				// a message with the data received.
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				if _, err := bot.Request(callback); err != nil {
						panic(err)
				}

				log.Printf("callback registered")
				log.Printf("args: " + update.CallbackQuery.Message.CommandArguments())
				log.Printf("test: " + update.CallbackData())

				command := strings.Split(update.CallbackData(), " ")[0]
				args := strings.Join(strings.Split(update.CallbackData(), " ")[1:], " ")
				switch command {
					case "/mood": {
						mood, err := strconv.Atoi(args)
						if err != nil {
							panic(err)
						} else {
							PostMood(mood,bot)
						}
						
					}
				}

				// And finally, send a message containing the data received.
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				if _, err := bot.Send(msg); err != nil {
						panic(err)
				}
		}
	}
}

