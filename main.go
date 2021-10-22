package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

func main() {

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
	} else {
		fmt.Println("too expensive")
	}
}
