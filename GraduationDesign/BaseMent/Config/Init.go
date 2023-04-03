package Config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/sftp"
	"github.com/redis/go-redis/v9"
	"github.com/tencentyun/cos-go-sdk-v5"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB
var DBR *redis.Client
var DBR2 *redis.Client
var TCCos *cos.Client
var Conf *Config
var SftpClient *sftp.Client
var FileRpc *grpc.ClientConn
var FileRpc2 *grpc.ClientConn
var Identify string
var GlobalUserCommandListen map[string]map[uint]bool

//var NowUser *middleware.Claims

type Config struct {
	PostgreSQL *PostgreSQLInfo
	Redis      *RedisInfo
	JwtK       *JWT
	FileGrpc   *FileGrpc
	SetModel   *SetModel
	Cache      *Cache
	Rank       *Rank
	TencentCos *TencentCos
	SFtp       *SFtp
}

type PostgreSQLInfo struct {
	TypeP      string `json:"type_p"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	DbUser     string `json:"db_user" gorm:"DEFAULT:root"`
	DbPassword string `json:"db_password"`
	DbName     string `json:"db_name"`
}

type RedisInfo struct {
	HostR       string `json:"host_r"`
	PortR       int    `json:"port_r"`
	DbUserR     string `json:"dbuser_r" gorm:"DEFAULT:root"`
	DbPassWordR string `json:"db_pass_word_r"`
	DBModel     int    `json:"db_model"`
}

type SetModel struct {
	AppNode  string
	AppNode2 string
}

type JWT struct {
	JwtKey string
}

type FileGrpc struct {
	Server []string
}

type Cache struct {
	CacheMusicListenCount    string
	CacheMusicLikeCount      string
	CacheMusicSortCount      string
	CacheMusicListLikeCount  string
	CacheMusicListSortCount  string
	CacheArticleReadCount    string
	CacheArticleCommentCount string
	CacheArticleLikeCount    string
	CacheArticleForward      string
	CacheCommentLikeCount    string
	CacheMusicDayListen      string
	CacheMusicListDayListen  string
}

type Rank struct {
	CacheRankMusicListYear  string
	CacheRankMusicListMonth string
	CacheRankMusicListWeek  string
	CacheRankMusicListDay   string
	CacheRankMusicYear      string
	CacheRankMusicMonth     string
	CacheRankMusicWeek      string
	CacheRankMusicDay       string
}

type TencentCos struct {
	Url       string
	SecretId  string
	SecretKey string
}

type SFtp struct {
	User     string
	Password string
	Host     string
	Port     string
}

func InitGlobalUserCommandListen() {
	GlobalUserCommandListen = make(map[string]map[uint]bool)
}

func InitsSQL() {
	InitsConfig()
	InitGlobalUserCommandListen()
	InitsPSQL()
	InitRedis()
	InitTencentCos()
	InitFtp()
	//InitFileRpc()
}

func InitsConfig() {
	path, _ := os.Getwd()
	configPath := path + "/BaseMent/Config/config.toml"
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
