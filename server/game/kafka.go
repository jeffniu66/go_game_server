package game

import (
	"github.com/Shopify/sarama"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"time"
)

var Producer *KafkaServer

type KafkaServer struct {
	Client            sarama.Client
	AsyncProducer     sarama.AsyncProducer
	SendTopic         *string
	PartitionConsumer sarama.PartitionConsumer
	ConsumeTopic      *string
}

func StartKafka() {
	Producer = &KafkaServer{}
	Producer.SendTopic = global.Game2Battle
	Producer.ConsumeTopic = global.Battle2Game
	config := sarama.NewConfig()

	address := global.MyConfig.Read("kafka", "address")
	client, err := sarama.NewClient([]string{address}, config)
	Producer.Client = client
	if err != nil {
		logger.Log.Errorln("kafka NewClient err: ", err)
	}

	asyncProducer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		logger.Log.Errorln("NewAsyncProducerFromClient err: ", err)
	}
	Producer.AsyncProducer = asyncProducer
	logger.Log.Infoln("-------------------- kafka start success!")
	if Producer.Client.Closed() {
		logger.Log.Errorln("init kafka 中断了，正在重连中...")
		time.AfterFunc(10*time.Second, StartKafka)
	}
}
