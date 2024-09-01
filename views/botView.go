package views

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendOrganizationPrompt(chatID int64) tgbotapi.MessageConfig {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("KAL"),
			tgbotapi.NewKeyboardButton("KBL"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("End Session"),
			tgbotapi.NewKeyboardButton("Next Stock Code"),
		),
	)

	// Создаем сообщение и присваиваем клавиатуру
	msg := tgbotapi.NewMessage(chatID, "Выберите организацию?")
	msg.ReplyMarkup = keyboard

	return msg
}

func SendEmployeePrompt(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Введите табельный номер сотрудника:")
}

func SendWelcome(chatID int64, employeeName string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Добро пожаловать, "+employeeName+". Введите сток код.")
}

func SendInvalidOrganization(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Пожалуйста, укажите организацию: KAL или KBL.")
}

func SendUnauthorized(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Вы не авторизованы. Обратитесь к администратору.")
}

func SendInvalidStockCode(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Неверный сток код. Пожалуйста, убедитесь, что сток код состоит из 6 цифр и введите его заново.")
}

func SendNextStockCodePrompt(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Введите следующий сток код.")
}

func SendStockCodeExists(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Этот сток код уже существует. Пожалуйста, введите следующий сток код.")
}

func SendPhotoPrompt(chatID int64) tgbotapi.MessageConfig {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("End Session"),
			tgbotapi.NewKeyboardButton("Next Stock Code"),
		),
	)

	// Создаем сообщение и присваиваем клавиатуру
	msg := tgbotapi.NewMessage(chatID, "Сток код сохранен. Пожалуйста, отправьте фотографию.")
	msg.ReplyMarkup = keyboard

	return msg
}

func SendPhotoUploaded(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Фотография успешно загружена и сохранена. Пожалуйста, введите следующий сток код или нажмите завершить сессию.")
}

func SendSessionEnded(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Сеанс завершен. Запускаю новый сеанс. Пожалуйста, укажите организацию.")
}

func SendNoActiveSession(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, "Нет активного сеанса. Для начала нового сеанса используйте команду /start.")
}
