package gotag

type MQConfig struct {
	Id       string
	Host     string
	Port     string
	Debug    string
	Qos      byte
	Retained bool
}

type OnTagCallback func(
	sourceName string,
	tagName string,
	value *Value,
	valueType int32,
	timestamp uint64,
	unit string)

type MsgQueueBase interface {
	Publish(topic string, payload []byte) error
	Subscribe(topic string) error
	UnSubscribe(topic string) error
	SubscribeCallback(ontag OnTagCallback) error
	SubscribeChannel() chan *Tag
	Close() error
}
