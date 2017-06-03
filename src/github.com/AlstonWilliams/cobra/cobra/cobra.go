package main

import (
	configPackage "github.com/AlstonWilliams/cobra/config"
	"path/filepath"
	"os"
	"github.com/sirupsen/logrus"
	"time"
	"strings"
	"github.com/AlstonWilliams/cobra/callerd"
	"github.com/AlstonWilliams/cobra/notify"
	"flag"
	"fmt"
)

var rules_root string
var rules []*configPackage.Rule

func main(){

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})


	configPath := flag.String("f", "./config.toml", "Config file path")

	getVersion := flag.Bool("v", false, "Show version")

	flag.Parse()

	if *getVersion == true{
		fmt.Println("cobra 1.0.0 (x86_64-pc-linux-gun)")
		os.Exit(0)
	}

	absolutePathOfConfigPath, err := filepath.Abs(*configPath)

	if err != nil{
		logrus.Error("failed to convert relative path to absolute path")
		os.Exit(0)
	}

	if _, err := os.Stat(absolutePathOfConfigPath); err != nil && os.IsNotExist(err) {
		logrus.Error("config file " + absolutePathOfConfigPath + " doesn't exist")
		os.Exit(1)
	}

	var config *configPackage.Config

	config = configPackage.NewConfig()
	parse_result := config.Parse(absolutePathOfConfigPath)
	if !parse_result {
		os.Exit(1)
	}

	rules_root = config.Rules_folder

	parseRules(config.Rules_folder)

	lastExecuteTime := time.Now().Add(- time.Minute * time.Duration(config.Interval))

	for true  {
		if time.Since(lastExecuteTime).Minutes() > float64(config.Interval){
			verifyInterfaceAndNotifyUserIfError(config)
			lastExecuteTime = time.Now()
		}
	}
}

func verifyInterfaceAndNotifyUserIfError(config *configPackage.Config){
	caller := callerd.NewCaller()

	for _, entry := range rules {
		go func(entry *configPackage.Rule) {
			params := make(map[string]string)
			urlParamItems := strings.Split(entry.Params, "&")
			for _, urlParamItems := range urlParamItems {
				urlParamKeyAndValue := strings.Split(urlParamItems, "=")
				params[urlParamKeyAndValue[0]] = urlParamKeyAndValue[1]
			}

			result, code := caller.Call(strings.ToLower(entry.Method), entry.Url, params)

			if result != true {
				if config.Telephone != "" {
					notify.NewNotify(notify.NOTIFY_TELEPHONE, config).Notify(entry.Url, code)
				}
				if config.Email != "" {
					notify.NewNotify(notify.NOTIFY_EMAIL, config).Notify(entry.Url, code)
				}
			}
		}(entry)
	}
}

func parseRules(path string) {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		logrus.Error("rules folder doesn't exist")
		os.Exit(1)
	}

	err := filepath.Walk(path, visit)
	if err != nil {
		logrus.Error("Error occurs when iterate the rules folder")
		os.Exit(1)
	}
}

func visit(path string, f os.FileInfo, err error) error {

	if path == rules_root{
		return nil
	}

	var rule *configPackage.Rule
	rule = configPackage.NewRule()
	rule.Parse(path)

	rules = append(rules, rule)

	return nil
}
