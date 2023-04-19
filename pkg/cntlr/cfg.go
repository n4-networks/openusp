package cntlr

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	_                  = iota
	SERVER_MODE_NORMAL = iota
	SERVER_MODE_TLS
	SERVER_MODE_NORMAL_AND_TLS
)

type GrpcCfg struct {
	Port string `yaml:"port"`
}

type CacheCfg struct {
	ServerAddr string `yaml:"serverAddr"`
}

type DbCfg struct {
	ServerAddr string `yaml:"serverAddr"`
	Name       string `yaml:"name"`
	UserName   string `yaml:"userName"`
	Passwd     string `yaml:"passwd"`
}

type UspCfg struct {
	EndpointId   string `yaml:"endpointId"`
	ProtoVersion string `yaml:"protoVersion"`
}

type Cfg struct {
	Mtp   MtpCfg   `yaml:"mtp"`
	Cache CacheCfg `yaml:"cache"`
	Grpc  GrpcCfg  `yaml:"grpc"`
	Usp   UspCfg   `yaml:"usp"`
	Db    DbCfg    `yaml:"db"`
}

func (c *Cntlr) config(confFile string) error {

	f, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Fatal("Could not open the conf file, going with default, error: ", err)
		return err
	}
	cfg := &Cfg{}
	if err := yaml.Unmarshal(f, cfg); err != nil {
		log.Fatal("Could not read conf parameters, error: ", err)
		return err
	}
	/*
		log.Printf("HTTP Configuration\n")
		log.Printf("%+v\n", cfg.Http)
		log.Printf("CoAP Configuration\n")
		log.Printf("%+v\n", cfg.CoAP)
		log.Printf("STOMP Configuration\n")
		log.Printf("%+v\n", C.Cfg.Usp.EndpointId)
	*/
	if db, ok := os.LookupEnv("DB_ADDR"); ok {
		cfg.Db.ServerAddr = db
	}

	if dbUser, ok := os.LookupEnv("DB_USER"); ok {
		cfg.Db.UserName = dbUser
	} else {
		log.Println("DB User name is not set")
		return errors.New("DB_USER is not set")
	}

	if dbPasswd, ok := os.LookupEnv("DB_PASSWD"); ok {
		cfg.Db.Passwd = dbPasswd
	} else {
		log.Println("DB User password is not set")
		return errors.New("DB_PASSWD is not set")
	}
	log.Printf("%+v\n", cfg.Db)

	if stomp, ok := os.LookupEnv("STOMP_ADDR"); ok {
		cfg.Mtp.Stomp.ServerAddr = stomp
	}

	if mqttBroker, ok := os.LookupEnv("MQTT_ADDR"); ok {
		cfg.Mtp.Mqtt.BrokerAddr = mqttBroker
	}

	if cache, ok := os.LookupEnv("CACHE_ADDR"); ok {
		cfg.Cache.ServerAddr = cache
	}

	c.Cfg = cfg
	return nil
}
