package mail

import (
	"github.com/MeM0rd/q-api-gateway/pkg/logger"
	"log"
	"net/smtp"
	"os"
)

type EmailSender interface {
	Send(
		to []string,
		msg []byte,
	)
}

type GmailSender struct {
	EmailAddress  string
	EmailPassword string
	AuthAddress   string
	ServerAddress string
	Logger        *logger.Logger
}

func NewGmailSender() *GmailSender {
	return &GmailSender{
		EmailAddress:  os.Getenv("GMAIL_NAME"),
		EmailPassword: os.Getenv("GMAIL_PASS"),
		AuthAddress:   os.Getenv("GMAIL_AUTH_ADDR"),
		ServerAddress: os.Getenv("GMAIL_SERVER_ADDR"),
		Logger:        logger.New(),
	}
}

func PrepareGmail(to []string, template string) {
	gs := NewGmailSender()

	t, err := findTemplate(template)
	if err != nil {
		gs.Logger.Infof("error template finding: %v", err)
		return
	}

	gs.Send(to, []byte(t))
}

func (gs *GmailSender) Send(to []string, msg []byte) {
	auth := smtp.PlainAuth("",
		gs.EmailAddress,
		gs.EmailPassword,
		gs.AuthAddress,
	)

	err := smtp.SendMail(
		gs.ServerAddress,
		auth,
		gs.EmailAddress,
		to,
		msg,
	)
	if err != nil {
		log.Printf("error sending email")
		return
	}
}
