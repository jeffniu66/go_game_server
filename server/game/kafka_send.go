package game

import (
	"github.com/Shopify/sarama"
	"go_game_server/server/logger"
	"runtime/debug"
)

func KafkaSend(msgData string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorln("kafka_send socket error: ", err, string(debug.Stack())) // 这里的err其实就是panic传入的内容
			debug.PrintStack()
		}
	}()
	asyncProducer := Producer.AsyncProducer
	//data, _ := proto.Marshal(msgData)
	msg := sarama.StringEncoder(msgData)
	asyncProducer.Input() <- &sarama.ProducerMessage{Topic: *Producer.SendTopic, Value: msg}
	logger.Log.Infoln("kafka send message : ", msg)
}
