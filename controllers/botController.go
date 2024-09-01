package controllers

import (
	"fmt"
	"kal-bot-go/models"
	"kal-bot-go/utils"
	"kal-bot-go/views"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userId := message.From.ID
	models.StartNewSession(userId)
	msg := views.SendOrganizationPrompt(message.Chat.ID)
	bot.Send(msg)
}

func HandleTextMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userId := message.From.ID
	session := models.GetSession(userId)
	text := message.Text

	if text == "End Session" {
		models.EndSession(userId)
		bot.Send(views.SendSessionEnded(message.Chat.ID))
		models.StartNewSession(userId)
		bot.Send(views.SendOrganizationPrompt(message.Chat.ID))
		return
	}

	if text == "Next Stock Code" {
		if session.Step != "getStockCode" {
			bot.Send(views.SendNoActiveSession(message.Chat.ID))
			return
		}
		session.StockCode = ""
		bot.Send(views.SendNextStockCodePrompt(message.Chat.ID))
		return
	}

	switch session.Step {
	case "getOrganization":
		// Логика для выбора организации (KAL или KBL)
		if text == "KAL" || text == "KBL" {
			session.Organization = text
			session.Step = "getEmployeeNumber"
			bot.Send(views.SendEmployeePrompt(message.Chat.ID))
		} else {
			bot.Send(views.SendInvalidOrganization(message.Chat.ID))
		}
	case "getEmployeeNumber":
		// Логика для получения табельного номера сотрудника
		employeeName := models.GetEmployeeName(text)
		if employeeName != "" {
			session.Name = employeeName
			session.Step = "getStockCode"
			bot.Send(views.SendWelcome(message.Chat.ID, employeeName))
		} else {
			bot.Send(views.SendUnauthorized(message.Chat.ID))
		}
	case "getStockCode":
		// Логика для проверки сток-кода
		if len(text) == 6 {
			exists, err := models.CheckStockCodeInSheet(text)
			if err != nil {
				log.Printf("Error checking stock code: %v", err)
				return
			}
			if exists {
				bot.Send(views.SendStockCodeExists(message.Chat.ID))
			} else {
				session.StockCode = text
				session.Step = "getPhoto"
				bot.Send(views.SendPhotoPrompt(message.Chat.ID))
			}
		} else {
			bot.Send(views.SendInvalidStockCode(message.Chat.ID))
		}
	}
}

func HandlePhotoMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userId := message.From.ID
	session := models.GetSession(userId)

	if session.Step == "getPhoto" {
		photo := message.Photo[len(message.Photo)-1] // Берем фото с наибольшим разрешением
		fileLink, err := bot.GetFileDirectURL(photo.FileID)
		if err != nil {
			log.Printf("Error getting file link: %v", err)
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Ошибка при получении ссылки на файл"))
			return
		}

		filePath, err := models.DownloadFile(fileLink)
		if err != nil {
			log.Printf("Error downloading file: %v", err)
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Ошибка при загрузке файла"))
			return
		}

		err = models.CompressAndResizeImage(filePath)
		if err != nil {
			log.Printf("Error processing image: %v", err)
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Ошибка при обработке изображения"))
			return
		}

		newFileName := fmt.Sprintf("%s.jpg", session.StockCode)
		folderId := os.Getenv(session.Organization + "_FOLDER_ID")
		fileDriveId, err := models.UploadFileToDrive(filePath, newFileName, folderId)
		if err != nil {
			log.Printf("Error uploading file to Google Drive: %v", err)
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Ошибка при загрузке файла в Google Drive"))
			return
		}

		// Форматируем текущую дату как "дд/мм/гггг"
		date := time.Now().Format("02/01/2006")

		// Обновляем строку для записи в Google Sheets
		row := []interface{}{
			session.Name,      // ФИО
			session.StockCode, // Сток код
			date,              // Дата
			fmt.Sprintf("https://drive.google.com/uc?id=%s", fileDriveId), // Ссылка на Google Диск фото
			session.Organization, // Организация
		}

		err = models.AppendRowToSheet(row)
		if err != nil {
			log.Printf("Error appending row to Google Sheets: %v", err)
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Ошибка при записи данных в Google Sheets"))
			return
		}

		utils.DeleteFile(filePath)
		bot.Send(views.SendPhotoUploaded(message.Chat.ID))
		session.Step = "getStockCode"
	}
}
