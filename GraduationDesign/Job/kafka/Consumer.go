package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
)

type AAAConsumerGroupHandler struct{}

func (AAAConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (AAAConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// 这个方法用来消费消息的
func (h AAAConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 获取消息
	for msg := range claim.Messages() {
		fmt.Printf("topic:%q partition:%d offset:%d value:%s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		// 将消息标记为已使用
		sess.MarkMessage(msg, "")
	}
	return nil
}

// 接收数据
func Consumer1() {
	// 先初始化 kafka
	config := sarama.NewConfig()
	// Version 必须大于等于  V0_10_2_0
	config.Version = sarama.V0_10_2_1
	config.Consumer.Return.Errors = true
	fmt.Println("start connect kafka")
	// 开始连接kafka服务器
	group, err := sarama.NewConsumerGroup([]string{"43.139.123.205:9095"}, "AAA-group", config)

	if err != nil {
		fmt.Println("连接kafka失败：", err)
		return
	}
	// 检查错误
	go func() {
		for err := range group.Errors() {
			fmt.Println("分组错误 : ", err)
		}
	}()

	ctx := context.Background()
	fmt.Println("开始获取消息")
	// for 是应对 consumer rebalance
	for {
		// 需要监听的主题
		topics := []string{"web_log"}
		handler := AAAConsumerGroupHandler{}
		// 启动kafka消费组模式，消费的逻辑在上面的 ConsumeClaim 这个方法里
		err := group.Consume(ctx, topics, handler)

		if err != nil {
			fmt.Println("消费失败; err : ", err)
			return
		}
	}

}
