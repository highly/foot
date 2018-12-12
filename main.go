package main

import (
	"fmt"
	"github.com/highly/foot/config"
	"github.com/highly/foot/rabbitMQ"
	"github.com/highly/foot/utils"
	"github.com/streadway/amqp"
	"path/filepath"
	"time"
)

func main() {
	usingConfig()

	usingMQ()

	fmt.Println("blocking...")
	select {}
}

func usingMQ() {
	mq1 := initMq("groupname1")
	handler := func(content amqp.Delivery) bool {
		fmt.Println(
			"groupname1 => ",
			"content: ", string(content.Body),
			"consumerTag: ", content.ConsumerTag,
			"deliveryTag: ", content.DeliveryTag,
		)
		return true
	}
	if err, _ := mq1.Consume(handler); err != nil {
		panic("mq1 consume failed")
	}

	mq2 := initMq("groupname2")
	handler2 := func(content amqp.Delivery) bool {
		fmt.Println(
			"groupname2 => ",
			"content: ", string(content.Body),
			"consumerTag: ", content.ConsumerTag,
			"deliveryTag: ", content.DeliveryTag,
		)
		return true
	}
	err, mq2Showdown := mq2.Consume(handler2)
	if err != nil {
		panic("mq2 consume failed")
	}

	dataList := []string{"message1", "message2", "message3"}
	for _, v := range dataList {
		err, _ = mq1.Publish(v)
		if err != nil {
			fmt.Println("mq1 publish failed: ", err.Error())
		}
		time.Sleep(2 * time.Second)
	}
	_, mq1Shutdown := mq1.Publish("closing")

	mq1Shutdown()
	mq2Showdown()
}

func initMq(exchangeGroupName string) *rabbitMQ.RabbitMQ {
	mq, err := rabbitMQ.New(exchangeGroupName)
	if err != nil {
		panic("rabbitMQ connection init failed")
	}
	return mq
}

func usingConfig() {
	config.New().EnvPrefix("TEST").ConfigPath(filepath.Join(utils.BaseDir(), "config")).Load("cnf")
	// effective immediately
	//config.New().ConfigPath(filepath.Join(BaseDir(), "config")).Load("cnf").Watching()
	fmt.Println(config.GetStringMapString("rabbitmq.exchange.default")["exchangename"])
}
