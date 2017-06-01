package sms

import (
	"net/url"
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"os"
	"encoding/json"
	configPackage "github.com/AlstonWilliams/cobra/config"
	"fmt"
)

type Notify_sms_yunpian struct {
	config	*configPackage.Config
}

func NewNotify_yunpian(config *configPackage.Config) Notify_sms_yunpian{
	return Notify_sms_yunpian{config}
}

type yunpian_return_structure_result struct {
	count int
	fee float64
	sid string
}

type yunpian_return_structure struct {
	code int
	msg string
	result yunpian_return_structure_result
}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func (yunpian Notify_sms_yunpian) Notify(service_url, error_code string) {
	urlToExecute := yunpian.config.Yunpian_url

	form := url.Values{}
	form.Add("apikey", yunpian.config.Yunpian_apikey)
	form.Add("mobile", yunpian.config.Telephone)
	form.Add("text", fmt.Sprintf(yunpian.config.Textformat, service_url, error_code))

	body := bytes.NewBufferString(form.Encode())
	rsp, err := http.Post(urlToExecute, "application/x-www-form-urlencoded", body)
	if err != nil{
		logrus.Error("Error occurs when send sms to user")
		return
	}
	defer rsp.Body.Close()

	body_byte, err := ioutil.ReadAll(rsp.Body)
	if err != nil{
		logrus.Error("Error occurs when get yunpian's response")
	}

	var yunpian_response yunpian_return_structure
	err = json.Unmarshal(body_byte, &yunpian_response)
	if err != nil {
		logrus.Error("Can't decode yunpian's response")
	} else {
		if yunpian_response.code != 0 {
			logrus.Error("Un-correct way to call yunpian: code is " + string(yunpian_response.code) + " msg is : " + yunpian_response.msg)
		}
	}


}
