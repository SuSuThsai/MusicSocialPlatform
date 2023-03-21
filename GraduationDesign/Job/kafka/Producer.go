package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type Producer struct {
	topic    string
	addr     []string
	config   *sarama.Config
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string, topic string) *Producer {
	p := Producer{}
	p.config = sarama.NewConfig()
	//Whether to give a success call
	p.config.Producer.Return.Successes = true
	p.config.Producer.Return.Errors = true
	//set when to reply
	p.config.Producer.RequiredAcks = sarama.WaitForAll
	//user hash-key is better explain random
	p.config.Producer.Partitioner = sarama.NewHashPartitioner
	p.addr = brokers
	p.topic = topic

	producer, err := sarama.NewSyncProducer(p.addr, p.config) //Initialize the client
	if err != nil {
		panic(err.Error())
		return nil
	}
	p.producer = producer
	return &p
}

func (p *Producer) SendMessage(m proto.Message, key string, operationID string) (int32, int64, error) {
	//log.Info(operationID, "SendMessage", "key ", key, m.String(), p.producer)
	kMsg := &sarama.ProducerMessage{}
	kMsg.Topic = p.topic
	kMsg.Key = sarama.StringEncoder(key)
	bMsg, err := proto.Marshal(m)
	if err != nil {
		//log.Error(operationID, "", "proto marshal err = %s", err.Error())
		return -1, -1, err
	}
	if len(bMsg) == 0 {
		//log.Error(operationID, "len(bMsg) == 0 ")
		return 0, 0, errors.New("len(bMsg) == 0 ")
	}
	kMsg.Value = sarama.ByteEncoder(bMsg)
	//log.Info(operationID, "ByteEncoder SendMessage begin", "key ", kMsg, p.producer, "len: ", kMsg.Key.Length(), kMsg.Value.Length())
	if kMsg.Key.Length() == 0 || kMsg.Value.Length() == 0 {
		//log.Error(operationID, "kMsg.Key.Length() == 0 || kMsg.Value.Length() == 0 ", kMsg)
		return -1, -1, errors.New("key or value == 0")
	}
	a, b, c := p.producer.SendMessage(kMsg)
	//log.Info(operationID, "ByteEncoder SendMessage end", "key ", kMsg.Key.Length(), kMsg.Value.Length(), p.producer)
	if c == nil {
		//promePkg.PromeInc(promePkg.SendMsgCounter)
	}
	return a, b, c
}
