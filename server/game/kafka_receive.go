package game

import (
	"github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
	"go_game_server/server/global"
	"go_game_server/server/logger"
)

func KafkaReceive() {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true     // 必须读取消费通道的返回值
	config.Group.Return.Notifications = true // 必须从相应的通道读取内容

	// init consumer
	address := global.MyConfig.Read("kafka", "address")
	brokers := []string{address}
	topics := []string{*Producer.ConsumeTopic}
	consumer, err := cluster.NewConsumer(brokers, *Producer.ConsumeTopic, topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	// consume messages, watch signals
	for {
		select {
		case err = <-consumer.Errors():
			logger.Log.Errorf("group(%s) topic(%s)  catch error(%v)", *Producer.ConsumeTopic, *Producer.ConsumeTopic, err)
			//return
		case ntf := <-consumer.Notifications():
			log.Printf("Rebalanced: %+v\n", ntf)
		case msg, _ := <-consumer.Messages():
			// 手动确认:主题分区的偏移量标记为已处理
			consumer.MarkPartitionOffset(*Producer.ConsumeTopic, msg.Partition, msg.Offset, "")
			logger.Log.Infoln("kafka receive message : ", msg)
		}
	}
}
