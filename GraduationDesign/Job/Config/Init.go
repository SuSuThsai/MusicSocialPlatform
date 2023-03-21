package Config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var Conf *Config
var err error

type Config struct {
	Etcd  *ETCD
	Kafka *Kafka
	Env   *Env
}

type Kafka struct {
	MsgChat  *MsgChat
	MsgAudio *MsgAudio
	MsgSql   *MsgSql
}

type MsgChat struct {
	Topic   string
	Group   string
	Brokers []string
}

type MsgAudio struct {
	Topic   string
	Group   string
	Brokers []string
}

type MsgSql struct {
	Topic   string
	Group   string
	Brokers []string
}

type ETCD struct {
	EtcdAddr []string
}

type Env struct {
	Env  string
	Host string
}

func Init() {
	path, _ := os.Getwd()
	configPath := path + "/Job/Config/config.toml"
	Conf = Default()
	_, err = toml.DecodeFile(configPath, &Conf)
	if err != nil {
		log.Panic("读取初始化文件失败！", err.Error())
	}
	return
}

func Default() *Config {
	return &Config{}
}
