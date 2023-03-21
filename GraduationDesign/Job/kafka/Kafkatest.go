package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

var address = []string{"43.139.123.205:9095"}

func Send() {
	// 配置
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应，只有上面的RequireAcks设置不是NoReponse这里才有用
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	// 版本
	config.Version = sarama.V0_10_2_1

	fmt.Println("start make producer")
	//使用配置，新建一个异步生产者
	producer, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		log.Printf("new async producer error: %s \n", err.Error())
		return
	}
	defer producer.AsyncClose()

	// 循环判断哪个通道发送过来数据
	fmt.Println("start goroutine")
	go func(p sarama.AsyncProducer) {
		for {
			select {
			case suc := <-p.Successes():
				fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
			case fail := <-p.Errors():
				fmt.Println("error: ", fail.Error())
			}
		}
	}(producer)

	var value string
	for i := 0; ; i++ {
		// 每隔两秒发送一条消息
		time.Sleep(2 * time.Second)

		// 创建消息
		value = fmt.Sprintf("async message, index = %d", i)
		// 注意：这里的msg必须得是新构建的变量，不然你会发现发送过去的消息内容都是一样的，因为批次发送消息的关系
		msg := &sarama.ProducerMessage{
			Topic: "web_log",
			Value: sarama.ByteEncoder(value),
		}

		// 使用通道发送
		producer.Input() <- msg
	}
}
