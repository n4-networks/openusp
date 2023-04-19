package cli

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/go-stomp/stomp"
)

type cfgHandler struct {
	restSrvAddr string
	stompAddr   string
	histFile    string
	connTimeout time.Duration
	logSetting  string
}
type restHandler struct {
	client *http.Client
}

type shHandler struct {
	shell    *ishell.Shell
	histFile string
	cmds     map[string]*ishell.Cmd
}

type stompHandler struct {
	client *stomp.Conn
}

type Cli struct {
	agent      agentInfo
	stomp      stompHandler
	sh         shHandler
	cfg        cfgHandler
	rest       restHandler
	lastCmdErr error
}

func (cli *Cli) GetLastCmdErr() error {
	return cli.lastCmdErr
}

func (cli *Cli) ClearLastCmdErr() {
	cli.lastCmdErr = nil
}

func (cli *Cli) Init() error {

	if err := cli.config(); err != nil {
		log.Println("Could not configure CLI, err:", err)
		return err
	}

	// Initialize logging
	if err := cli.loggingInit(); err != nil {
		log.Println("Logging settings could not be applied")
	}

	// Initialization rest client
	if err := cli.restInit(); err != nil {
		log.Println("Could not initialize rest client:", err)
		return err
	}

	// Initialization of Agent Parameters
	if err := cli.initCliWithAgentParams(); err != nil {
		log.Println("Could not set agent information:", err)
	}
	log.Println("CLI version:", getVer())

	// Initialize shell
	cli.sh.shell = ishell.New()

	// Set default Prompt
	cli.sh.shell.SetPrompt("OpenUsp-Cli>> ")
	cli.sh.histFile = "history"
	cli.sh.shell.SetHistoryPath(cli.sh.histFile)

	// Initialize shell Cmds
	cli.sh.cmds = make(map[string]*ishell.Cmd)

	// Initialize shell Cmds
	cli.sh.cmds = make(map[string]*ishell.Cmd)

	// Register verb cmds
	cli.registerVerbs()

	// MTP and DB
	cli.registerNounsMtp()
	cli.registerNounsDb()
	cli.registerNounsStomp()

	// CLI related
	cli.registerNounsHistory()
	cli.registerNounsLogging()
	cli.registerNounsVersion()

	// Agent
	cli.registerNounsAgent()

	// Device Model
	cli.registerNounsDevice()
	cli.registerNounsDevInfo()
	cli.registerNounsBridging()
	cli.registerNounsDhcpv4()
	cli.registerNounsEth()
	cli.registerNounsIp()
	cli.registerNounsNat()
	cli.registerNounsWiFi()
	cli.registerNounsTime()
	cli.registerNounsNw()

	// Basic low level
	cli.registerNounsDatamodel()
	cli.registerNounsCommand()
	cli.registerNounsParam()
	cli.registerNounsInstance()

	return nil
}

func (cli *Cli) Run() {
	cli.sh.shell.Println("**************************************************************")
	cli.sh.shell.Println("                          OpenUsp Cli")
	cli.sh.shell.Println("**************************************************************")
	cli.sh.shell.Run()
}

func (cli *Cli) config() error {

	if restSrvAddr, ok := os.LookupEnv("REST_SRV_ADDR"); ok {
		cli.cfg.restSrvAddr = restSrvAddr
	} else {
		cli.cfg.restSrvAddr = "http://localhost:8080"
	}

	cli.cfg.connTimeout = 10 * time.Second
	cli.cfg.histFile = "history"

	if logging, ok := os.LookupEnv("CLI_LOGGING"); ok {
		cli.cfg.logSetting = logging
	} else {
		cli.cfg.logSetting = "none"
	}
	return nil
}

func (cli *Cli) ProcessCmd(args string) error {
	log.Println("Running cli command in non interactive mode")
	log.Println("Processing cmd:", args)
	tok := strings.Split(args, " ")
	return cli.sh.shell.Process(tok...)
}

func (cli *Cli) SetOut(writer io.Writer) error {
	cli.sh.shell.SetOut(writer)
	return nil
}
