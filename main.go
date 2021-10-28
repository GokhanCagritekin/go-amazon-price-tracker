package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

type Mail struct {
	From     string
	To       []string
	Password string
}

type Tracks struct {
	DesiredPrice float64 `json:"desiredPrice"`
}

func LoadMailFields() Mail {
	vi := viper.New()
	vi.SetConfigFile(".env")
	vi.ReadInConfig()
	From := vi.GetString("from")
	To := vi.GetStringSlice("to")
	Password := vi.GetString("password")
	return Mail{From: From, Password: Password, To: To}
}

func SendEmail(from string, password string, to []string, message []byte) {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}
}

func doEvery(ctx context.Context, d time.Duration, f func(time.Time)) error {
	ticker := time.Tick(d)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case x := <-ticker:
			go f(x)
		}
	}
}

func getDesiredPrice(str string) (desiredPrice float64) {
	arr := strings.Split(str, ":")
	arr2 := strings.Split(arr[1], "}")
	desiredPrice, err := strconv.ParseFloat(arr2[0], 64)
	if err != nil {
		log.Println(err)
	}
	return desiredPrice
}

func checkPrices(t time.Time) {
	client := GetRedisClient()

	keys, _ := client.Keys(ctx, "*http*").Result()
	for i := 0; i < len(keys); i++ {
		desiredPrice := getDesiredPrice(fmt.Sprint(client.Get(ctx, keys[i]).Result()))
		go check(keys[i], desiredPrice, client)
	}
}

func check(url string, desiredPrice float64, client *redis.Client) {
	collector := colly.NewCollector(
		colly.AllowedDomains("amazon.com", "www.amazon.com"),
	)

	priceStr := ""

	collector.OnHTML("#priceblock_ourprice", func(element *colly.HTMLElement) {
		priceStr = element.Text
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println(err.Error())
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})
	collector.Visit(url)

	price, err := strconv.ParseFloat(strings.ReplaceAll(priceStr, "$", ""), 64)
	if err != nil {
		log.Println(err)
	}

	if price < desiredPrice {
		mailFields := LoadMailFields()
		message := []byte("Price of the Product " + url + "is now lower than " + fmt.Sprint(desiredPrice))
		SendEmail(mailFields.From, mailFields.Password, mailFields.To, message)
		client.Del(ctx, url)
	}
}

func GetRedisClient() (c *redis.Client) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6000",
		Password: "",
		DB:       0,
	})
	return client
}

var ctx = context.Background()

func Addtrack() {

	client := GetRedisClient()

	desiredPrice, _ := strconv.ParseFloat(os.Args[1], 64)
	url := os.Args[2]

	_, err := client.Get(ctx, url).Result()

	if err == redis.Nil {
		track := Tracks{DesiredPrice: desiredPrice}

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
	}
}

func main() {
	if len(os.Args) > 1 {
		Addtrack()
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		doEvery(ctx, 2000*time.Millisecond, checkPrices)
	}
}
