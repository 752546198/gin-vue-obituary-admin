package viper

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"time"
)

// server 基服务本配置结构
type serverStruct struct {
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"readTimeout"`
	WriteTimeout    time.Duration `mapstructure:"writeTimeout"`
	JWTSecret       string        `mapstructure:"jwtSecret"`
	JWTExpire       int           `mapstructure:"jwtExpire"`
	PrefixURL       string        `mapstructure:"PrefixUrl"`
	StaticRootPath  string        `mapstructure:"staticRootPath"`
	UploadImagePath string        `mapstructure:"uploadImagePath"`
	ImageFormats    []string      `mapstructure:"imageFormats"`
	UploadLimit     float64       `mapstructure:"uploadLimit"`
}

// database 数据库配置结构
type databaseStruct struct {
	DBType      string `mapstructure:"dbType"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DBName      string `mapstructure:"dbName"`
	TablePrefix string `mapstructure:"tablePrefix"`
	Debug       bool   `mapstructure:"debug"`
}

// redis 配置结构
type redisStruct struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	Password    string        `mapstructure:"password"`
	DBNum       int           `mapstructure:"db"`
	MaxIdle     int           `mapstructure:"maxIdle"`
	MaxActive   int           `mapstructure:"maxActive"`
	IdleTimeout time.Duration `mapstructure:"idleTimeout"`
}

// pika 配置结构
type pikaStruct struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	Password    string        `mapstructure:"password"`
	DBNum       int           `mapstructure:"db"`
	MaxIdle     int           `mapstructure:"maxIdle"`
	MaxActive   int           `mapstructure:"maxActive"`
	IdleTimeout time.Duration `mapstructure:"idleTimeout"`
}

// logger 配置结构
type loggerStruct struct {
	Level        string `mapstructure:"level"`
	Pretty       bool   `mapstructure:"pretty"`
	Color        bool   `mapstructure:"color"`
	Console      bool   `mapstructure:"console"`
	DebugLogPath string `mapstructure:"debugLogPath"`
	InfoLogPath  string `mapstructure:"infoLogPath"`
	ErrorLogPath string `mapstructure:"errorLogPath"`
}

var (
	ServerConf   = &serverStruct{}
	DatabaseConf = &databaseStruct{}
	RedisConf    = &redisStruct{}
	PikaConf     = &pikaStruct{}
	LoggerConf   = &loggerStruct{}
)

// 初始化 viper，生成服务配置
func Setup() {
	viper.SetConfigType("YAML")
	// 读取 gin 模式，根据模式选择读取对应的配置文件
	data, err := ioutil.ReadFile("config/release.yaml")
	if gin.Mode() == gin.DebugMode {
		data, err = ioutil.ReadFile("config/debug.yaml")
	} else if gin.Mode() == gin.TestMode {
		data, err = ioutil.ReadFile("config/test.yaml")
	}
	if err != nil {
		log.Fatalf("Read the 'config.yaml' fail，msg: %v\n", err)
	}
	// 解析配置文件并进行结构体赋值
	_ = viper.ReadConfig(bytes.NewBuffer(data))
	_ = viper.UnmarshalKey("server", ServerConf)
	_ = viper.UnmarshalKey("database", DatabaseConf)
	_ = viper.UnmarshalKey("redis", RedisConf)
	_ = viper.UnmarshalKey("pika", PikaConf)
	_ = viper.UnmarshalKey("logger", LoggerConf)
}
