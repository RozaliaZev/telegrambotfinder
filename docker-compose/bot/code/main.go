package main

import (
	"botFinder/code/database"
	"botFinder/code/search"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"
    "strings"
	"github.com/Syfaro/telegram-bot-api"
)


var replyMarkup = tgbotapi.NewReplyKeyboard(
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("Статистика"),
    ),
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("Выбор города в России"),
    ),
	tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("Выбор города в Испании"),
    ),
	tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("Выбор города в США"),
    ),
	tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("Выбор города в Германии"),
    ),
)

func telegramBot() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}


		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":

				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, this bot will show you the weather at the moment in any city. Enter the name of the city!")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я могу показать вам погоду в любом городе мира. Просто укажите мне название города, и я дам вам подробную информацию о погоде в этом районе.")
				bot.Send(msg)

                msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
		        msg.ReplyMarkup = replyMarkup
                bot.Send(msg)

			case "/statistics","Статистика":

				statistics, pieDiagram, err := database.GetStatistics(update.Message.Chat.ID)

				if err != nil && pieDiagram == "" {
					log.Println(err)
					//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Chart output error, but bot still working.")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка вывода диаграммы, но вы можете продолжать пользоваться ботом.")
					bot.Send(msg)
				} else if err != nil {
					log.Println(err)
					//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error, but bot still working.")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка базы данных, но вы можете продолжать пользоваться ботом.")
					bot.Send(msg)
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, statistics)
				bot.Send(msg)

				data, _ := ioutil.ReadFile(pieDiagram)
				b := tgbotapi.FileBytes{Name: pieDiagram, Bytes: data}

				imgMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, b)
				bot.Send(imgMsg)

				err = os.Remove(pieDiagram)
				if err != nil {
					log.Println(err)
				}


            case "Выбор города в России":
                
                cityMarkup := tgbotapi.NewReplyKeyboard(
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("Москва"),
                        tgbotapi.NewKeyboardButton("Саратов"),
                        tgbotapi.NewKeyboardButton("Санкт-Петербург"),
                        tgbotapi.NewKeyboardButton("Казань"),
                    ),
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("назад"),
                    ),
                )
    
                cityMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберете город")
                cityMsg.ReplyMarkup = cityMarkup
    
                bot.Send(cityMsg)

			case "Выбор города в Испании":
                
                cityMarkup := tgbotapi.NewReplyKeyboard(
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("Мадрид"),
                        tgbotapi.NewKeyboardButton("Барселона"),
                        tgbotapi.NewKeyboardButton("Валенсия"),
                        tgbotapi.NewKeyboardButton("Гранада"),
                    ),
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("назад"),
                    ),
                )
    
                cityMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберете город")
                cityMsg.ReplyMarkup = cityMarkup
    
                bot.Send(cityMsg)

			case "Выбор города в США":
                
                cityMarkup := tgbotapi.NewReplyKeyboard(
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("Нью‑Йорк"),
                        tgbotapi.NewKeyboardButton("Лос‑Анджелес"),
                        tgbotapi.NewKeyboardButton("Чикаго"),
                        tgbotapi.NewKeyboardButton("Майами"),
                    ),
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("назад"),
                    ),
                )
    
                cityMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберете город")
                cityMsg.ReplyMarkup = cityMarkup
    
                bot.Send(cityMsg)

			case "Выбор города в Германии":
                
                cityMarkup := tgbotapi.NewReplyKeyboard(
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("Берлин"),
                        tgbotapi.NewKeyboardButton("Гамбург"),
                        tgbotapi.NewKeyboardButton("Мюнхен"),
                        tgbotapi.NewKeyboardButton("Лейпциг"),
                    ),
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("назад"),
                    ),
                )
    
                cityMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберете город")
                cityMsg.ReplyMarkup = cityMarkup
    
                bot.Send(cityMsg)

            case "назад":
                msgBack := tgbotapi.NewMessage(update.Message.Chat.ID, "Переход в главное меню")
                msgBack.ReplyMarkup = replyMarkup

                bot.Send(msgBack)

			default:

                if strings.HasPrefix(update.Message.Text, "/") {
                    command := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
                    bot.Send(command)
                }

				request := "https://www.google.com/search?q=погода+" + update.Message.Text
				message, err := search.Search(request)

				if err != nil || message == "" {
					//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Oops.. I don't know the weather in this city yet or I can't open the page")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Упс.. Погоду в этом городе я пока не могу узнать или я не могу открыть интернет-страницу")
					bot.Send(msg)
				}

				if os.Getenv("DB_SWITCH") == "on" {

					if err := database.CollectData(update.UpdateID, update.Message.Chat.ID, update.Message.Text, message); err != nil {
						//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error, but bot still working.")
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка базы данных, но вы можете продолжать пользоваться ботом.")
						bot.Send(msg)
					}

				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				bot.Send(msg)

			}
		} else {
            
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите название города.")
			bot.Send(msg)

		}
	}
}

func main() {

	time.Sleep(1 * time.Minute)

	if os.Getenv("CREATE_TABLE") == "yes" {

		if os.Getenv("DB_SWITCH") == "on" {

			if err := database.CreateTable(); err != nil {
				panic(err)
			}
		}
	}

	time.Sleep(1 * time.Minute)
	log.Println("start work")
	telegramBot()
}
