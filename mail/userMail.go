package mail

import (
	"Enterprise/data"
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	_ "gopkg.in/gomail.v2"
	"html/template"
	"net/smtp"
	"os"
	"strconv"
)

var from = os.Getenv("MAIL_FROM")
var password = os.Getenv("MAIL_PASSWORD")
var smtpHost = os.Getenv("MAIL_HOST")
var smtpPort = os.Getenv("MAIL_PORT")

var mimeHeaders = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func EmailFromGoMail(email string, subject string, mailInputs *data.MailInputs, mailPath string) error {
	var apiKey = os.Getenv("BREVO_API_KEY")
	var brevoFrom = os.Getenv("BREVO_FROM")
	var brevoHost = os.Getenv("BREVO_HOST")
	var brevoPort = os.Getenv("BREVO_PORT")

	tmpl, err := template.ParseFiles(mailPath)
	if err != nil {
		return err
	}

	data := struct {
		Email    string
		Code     string
		Username string
		Link     string
	}{
		Email:    mailInputs.Email,
		Code:     mailInputs.Code,
		Username: mailInputs.Username,
		Link:     mailInputs.Link,
	}

	// Generate the email body by applying the template with the data
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", brevoFrom)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body.String())

	port, _ := strconv.Atoi(brevoPort)
	n := gomail.NewDialer(brevoHost, port, brevoFrom, apiKey)
	if err := n.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func EmailLogics(subject, templatePath string, emailDto *data.MailInputs, templateData interface{}) error {

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	to := []string{emailDto.Email}

	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	err = t.Execute(&body, templateData)
	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpHost+":"+string(smtpPort), auth, from, to, body.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func VerifyEmail(emailDto *data.MailInputs) error {
	templateData := struct {
		Code     string
		Username string
		Link     string
	}{
		Code:     emailDto.Code,
		Username: emailDto.Username,
		Link:     os.Getenv("FRONTEND_URL") + "/verify",
	}
	return EmailLogics("Verify Account", "mail/templates/verify.html", emailDto, templateData)
}

func ResetPassword(emailDto *data.MailInputs) error {
	templateData := struct {
		Email    string
		Code     string
		Username string
		Link     string
	}{
		Email:    emailDto.Email,
		Code:     emailDto.Code,
		Username: emailDto.Username,
		Link:     os.Getenv("BACKEND_URL") + "api/user/reset-password?code=" + emailDto.Code,
	}
	return EmailLogics("Reset Password", "mail/templates/reset.html", emailDto, templateData)
}
