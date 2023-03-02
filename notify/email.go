package notify

import (
	"fmt"
	"net/smtp"
	"os"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
)

// EmailConfig holds all email configuration data
type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

// Email holds all email data
type Email struct {
	Recipient, Subject, HTMLBody, TextBody, Sender string
	Config                                         EmailConfig
}

// Prepare prepares the email notification, by using the passed booking, to set the Email fields
func (e *Email) Prepare(booking managerv1.Booking) Notifier {
	e.Subject = "Notice: Your resource instances will be stopped in 20 minutes."
	e.HTMLBody = fmt.Sprintf("<p>Your booking for resource <strong>%s</strong> expires in 20 minutes and the resource will be stopped. Please, extend the booking if you want to keep the resource instances running.</p>", booking.Spec.ResourceName)
	e.TextBody = fmt.Sprintf("Your booking for resource %s expires in 20 minutes and the resource will be stopped. Please, extend the booking if you want to keep the resource instances running.", booking.Spec.ResourceName)
	e.Sender = os.Getenv("SMTP_SENDER")

	e.Config.Host = os.Getenv("SMTP_HOST")
	e.Config.Port = os.Getenv("SMTP_PORT")
	e.Config.Username = os.Getenv("SMTP_USER")
	e.Config.Password = os.Getenv("SMTP_PASSWORD")

	return e
}

// Send sends the email notification
func (e *Email) Send() error {
	from := fmt.Sprintf("From: %s\r\n", e.Sender)
	to := fmt.Sprintf("To: %s\r\n", e.Recipient)
	subj := fmt.Sprintf("Subject: %s\r\n\r\n", e.Subject)
	body := fmt.Sprintf("%s\r\n", e.TextBody)

	msg := []byte(fmt.Sprintf("%s%s%s%s", from, to, subj, body))
	auth := smtp.PlainAuth("", e.Config.Username, e.Config.Password, e.Config.Host)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", e.Config.Host, e.Config.Port), auth, e.Sender, []string{e.Recipient}, msg)
	if err != nil {
		return err
	}

	return nil
}
