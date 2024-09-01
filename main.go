package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// Подключение телеграм API
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	// Подключение dotenv
	"github.com/joho/godotenv"
	// Подключение конфигов
	"kal-bot-go/config"
	// Подключение контроллеров
	"kal-bot-go/controllers"
)

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Инициализация Telegram-бота
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	// Инициализация сервисов Google
	config.InitGoogleServices()

	// Основной цикл обработки обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // игнорируем не Message обновления
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				controllers.HandleStartCommand(bot, update.Message)
			}
		} else if update.Message.Photo != nil {
			controllers.HandlePhotoMessage(bot, update.Message)
		} else {
			controllers.HandleTextMessage(bot, update.Message)
		}
	}

	// Можно добавить HTTP-сервер, если нужен, чтобы программа не завершалась
	fmt.Println("Запуск HTTP-сервера на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
