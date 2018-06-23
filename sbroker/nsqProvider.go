package sbroker

import (
	"fmt"
	"net/http"
	"os"

	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

var NSQConsumer *nsq.Consumer

type nsqProvider struct{}

func (np *nsqProvider) SubscribeMessage() error {
	c := utils.Config()
	if c.SBroker.NSQ.NsqlookupdHost != "" {
		config := nsq.NewConfig()
		channel := c.SBroker.NSQ.Channel
		hostname, err := os.Hostname()
		if err == nil {
			config.Hostname = hostname
			channel = hostname
		}

		NSQConsumer, err = nsq.NewConsumer(c.SBroker.NSQ.Topic, channel, config)
		if err != nil {
			logging.Log(zapcore.ErrorLevel, &logging.AppLog{
				Kind:    "messaging-subscribe",
				Message: err.Error(),
			})
			return errors.Wrap(err, "")
		}

		logging.Log(zapcore.InfoLevel, &logging.AppLog{
			Kind:    "messaging-subscribe",
			Message: fmt.Sprintf("%p", NSQConsumer),
		})

		NSQConsumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			// TODO
			return nil
		}))
		err = NSQConsumer.ConnectToNSQLookupd(c.SBroker.NSQ.NsqlookupdHost + ":" + c.SBroker.NSQ.NsqlookupdPort)
		if err != nil {
			logging.Log(zapcore.ErrorLevel, &logging.AppLog{
				Kind:    "messaging-subscribe",
				Message: err.Error(),
			})
			return errors.Wrap(err, "")
		}
	}

	return nil
}

func (np *nsqProvider) UnsubscribeMessage() error {
	if NSQConsumer == nil {
		return nil
	}

	c := utils.Config()
	hostname, err := os.Hostname()
	_, err = http.Post("http://"+c.SBroker.NSQ.NsqdHost+":"+c.SBroker.NSQ.NsqdPort+"/channel/delete?topic="+c.SBroker.NSQ.Topic+"&channel="+hostname, "text/plain", nil)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Kind:    "messaging-error",
			Message: err.Error(),
		})
		return errors.Wrap(err, "")
	}

	return nil
}
