package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {

	v1 := viper.New()
	v1.SetConfigFile(".env")
	v1.ReadInConfig()

	from := v1.GetString("from")
	password := v1.GetString("password")
	to := v1.GetStringSlice("to")

	message := []byte("My super secret message.")

	fmt.Printf("email message %s sent from %s to %s %s", message, from, to, password)

	// smtpHost := "smtp.gmail.com"
	// smtpPort := "587"

	// auth := smtp.PlainAuth("", from, password, smtpHost)

	// err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
