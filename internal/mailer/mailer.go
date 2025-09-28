package mailer

import (
	"bytes"
	"embed"
	"time"

	ht "html/template"
	tt "text/template"

	"github.com/wneessen/go-mail"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	client *mail.Client
	sender string
}

func New(host string, port int, username, password, sender string) (*Mailer, error) {
	client, err := mail.NewClient(
		host,
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	// Return a Mailer instance containing the client and sender information.
	mailer := &Mailer{
		client: client,
		sender: sender,
	}
	return mailer, nil
}

func (m *Mailer) Send(recipient string, templateFile string, data any) error {
	// Use the ParseFS() method from text/template to parse the required template file
	// from the embedded file system.
	textTmpl, err := tt.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Execute the named template "subject", passing in the dynamic data and storing the
	// result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Follow the same pattern to execute the "plainBody" template and store the result
	// in the plainBody variable.
	plainBody := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// Use the ParseFS() method from html/template this time to parse the required template
	// file from the embedded file system.
	htmlTmpl, err := ht.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// And execute the "htmlBody" template and store the result in the htmlBody variable.
	htmlBody := new(bytes.Buffer)
	err = htmlTmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Use the mail.NewMsg() function to initialize a new mail.Msg instance.
	// Then we use the To(), From() and Subject() methods to set the email recipient,
	// sender and subject headers, the SetBodyString() method to set the plain-text body,
	// and the AddAlternativeString() method to set the HTML body.
	msg := mail.NewMsg()
	err = msg.To(recipient)
	if err != nil {
		return err
	}
	err = msg.From(m.sender)
	if err != nil {
		return err
	}
	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())

	// Retry 3 times on sending email if fail.
	for i := 1; i <= 3; i++ {
		err = m.client.DialAndSend(msg)
		if err == nil {
			return nil
		}
		// If it didn't work, sleep for a short time and retry.
		if i != 3 {
			time.Sleep(500 * time.Millisecond)
		}
	}
	return err
}
