package main

import (
	"fmt"
	"github.com/highly/foot/config"
	"github.com/highly/foot/log"
	"github.com/highly/foot/orm"
	"github.com/highly/foot/rabbitMQ"
	"github.com/highly/foot/utils"
	"github.com/streadway/amqp"
	"path/filepath"
	"time"
)

type Mistake struct {
	orm.Model
	StuId      uint64 `json:"stu_id"`
	QuestionId uint64 `json:"question_id"`
}

func (t *Mistake) TableName() string {
	return "mistake"
}

func main() {
	usingConfig()

	//usingMQ()

	usingOrm()

	fmt.Println("blocking...")
	select {}
}

func usingOrm() {
	log.New(config.GetInt("logLevel"))
	db, _ := orm.NewGorm()

	var personList Mistake
	db.First(&personList)
	fmt.Println(personList)
}

func usingMQ() {
	mq1 := initMq("groupname1")
	if err, _ := mq1.Consume(func(content amqp.Delivery) bool {
		fmt.Println(
			"groupname1 => ",
			"content: ", string(content.Body),
			"consumerTag: ", content.ConsumerTag,
			"deliveryTag: ", content.DeliveryTag,
		)
		return true
	}); err != nil {
		panic("mq1 consume failed")
	}

	mq2 := initMq("groupname2")
	err, mq2Showdown := mq2.Consume(func(content amqp.Delivery) bool {
		fmt.Println(
			"groupname2 => ",
			"content: ", string(content.Body),
			"consumerTag: ", content.ConsumerTag,
			"deliveryTag: ", content.DeliveryTag,
		)
		return true
	})
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

	mq2Showdown()
	mq1Shutdown()
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
