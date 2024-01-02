package email

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"
	"text/template"

	"github.com/Dukler/ChallengeStori/parser"
	"github.com/Dukler/ChallengeStori/resources"
)

type EmailService struct {
	from     string
	password string
	to       []string
	smtp     string
}

func NewEmailService() *EmailService {
	recipients := os.Getenv("EMAIL_RECIPIENTS")
	to := strings.Split(recipients, ",")
	return &EmailService{
		from:     os.Getenv("EMAIL_ADDRESS"),
		password: os.Getenv("EMAIL_PASSWORD"),
		to:       to,
		smtp:     "smtp.gmail.com",
	}
}

func (e *EmailService) SendSummary(sum *parser.Summarizer, recipients ...string) {
	auth := smtp.PlainAuth("", e.from, e.password, e.smtp)
	var to []string
	if len(recipients) > 0 {
		to = recipients
	} else {
		to = e.to
	}
	toHeader := strings.Join(to, ",")

	templatePath, err := resources.GetPath("mail_template.html")
	if err != nil {
		fmt.Println("Error opening template")
		return
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal(err)
	}

	data := getCreditCardData(sum)

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Fatal(err)
	}

	contentType := "Content-Type: text/html; charset=UTF-8"
	mime := "MIME-version: 1.0;\n"

	msg := []byte("To: " + toHeader + "\r\n" +
		"From: " + e.from + "\r\n" +
		"Subject: Your credit card statement\r\n" +
		contentType + "\r\n" +
		mime + "\r\n" +
		"\r\n" +
		buf.String())

	err = smtp.SendMail(net.JoinHostPort(e.smtp, "25"), auth, e.from, to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getCreditCardData(sum *parser.Summarizer) CreditCardData {
	var txns []Transaction
	for k, v := range sum.TxnsByMonth {
		tx := Transaction{Month: k, Amount: v}
		txns = append(txns, tx)
	}
	return CreditCardData{
		Transactions:    txns,
		TotalBalance:    float64(sum.Balance) / 100,
		AvgCreditAmount: float64(sum.AvgCredit) / 100,
		AvgDebitAmount:  float64(sum.AvgDebit) / 100,
	}
}

type CreditCardData struct {
	Transactions    []Transaction
	TotalBalance    float64
	AvgCreditAmount float64
	AvgDebitAmount  float64
}

type Transaction struct {
	Month  string
	Amount int
}
