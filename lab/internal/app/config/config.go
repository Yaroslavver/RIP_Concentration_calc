package config

import (
	//"fmt"
	"os"
	"strconv"
	"time"
	
   "github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int
	JWT         JWTConfig
	Redis       RedisConfig
}

type JWTConfig struct {
	SigningMethod string
	ExpiresIn     time.Duration
	TokenSecret   string
}

type RedisConfig struct {
	Host        string
	Port        int
	Password    string
	DB          int
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

func NewConfig() (*Config, error) {
	var err error

   configName := "config"
   _ = godotenv.Load()
   if os.Getenv("CONFIG_NAME") != "" {
      configName = os.Getenv("CONFIG_NAME")
   }

   viper.SetConfigName(configName)
   viper.SetConfigType("toml")
   viper.AddConfigPath("config")
   viper.AddConfigPath(".")
   viper.WatchConfig()

   if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	// Переопределяем секрет из .env (если есть)
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWT.TokenSecret = secret
	}
	// Redis из .env
	if host := os.Getenv("REDIS_HOST"); host != "" {
		cfg.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		p, _ := strconv.Atoi(port)
		cfg.Redis.Port = p
	}
	if pwd := os.Getenv("REDIS_PASSWORD"); pwd != "" {
		cfg.Redis.Password = pwd
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		d, _ := strconv.Atoi(db)
		cfg.Redis.DB = d
	}

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	/*cfg := &Config{} // создаем объект конфига
	err = viper.Unmarshal(cfg) // читаем информацию из файла, 
	// конвертируем и затем кладем в нашу переменную cfg
	if err != nil {
		return nil, err
	}*/

	log.Info("config parsed")

	return cfg, nil
}