package emails

import (
	"bytes"
	"html/template"
	"net/smtp"
)

var (
	senderEmail    = "misterphoenix6@gmail.com"
	senderPassword = "mister_66"
	senderHost     = "smtp.gmail.com"
)

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func sendEmail(to []string, subjectLine, body string) (bool, error) {
	auth := smtp.PlainAuth("", senderEmail, senderPassword, senderHost)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + subjectLine + "!\n"
	msg := []byte(subject + mime + body)
	addr := "smtp.gmail.com:587"

	toAddresses := to
	if err := smtp.SendMail(addr, auth, senderEmail, toAddresses, msg); err != nil {
		return false, err
	}
	return true, nil
}