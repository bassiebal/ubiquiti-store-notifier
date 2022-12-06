package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bassiebal/ubiquiti-store-notifier/pkg/bot"
	"github.com/bassiebal/ubiquiti-store-notifier/pkg/scraper"
)

type Config struct {
	Ubuiquiti scraper.UbiquitiCredentials
	Telegram  bot.TelegramCredentials
}

func GetConfig() *Config {
	ubiquitiUsername := os.Getenv("UBIQUITI_USERNAME")
	if ubiquitiUsername == "" {
		log.Fatalf("Required environment variable UBIQUITI_USERNAME not set")
	}

	ubiquitiPassword := os.Getenv("UBIQUITI_PASSWORD")
	if ubiquitiPassword == "" {
		log.Fatalf("Required environment variable UBIQUITI_PASSWORD not set")
	}

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Fatalf("Required environment variable TELEGRAM_TOKEN not set")
	}

	telegramChatIDsString := os.Getenv("TELEGRAM_CHAT_IDS")
	if telegramChatIDsString == "" {
		log.Fatalf("Required environment variable TELEGRAM_CHAT_IDS not set")
	}

	telegramChatIDsSlice := strings.Split(telegramChatIDsString, ",")

	telegramChatIDs := []int64{}
	for _, chatIDString := range telegramChatIDsSlice {
		chatID, err := strconv.ParseInt(chatIDString, 10, 64)
		if err != nil {
			log.Fatalf("Cannot convert chatID %s to integer", chatIDString)
		}
		telegramChatIDs = append(telegramChatIDs, chatID)
	}

	telegramLogChatIDString := os.Getenv("TELEGRAM_LOG_CHAT_ID")
	if telegramLogChatIDString == "" {
		log.Fatalf("Required environment variable TELEGRAM_LOG_CHAT_ID not set")
	}
	logChatID, err := strconv.ParseInt(telegramLogChatIDString, 10, 64)
	if err != nil {
		log.Fatalf("Cannot convert chatID %s to integer", telegramLogChatIDString)
	}

	return &Config{
		Ubuiquiti: scraper.UbiquitiCredentials{
			Username: ubiquitiUsername,
			Password: ubiquitiPassword,
		},
		Telegram: bot.TelegramCredentials{
			Token:     telegramToken,
			ChatIDs:   telegramChatIDs,
			LogChatID: logChatID,
		},
	}
}
