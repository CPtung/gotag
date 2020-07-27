package gotag

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	logger "github.com/sirupsen/logrus"
)

type TpMqtt struct {
	MsgQueueBase
	log    *logger.Logger
	c      mqtt.Client
	ch     chan *Tag
	topics []string
	ontag  OnTagCallback
}

func (self *TpMqtt) Publish(topic string, payload []byte) error {
	if self.c == nil {
		return errors.New("tag client not found")
	}
	token := self.c.Publish(topic, 0, false, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (self *TpMqtt) Subscribe(topic string) error {
	if self.c == nil {
		return errors.New("tag client not found")
	}
	if self.ontag == nil {
		return errors.New("needs to assign a handler before subscribe topics")
	}
	if token := self.c.Subscribe(topic, 0, self.onMessageReceived); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	// add subscribed topic
	for i := range self.topics {
		if self.topics[i] == topic {
			return nil
		}
	}
	self.topics = append(self.topics, topic)
	return nil
}

func (self *TpMqtt) UnSubscribe(topic string) error {
	if self.c == nil {
		return errors.New("tag client not found")
	}
	if token := self.c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	// remove unsubscribed topic
	for i := range self.topics {
		if self.topics[i] == topic {
			self.topics = append(self.topics[:i], self.topics[i+1:]...)
			break
		}
	}
	return nil
}

func (self *TpMqtt) SubscribeChannel() chan *Tag {
	return self.ch
}

func (self *TpMqtt) SubscribeCallback(hnd OnTagCallback) error {
	if self.c == nil {
		return errors.New("tag client not found")
	}
	self.ontag = hnd
	return nil
}

func (self *TpMqtt) Close() error {
	if self.c != nil {
		close(self.ch)
		self.c.Disconnect(250)
	}
	return nil
}

func (self *TpMqtt) onMessageReceived(client mqtt.Client, message mqtt.Message) {
	if self.ontag == nil {
		return
	}
	if _, _, err := DecodeTopic(message.Topic()); err != nil {
		self.log.Errorf("%v: %v", err.Error(), message.Topic())
		return
	}
	t := &Tag{}
	if err := DecodePayload(message.Payload(), t); err != nil {
		self.log.Errorf("on message received error (%v)", err.Error())
		return
	}
	if self.ontag != nil {
		self.ontag(t.SourceName, t.TagName, t.Val, t.ValType, t.Ts, t.Unit)
	}

	select {
	case self.ch <- t:
	default:
	}

	t = nil
}

func (self *TpMqtt) OnConnectHandler(client mqtt.Client) {
	self.log.Info("mqtt client connected")
	for i := range self.topics {
		if token := self.c.Subscribe(self.topics[i], 0, self.onMessageReceived); token.Wait() && token.Error() != nil {
			self.log.Warnf("re-subscribe topic (%v) error: %v", self.topics[i], token.Error())
		}
	}
}

func (self *TpMqtt) OnDisconnectHandler(client mqtt.Client, err error) {
	self.log.Info("mqtt client disconnected")
	if err != nil {
		self.log.Infof("disconnect error (%v)", err)
	}
}

func (self *TpMqtt) OnPublishHandler(client mqtt.Client, message mqtt.Message) {
	self.log.Infof("Publish topic: %v\n", message.Topic())
	self.log.Infof("Publish msg: %v\n", message.Payload())
}

func (self *TpMqtt) SetLogLevel(level string) error {
	if level == "info" {
		self.log.SetLevel(logger.InfoLevel)
	} else if level == "debug" {
		self.log.SetLevel(logger.DebugLevel)
	} else if level == "warn" {
		self.log.SetLevel(logger.WarnLevel)
	} else if level == "error" {
		self.log.SetLevel(logger.ErrorLevel)
	} else {
		return errors.New("Level not defined")
	}
	return nil
}

func newMqtt(cfg *MQConfig) (*TpMqtt, error) {
	t := &TpMqtt{
		topics: []string{},
		ontag:  nil,
		ch:     make(chan *Tag, 512),
		log:    logger.New(),
	}
	// set debug log level
	t.SetLogLevel(cfg.Debug)
	// init mqtt client
	opts := mqtt.NewClientOptions()
	if cfg.Host == "" {
		opts.AddBroker("tcp://" + getEnv("APPMAN_TAGSERVICE_ADDR", "localhost") + ":" + cfg.Port)
	} else {
		opts.AddBroker("tcp://" + cfg.Host + ":" + cfg.Port)
	}
	opts.SetClientID(genId(8))
	opts.SetKeepAlive(30 * time.Second)
	opts.SetMaxReconnectInterval(3)
	opts.SetCleanSession(true)
	opts.SetOnConnectHandler(t.OnConnectHandler)
	opts.SetConnectionLostHandler(t.OnDisconnectHandler)
	opts.SetDefaultPublishHandler(t.OnPublishHandler)

	t.c = mqtt.NewClient(opts)
	if token := t.c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return t, nil
}
