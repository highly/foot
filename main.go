package main

import (
	"fmt"
	"github.com/highly/foot/config"
	"github.com/highly/foot/log"
	"github.com/highly/foot/orm"
	"github.com/highly/foot/rabbitMQ"
	"github.com/highly/foot/utils"
	"github.com/streadway/amqp"
	"os"
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
	initConfig()

	initLog()

	// usingMQ()

	// usingOrm()

	fmt.Println("blocking...")
	select {}
}

func usingOrm() {
	//log.New(config.Int("logLevel"))
	//db, _ := orm.NewGorm()
	//
	//var personList Mistake
	//db.First(&personList)
	//fmt.Println(personList)
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

func initLog() {
	hostname, _ := os.Hostname()

	opts := []log.Option{
		log.WithField("hostname", hostname),
		log.WithField("hostip", utils.Ip()),
		log.WithField("environment", config.String("environment")),
		log.WithField("service_name", config.String("serviceName")),
		log.WithField("version", config.String("version")),
		log.WithLevel(log.ToLevel(config.DefaultString("log.level", "info"))),
		log.WithPath(config.DefaultString("log.path", "test.log")),
		log.WithMaxSize(config.DefaultInt("log.maxSize", 50)),
		log.WithMaxBackups(config.DefaultInt("log.maxBackups", 5)),
		log.WithMaxAge(config.DefaultInt("log.age", 5)),
	}

	log.New(opts...)

	log.Info("log test")
}

func initConfig() {
	// 1)
	// config.EnvPrefix("TEST")
	// config.ConfigPath(filepath.Join(utils.BaseDir(), "config"))
	// config.Load("cnf")

	// 2)
	config.EnvPrefix("TEST")
	config.LoadFile(filepath.Join(utils.BaseDir(), "cnf.yaml"))

	// plus) effective immediately
	// config.Watching()

	fmt.Println("testing config: ", config.StringMap("rabbitmq.exchange.GroupName1")["exchangename"])
}
