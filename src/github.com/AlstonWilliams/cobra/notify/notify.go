package notify

import (
	"github.com/AlstonWilliams/cobra/notify/sms"
	"github.com/AlstonWilliams/cobra/notify/email"
	"github.com/AlstonWilliams/cobra/config"
)

const NOTIFY_TELEPHONE = "telephone"
const NOTIFY_EMAIL = "email"

type Notify interface {
	Notify(service_url, error_code string)
}

func NewNotify(telephoneOrEmail string, config *config.Config) Notify{
	if telephoneOrEmail == NOTIFY_TELEPHONE {
		return sms.NewNotify_yunpian(config)
	} else if telephoneOrEmail == NOTIFY_EMAIL{
		return email.NewNotify_mailgun(config)
	} else{
		return nil
	}
}