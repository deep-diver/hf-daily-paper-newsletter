package internal

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
)

type Head struct {
	Title                 string
	Description           string
	ImageURL              string
	CommunityLink         string
	CommunityLinkBtnTitle string
	BgColor               string
}

type Section struct {
	Title string
}

type ArticleTuple struct {
	Article1  Article
	Article2  Article
	LinkTitle string
	BgColor   string
}

type Email struct {
	Title         string
	FooterTitle   string
	Header        Head
	FirstSection  Section
	ArticleTuples []ArticleTuple // has to be
}

type Request struct {
	from     string
	to       []string
	subject  string
	body     string
	password string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func NewRequest(from string, password string, to []string, subject string) *Request {
	return &Request{
		from:     from,
		to:       to,
		subject:  subject,
		password: password,
	}
}

func (r *Request) parseTemplate(templatePath string, data interface{}) error {
	templateFilenames := GetTemplatesInDir(templatePath)

	tmpl, _ := template.ParseFiles(templateFilenames...)

	buffer := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buffer, "template.gohtml", data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail() bool {
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%d", "smtp.gmail.com", 587)

	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", r.from, r.password, "smtp.gmail.com"), "deep.diver.csp@gmail.com", r.to, []byte(body)); err != nil {
		fmt.Println(err)
		fmt.Println("fail on smtp.SendMail")
		return false
	}
	return true
}

func (r *Request) Send(templatePath string, items interface{}) {
	err := r.parseTemplate(templatePath, items)
	if err != nil {
		log.Fatal(err)
	}

	if ok := r.sendMail(); ok {
		fmt.Printf("Email has been sent to   %s\n", r.to)
	} else {
		fmt.Printf("Failed to send the email to %s\n", r.to)
	}
}
