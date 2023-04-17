package rest

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/n4-networks/usp/pkg/db"
	"github.com/n4-networks/usp/pkg/pb/mtpgrpc"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type Cfg struct {
	httpPort    string
	isTlsOn     bool
	mtpAddr     string
	dbAddr      string
	dbUserName  string
	dbPasswd    string
	connTimeout time.Duration
	logSetting  string
}

type mtpHandle struct {
	grpcIntf mtpgrpc.MtpGrpcClient
	grpcConn *grpc.ClientConn
	txMsgCnt uint64
}

func (m *mtpHandle) incTxMsgCnt() uint64 {
	m.txMsgCnt++
	return m.txMsgCnt
}

type dbHandle struct {
	client  *mongo.Client
	uspIntf *db.UspDb
}

type Rest struct {
	mtp    mtpHandle
	db     dbHandle
	cfg    Cfg
	router *mux.Router
}

func (re *Rest) Init() error {

	log.Println("Running re Server version:", getVer())

	log.Println("Reading config parameters...")
	if err := re.config(); err != nil {
		log.Println("Could not configure re Server, err:", err)
		return err
	}

	// Initialize logging
	log.Println("Initializing logging module...")
	if err := re.loggingInit(); err != nil {
		log.Println("Logging settings could not be applied")
	}
	// Connect o Db
	log.Println("Connecting to DB server @", re.cfg.dbAddr)
	if err := re.connectDb(); err != nil {
		log.Println("Error in connecting to DB:", err)
	}

	// Connect to MTP server
	log.Println("Connecting to MTP @", re.cfg.mtpAddr)
	if err := re.connectMtp(); err != nil {
		log.Println("Error in connecting to Mtp:", err)
	} else {
		log.Println("Connection to MTP...Success")
	}

	// Initialize Router
	if err := re.initRouter(); err != nil {
		log.Println("Error in initializing Router:", err)
	} else {
		log.Println("Initializing Router...Success")
	}
	log.Println("N4-RESt Server has been initialized")
	return nil
}

func (re *Rest) config() error {

	if httpPort, ok := os.LookupEnv("HTTP_PORT"); ok {
		re.cfg.httpPort = httpPort
	} else {
		re.cfg.httpPort = "8080"
	}

	isTlsOn, ok := os.LookupEnv("HTTP_TLS")
	if ok && isTlsOn == "1" {
		re.cfg.isTlsOn = true
	} else {
		re.cfg.isTlsOn = false
	}

	if dbAddr, ok := os.LookupEnv("DB_ADDR"); ok {
		re.cfg.dbAddr = dbAddr
	} else {
		re.cfg.dbAddr = ":27017"
	}

	if dbUserName, ok := os.LookupEnv("DB_USER"); ok {
		re.cfg.dbUserName = dbUserName
	} else {
		log.Println("DB_USER is not set")
		return errors.New("DB_USER not set")
	}

	if dbPasswd, ok := os.LookupEnv("DB_PASSWD"); ok {
		re.cfg.dbPasswd = dbPasswd
	} else {
		log.Println("DB_PASSWD is not set")
		return errors.New("DB_PASSWD not set")
	}

	if mtpGrpcAddr, ok := os.LookupEnv("MTP_GRPC_ADDR"); ok {
		re.cfg.mtpAddr = mtpGrpcAddr
	} else {
		re.cfg.mtpAddr = ":9001"
	}

	re.cfg.connTimeout = 10 * time.Second

	if logging, ok := os.LookupEnv("LOGGING"); ok {
		re.cfg.logSetting = logging
	} else {
		re.cfg.logSetting = "none"
	}
	return nil
}

func (re *Rest) loggingInit() error {
	log.SetPrefix("N4: ")
	switch re.cfg.logSetting {
	case "short":
		log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	case "long":
		log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	case "all":
		log.Println("Setting log for all")
		log.SetFlags(log.Lshortfile | log.Llongfile | log.Ldate | log.Ltime)
	default:
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	return nil
}
