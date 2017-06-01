package callerd

import (
	"github.com/AlstonWilliams/cobra/stringutil"
	"net/http"
	"github.com/sirupsen/logrus"
	"strconv"
	"os"
	"net/url"
	"bytes"
)

const HTTP_METHOD_GET  = "get"
const HTTP_METHOD_POST = "post"

const HTTP_STATUS_CODE_ERROR = 400

type Caller_http struct {

}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp:true})
}

func (caller_http Caller_http) Call(method string, url string, params map[string]string) (bool, string){
	if method == HTTP_METHOD_GET {
		return caller_http.callGet(url, params)
	} else if method == HTTP_METHOD_POST {
		return caller_http.callPost(url, params)
	} else {
		logrus.Error("HTTP method which is not supported now")
		return true, ""
	}
}

func (Caller_http Caller_http) callGet(url string, params map[string]string) (bool, string){

	urlToExecute := url + "?"
	for k, v := range params  {
		urlToExecute += (k + "=" + v) + "&"
	}
	urlToExecute = stringutil.TrimSuffix(urlToExecute, "&")

	logrus.Info("Calling " + urlToExecute)

	resp, err := http.Get(urlToExecute)
	if err != nil {
		logrus.Error("Maybe service is exit " + urlToExecute)
		return false, "不能连接到服务"
	}

	if resp.StatusCode > HTTP_STATUS_CODE_ERROR {
		logrus.Error(urlToExecute + " returns " + strconv.Itoa(resp.StatusCode) + ", has notified the user")
		return false, string(resp.StatusCode)
	} else {
		logrus.Info(url + " is work normally")
	}

	return true, ""

}

func (Caller_http Caller_http) callPost(urlToExecute string, params map[string]string) (bool, string) {
	form := url.Values{}
	for key, value := range params {
		form.Add(key, value)
	}

	logrus.Info("Calling " + urlToExecute)

	body := bytes.NewBufferString(form.Encode())
	rsp, err := http.Post(urlToExecute, "application/x-www-form-urlencoded", body)
	if err != nil {
		logrus.Error("Maybe service is exit " + urlToExecute)
		return false, "不能连接到服务"
	}

	if rsp.StatusCode > HTTP_STATUS_CODE_ERROR {
		logrus.Error(urlToExecute + " returns " + strconv.Itoa(rsp.StatusCode) + ", has notified the user")
		return false, string(rsp.StatusCode)
	} else {
		logrus.Info(urlToExecute + " is work normally")
	}

	return true, ""
}
