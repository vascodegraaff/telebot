package controllers

import (
	"encoding/json"
	"example/user/hello/models"
	"io/ioutil"
	"log"
	"github.com/google/uuid"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
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

var rangeKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("1"),
			tgbotapi.NewKeyboardButton("2"),
			tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("4"),
			tgbotapi.NewKeyboardButton("5"),
	),
)

func SendMessage(bot *tgbotapi.BotAPI, text string, replyType string){
	message := tgbotapi.NewMessage(5383565084, text)
	switch replyType {
	case "range":
		message.ReplyMarkup = rangeKeyboard
	}
	// case "number":
	// 	message.ReplyMarkup = 
	// message.ReplyMarkup = moodKeyboard

	bot.Send(message)
	log.Printf("Message sent: " + text)

}

func SetJobs(bot *tgbotapi.BotAPI) {
	file, err := ioutil.ReadFile("/Users/vasco/Projects/bot/question.json")
	if err != nil {
		panic("unable to read file")
	}

	var questionSet []models.QuestionSet
	_ = json.Unmarshal([]byte(file), &questionSet)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	c := cron.New()

	for i, set := range questionSet {
		questionSet[i].Id = uuid.New()
		log.Printf("set_name: %s\n", set.Set_name)
		log.Printf("description: %s\n", set.Description)
		log.Printf("schedule type: %s\n", set.Schedule.T)
		for j, question := range set.Questions {
			questionSet[i].Questions[j].Id = uuid.New()
			switch questionSet[i].Schedule.T{
			case "cron":
				c.AddFunc(questionSet[i].Schedule.Value,  func() {
					SendMessage(bot, question.Question, question.ReplyType)
					log.Printf("cron job executed")
				})
			// case "random":
				
			}
			log.Printf("question id: %s\n", questionSet[i].Questions[j].Id)
			log.Printf("question: %s\n", question.Question)
			log.Printf("reply type: %s\n", question.ReplyType)
		}
	}
	c.Start()
}