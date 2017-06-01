package email

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"os"
	"bytes"
	"mime/multipart"
	configPackage "github.com/AlstonWilliams/cobra/config"
	"fmt"
	"strings"
	"strconv"
)

type Notify_email_mailgun struct {
	config *configPackage.Config
}

func NewNotify_mailgun(config *configPackage.Config) Notify_email_mailgun{
	return Notify_email_mailgun{config}
}

type mailgun_return_structure struct {
	id string
	message string
}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func (mailgun Notify_email_mailgun) Notify(service_url, error_code string) {

	users := strings.Split(mailgun.config.Email, ",")
	for _, user := range users {
		mailgun.send(user, service_url, error_code)
	}

}

func (mailgun Notify_email_mailgun) send(user, service_url, error_code string){
	params := &bytes.Buffer{}
	writer := multipart.NewWriter(params)
	writer.WriteField("from", mailgun.config.Mailgun_from)
	writer.WriteField("to", user)
	writer.WriteField("subject", mailgun.config.Mailgun_subject)
	writer.WriteField("text", fmt.Sprintf(mailgun.config.Textformat, service_url, error_code))
	writer.Close()

	req, err := http.NewRequest("POST", mailgun.config.Mailgun_url, params)
	if err != nil {
		logrus.Error("Failed to make mailgun connection")
		return
	}
	req.SetBasicAuth("api", mailgun.config.Mailgun_api)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var i int = 1
	for i <= 3 {
		_, err := http.DefaultClient.Do(req)

		if err != nil {
			logrus.Error("Error occurs when send email to " + user + ", try " + strconv.Itoa(i) + " times")
			i++
			continue
		} else {
			logrus.Info("Successfully send mail to " + user)
			return
		}
	}
}