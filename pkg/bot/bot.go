package bot

import (
	"fmt"

	"github.com/bassiebal/ubiquiti-store-notifier/pkg/scraper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramCredentials struct {
	Token     string
	LogChatID int64
	ChatIDs   []int64
}

func SendUpdate(config *TelegramCredentials, product *scraper.Product) error {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return err
	}

	stockMsg := "Back in stock"
	if !product.Available {
		stockMsg = "Out of stock"
	}

	message := fmt.Sprintf(`<a href="%s"><b>%s</b></a> <b> %s </b> $%v`, fmt.Sprintf("https://store.ui.com/%s", product.Link), product.Name, stockMsg, product.Price)

	for _, chatID := range config.ChatIDs {
		msg := tgbotapi.NewMessage(chatID, message)
		msg.ParseMode = "HTML"

		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendError(config *TelegramCredentials, message error) error {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(config.LogChatID, fmt.Sprintf("Error: %v", message))
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
