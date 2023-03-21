package kafka

import (
	"github.com/Shopify/sarama"
)

type fcb func(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession)

type PersistentSqlConsumerHandler struct {
	msgHandle                  map[string]fcb
	persistentSqlConsumerGroup *MConsumerGroup
}

//func (pc *PersistentSqlConsumerHandler) Init() {
//	pc.msgHandle = make(map[string]fcb)
//	pc.msgHandle[Config.Conf.Kafka.MsgSql.Topic] = pc.handleChatWs2Mysql
//	pc.persistentSqlConsumerGroup = NewMConsumerGroup(&MConsumerGroupConfig{KafkaVersion: sarama.KafkaVersion{},
//		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{Config.Conf.Kafka.MsgSql.Topic},
//		Config.Conf.Kafka.MsgSql.Brokers, Config.Conf.Kafka.MsgSql.Group)
//}

//	func (pc *PersistentSqlConsumerHandler) handleChatWs2Mysql(cMsg *sarama.ConsumerMessage, msgKey string, _ sarama.ConsumerGroupSession) {
//		msg := cMsg.Value
//		//log.NewInfo("msg come here mysql!!!", "", "msg", string(msg), msgKey)
//		var tag bool
//		msgFromMQ := pbMsg.MsgDataToMQ{}
//		err := proto.Unmarshal(msg, &msgFromMQ)
//		if err != nil {
//			//log.NewError(msgFromMQ.OperationID, "msg_transfer Unmarshal msg err", "msg", string(msg), "err", err.Error())
//			return
//		}
//		//log.Debug(msgFromMQ.OperationID, "proto.Unmarshal MsgDataToMQ", msgFromMQ.String())
//		//Control whether to store history messages (mysql)
//		isPersist := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, constant.IsPersistent)
//		//Only process receiver data
//		if isPersist {
//			switch msgFromMQ.MsgData.SessionType {
//			case constant.SingleChatType, constant.NotificationChatType:
//				if msgKey == msgFromMQ.MsgData.RecvID {
//					tag = true
//				}
//			case constant.GroupChatType:
//				if msgKey == msgFromMQ.MsgData.SendID {
//					tag = true
//				}
//			case constant.SuperGroupChatType:
//				tag = true
//			}
//			if tag {
//				//log.NewInfo(msgFromMQ.OperationID, "msg_transfer msg persisting", string(msg))
//				if err = im_mysql_msg_model.InsertMessageToChatLog(msgFromMQ); err != nil {
//					//log.NewError(msgFromMQ.OperationID, "Message insert failed", "err", err.Error(), "msg", msgFromMQ.String())
//					return
//				}
//			}
//
//		}
//	}
func (PersistentSqlConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (PersistentSqlConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (pc *PersistentSqlConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		//log.NewDebug("", "kafka get info to mysql", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			pc.msgHandle[msg.Topic](msg, string(msg.Key), sess)
		} else {
			//log.Error("", "msg get from kafka but is nil", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
