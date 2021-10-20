package main

import (
	"log"
	"net/smtp"
)

func main() {

	from := "my_email@gmail.com"
	password := "super_secret_password"
	to := []string{"recipient@email.com"}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("My super secret message.")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}
}
