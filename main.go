package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

func main() {

	collector := colly.NewCollector(
		colly.AllowedDomains("amazon.com", "www.amazon.com"),
	)

	priceStr := ""

	collector.OnHTML("#priceblock_ourprice", func(element *colly.HTMLElement) {
		priceStr = element.Text
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})
	collector.Visit("https://www.amazon.com/Oculus-Quest-Advanced-All-One-Virtual/dp/B099VMT8VZ")

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

	fmt.Printf("email message: %s sent from %s to %s %s", message, from, to, password)

	// smtpHost := "smtp.gmail.com"
	// smtpPort := "587"

	// auth := smtp.PlainAuth("", from, password, smtpHost)

	// err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
