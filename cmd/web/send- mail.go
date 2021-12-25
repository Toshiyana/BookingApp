package main

import (
	"log"
	"time"

	"github.com/Toshiyana/BookingApp/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	// listen to the mail channel
	// Using go routine, fires asynchronously in background
	// Every time, I can get and send messages without stopping application
	// (Because sending a mail is take times, about 10 seconds)

	go func() {
		for {
			// using inifinity for loop, listen all the time for incoming data
			msg := <-app.MailChan // listen to the mail channel
			sendMsg(msg)
		}
	}()

}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost" // install a dummy mail server on local machines for practice
	server.Port = 1025        // port for dummy server

	// Whether keeping connection to the email server all the time or not
	// choose false because we only make a connection when sending and giving email
	server.KeepAlive = false

	// if you can't connect within 10 seconds, give up connecting and sending
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// In production, you set the below values, but we don't because of a dummy server
	// server.Username =
	// server.Password =
	// server.Encryption =

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent!")
	}

}
