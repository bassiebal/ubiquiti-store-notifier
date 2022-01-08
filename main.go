package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/bassiebal/ubiquiti-store-notifier/pkg/bot"
	"github.com/bassiebal/ubiquiti-store-notifier/pkg/database"
	"github.com/bassiebal/ubiquiti-store-notifier/pkg/scraper"
)

func main() {
	config := getConfig()

	db, err := database.Connect("./database.db")
	if err != nil {
		logError(err)
		log.Fatal(err)
	}
	defer db.Close()

	products, err := scraper.GetProducts(config.Ubuiquiti)
	if err != nil {
		logError(err)
		log.Fatal(err)
	}
	log.Printf("Retrieved %d products from the store", len(products))

	updateTimestamp := time.Now().Unix()
	for _, product := range products {
		dbProduct := scraper.Product{}
		err = db.Get(&dbProduct, `
			SELECT name, price, available, link
			FROM products
			WHERE name = ?
			ORDER BY inserted_at DESC
			LIMIT 1`, product.Name)
		if err != nil && err != sql.ErrNoRows {
			logError(fmt.Errorf("Unable to retrieve product from database: %v", err))
		}

		if reflect.DeepEqual(product, dbProduct) {
			continue
		}

		log.Printf("Change for product: %s, with price: %v and availability: %v", product.Name, product.Price, product.Available)
		_, err = db.NamedExec(fmt.Sprintf(`
				INSERT INTO products (name, price, link, available, inserted_at)
				VALUES (:name, :price, :link, :available, %d)
			`, updateTimestamp), product)
		if err != nil {
			logError(fmt.Errorf("Error inserting product: %v", err))
		}

		if product.Available && !dbProduct.Available {
			err = bot.SendUpdate(&config.Telegram, &product)
			if err != nil {
				logError(fmt.Errorf("Error sending telegram update: %v", err))
			}
		}

	}

}

func logError(message error) {
	err := bot.SendError(&getConfig().Telegram, message)
	if err != nil {
		log.Printf("Could not send error to telegram: %v\n", err)
	}
	log.Println(err)
}
