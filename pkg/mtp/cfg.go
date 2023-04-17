package mtp

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

type HttpCfg struct {
	Mode     string `yaml:"mode"`
	Port     string `yaml:"port"`
	TLSPort  string `yaml:"tlsPort"`
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

type CoAPServerCfg struct {
	Mode     string `yaml:"mode"`
	Port     string `yaml:"port"`
	DTLSPort string `yaml:"dtlsPort"`
}

type CoAPClientCfg struct {
	Mode           string `yaml:"mode"`
	ServerAddr     string `yaml:"serverAddr"`
	ServerAddrDTLS string `yaml:"serverAddrDTLS"`
}
type CoAPCfg struct {
	Server CoAPServerCfg `yaml:"server"`
	Client CoAPClientCfg `yaml:"client"`
}

type StompCfg struct {
	Mode            string `yaml:"mode"`
	ServerAddr      string `yaml:"serverAddr"`
	ServerAddrTLS   string `yaml:"serverAddrTLS"`
	AgentQueue      string `yaml:"agentQueue"`
	ControllerQueue string `yaml:"controllerQueue"`
	UserName        string `yaml:"userName"`
	Passwd          string `yaml:"passwd"`
	RetryCount      int    `yaml:"retryCount"`
}

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

type MqttCfg struct {
	Mode       string `yaml:"mode"`
	BrokerAddr string `yaml:"brokerAddr"`
	Topic      string `yaml:"topic"`
	UserName   string `yaml:"userName"`
	Passwd     string `yaml:"passwd"`
}

type Cfg struct {
	Http  HttpCfg  `yaml:"http"`
	CoAP  CoAPCfg  `yaml:"coap"`
	Stomp StompCfg `yaml:"stomp"`
	Mqtt  MqttCfg  `yaml:"mqtt"`
	Cache CacheCfg `yaml:"cache"`
	Grpc  GrpcCfg  `yaml:"grpc"`
	Usp   UspCfg   `yaml:"usp"`
	Db    DbCfg    `yaml:"db"`
}

func (m *Mtp) config(confFile string) error {

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
		cfg.Stomp.ServerAddr = stomp
	}

	if mqttBroker, ok := os.LookupEnv("MQTT_ADDR"); ok {
		cfg.Mqtt.BrokerAddr = mqttBroker
	}

	if cache, ok := os.LookupEnv("CACHE_ADDR"); ok {
		cfg.Cache.ServerAddr = cache
	}

	m.Cfg = cfg
	return nil
}
