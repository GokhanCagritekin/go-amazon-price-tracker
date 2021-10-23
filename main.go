package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

type Tracks struct {
	DesiredPrice float64 `json:"desiredPrice"`
}

var ctx = context.Background()

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:        "localhost:6000",
		Password:    "",
		DB:          0,
		IdleTimeout: 5 * time.Minute,
	})

	desiredPrice, _ := strconv.ParseFloat(os.Args[1], 64)
	url := os.Args[2]

	collector := colly.NewCollector()

	priceStr := ""

	collector.OnHTML("#priceblock_ourprice", func(element *colly.HTMLElement) {
		priceStr = element.Text
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})
	collector.Visit(url)

	price, err := strconv.ParseFloat(strings.ReplaceAll(priceStr, "$", ""), 64)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("price is %v\n", price)

	v1 := viper.New()
	v1.SetConfigFile(".env")
	v1.ReadInConfig()

	from := v1.GetString("from")
	password := v1.GetString("password")
	to := v1.GetStringSlice("to")

	message := []byte("My super secret message.")

	if price < desiredPrice {
		fmt.Printf("email message: %s sent from %s to %s %s", message, from, to, password)
		// smtpHost := "smtp.gmail.com"
		// smtpPort := "587"

		// auth := smtp.PlainAuth("", from, password, smtpHost)

		// err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		client.Del(ctx, url)
	} else {
		fmt.Println("too expensive")

		_, err := client.Get(ctx, url).Result()
		if err == redis.Nil {
			fmt.Printf("%s does not exist\n", url)
			track := Tracks{DesiredPrice: desiredPrice}
			fmt.Println(track)
			json, err := json.Marshal(track)
			if err != nil {
				log.Println(err)
			}

			err = client.Set(ctx, url, json, 0).Err()

			if err != nil {
				log.Println(err)
			}
		} else if err != nil {
			panic(err)
		} else {
			fmt.Printf("%s does exist\n", url)
		}
	}
	keys, _ := client.Keys(ctx, "*http*").Result()
	for i := 0; i < len(keys); i++ {
		fmt.Println(keys[i], client.Get(ctx, keys[i]))
	}

}
