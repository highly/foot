package rabbitMQ

import (
	"encoding/json"
	"fmt"
	"github.com/highly/foot/config"
	"github.com/streadway/amqp"
	"log"
)

const ContentTypePlain = "text/plain"
const ContentTypeJson = "application/json"

type CallBack func(content amqp.Delivery) bool

type exchangeInfo struct {
	exchangeName string
	queueName    string
	routingKey   string
	consumerTag  string
	handler      CallBack
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
	info    exchangeInfo
}

func New(exchangeGroupName string) (*RabbitMQ, error) {
	mq := &RabbitMQ{
		conn:    nil,
		channel: nil,
		done:    make(chan error),
	}

	var err error
	mq.conn, err = amqp.DialConfig(mq.amqpURI(), mq.dialConfig())
	if err != nil {
		return nil, fmt.Errorf("mq connection failed, [%s]", err)
	}
	mq.channel, err = mq.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("mq channel init failed, [%s]", err)
	}
	return mq.exchange(exchangeGroupName)
}

func (mq *RabbitMQ) Publish(content string) (err error, shutdown func() error) {
	if err := mq.channel.Publish(
		mq.info.exchangeName,
		mq.info.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  mq.getContentType(content),
			Body:         []byte(content),
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		return fmt.Errorf("mq exchange publish: %s", err), mq.publishClose
	}

	return nil, mq.publishClose
}

func (mq *RabbitMQ) Consume(callBack func(content amqp.Delivery) bool) (err error, shutdown func() error) {
	deliveries, err := mq.channel.Consume(
		mq.info.queueName, mq.info.consumerTag, false, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("mq queue consume failed, [%s]", err), mq.mqClose
	}

	mq.info.handler = callBack
	go mq.handleConsume(deliveries, mq.done)

	return nil, mq.mqClose
}

func (mq *RabbitMQ) amqpURI() string {
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		config.Get("rabbitmq.connection.default.user"),
		config.Get("rabbitmq.connection.default.password"),
		config.Get("rabbitmq.connection.default.host"),
		config.GetInt("rabbitmq.connection.default.port"),
	)
	log.Println("mq host url: ", amqpURI)
	return amqpURI
}

func (mq *RabbitMQ) dialConfig() amqp.Config {
	vhost := "/"
	configVhost := config.GetString("rabbitmq.connection.default.vhost")
	if configVhost != "" {
		vhost = configVhost
	}
	return amqp.Config{
		Vhost: vhost,
	}
}

func (mq *RabbitMQ) exchange(exchangeGroupName string) (*RabbitMQ, error) {
	if err := mq.exchangeCnf(exchangeGroupName); err != nil {
		return nil, err
	}
	if err := mq.exchangeDeclare(); err != nil {
		return nil, err
	}
	if err := mq.queueDeclare(); err != nil {
		return nil, err
	}
	if err := mq.queueBind(); err != nil {
		return nil, err
	}
	return mq, nil
}

func (mq *RabbitMQ) exchangeCnf(exchangeGroupName string) error {
	exchangeConfig := config.GetStringMapString(fmt.Sprintf("rabbitmq.exchange.%s", exchangeGroupName))
	if exchangeConfig == nil {
		return fmt.Errorf("mq exhange group [%s] not found", exchangeGroupName)
	}
	if exchangeConfig["exchangename"] == "" {
		return fmt.Errorf("exchnage name in group [%s] can not be empty", exchangeGroupName)
	}

	mq.info = exchangeInfo{
		exchangeName: exchangeConfig["exchangename"],
		queueName:    exchangeConfig["queuename"],
		routingKey:   exchangeConfig["routingkey"],
		consumerTag:  exchangeConfig["consumertag"],
	}
	return nil
}

func (mq *RabbitMQ) exchangeDeclare() error {
	if err := mq.channel.ExchangeDeclare(
		mq.info.exchangeName, "direct", true, false, false, false, nil,
	); err != nil {
		return fmt.Errorf("mq exchange Declare failed, [%s]", err)
	}
	return nil
}

func (mq *RabbitMQ) queueDeclare() error {
	if mq.info.queueName == "" {
		return nil
	}
	if _, err := mq.channel.QueueDeclare(
		mq.info.queueName, true, false, false, false, nil,
	); err != nil {
		return fmt.Errorf("mq queue Declare failed, [%s]", err)
	}
	return nil
}

func (mq *RabbitMQ) queueBind() error {
	if mq.info.queueName == "" {
		return nil
	}
	if err := mq.channel.QueueBind(
		mq.info.queueName, mq.info.routingKey, mq.info.exchangeName, false, nil,
	); err != nil {
		return fmt.Errorf("mq queue Bind failed, [%s]", err)
	}
	return nil
}

func (mq *RabbitMQ) isJson(body string) bool {
	var temp map[string]interface{}
	return json.Unmarshal([]byte(body), &temp) == nil
}

func (mq *RabbitMQ) getContentType(body string) string {
	if mq.isJson(body) {
		return ContentTypeJson
	}
	return ContentTypePlain
}

func (mq *RabbitMQ) mqClose() error {
	if err := mq.channel.Cancel(mq.info.consumerTag, true); err != nil {
		return fmt.Errorf("mq RabbitMQ cancel failed, [%s]", err)
	}

	if err := mq.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error, [%s]", err)
	}

	defer log.Println("AMQP shutdown OK")
	return <-mq.done
}

func (mq *RabbitMQ) publishClose() error {
	go func() {
		mq.done <- nil
	}()
	return mq.mqClose()
}

func (mq *RabbitMQ) handleConsume(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		if mq.info.handler(d) {
			d.Ack(false)
		}
	}
	done <- nil
}
