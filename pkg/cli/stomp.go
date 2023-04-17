package cli

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/go-stomp/stomp"
)

const (
	DEFAULT_MAX_STOMP_MSG_TO_REMOVE = 1
)

func (cli *Cli) registerNounsStomp() {
	cmds := []noun{
		{"reconnect", "stomp", connectStompHelp, cli.connectStompCmd},
		{"show", "stomp", showStompHelp, cli.showStomp},
		{"remove", "stomp", removeStompHelp, cli.removeStomp},
		{"remove.stomp", "conn", removeStompConnHelp, cli.removeStompConn},
		{"remove.stomp", "msgs", removeStompMsgsHelp, cli.removeStompMsgs},
	}
	cli.registerNouns(cmds)
}

const removeStompHelp = "remove stomp conn|queue|..."

func (cli *Cli) removeStomp(c *ishell.Context) {
	c.Printf(removeStompHelp)
}

const removeStompMsgsHelp = "remove stomp msgs <number-of-msgs>"

func (cli *Cli) removeStompMsgs(c *ishell.Context) {
	if cli.stomp.client == nil {
		c.Println("Not connected to STOMP server, use connect command to connect")
		return
	}
	if !cli.agent.isSet.epId {
		c.Println("Agent id has not been set, use set agentid command to configure")
		return
	}
	queueName := cli.agent.epId
	sub, err := cli.stomp.client.Subscribe(queueName, stomp.AckAuto)
	if err != nil {
		c.Println("Error: could not subscribe to", queueName, err.Error())
		return
	}
	var numOfMsg int64
	if len(c.Args) > 0 {
		numOfMsg, _ = strconv.ParseInt(c.Args[0], 0, 64)
	} else {
		numOfMsg = DEFAULT_MAX_STOMP_MSG_TO_REMOVE
	}
	log.Println("Number of msgs to be removed:", numOfMsg)

	msgCount := 0

	for i := 1; i <= int(numOfMsg); i++ {
		select {
		case <-time.After(1 * time.Second):
			break
		case msg := <-sub.C:
			msg = msg // this is to avoid "variable declared but not used" error
			msgCount++
		}
	}
	c.Println("Number of msg removed:", msgCount)
	if err := sub.Unsubscribe(); err != nil {
		log.Println("Could not ubsubscribe to the agent queue")
	}
}

const removeStompConnHelp = "remove stomp conn"

func (cli *Cli) removeStompConn(c *ishell.Context) {
	if cli.stomp.client != nil {
		c.Println("Disconncting from server:", cli.cfg.stompAddr)
		if err := cli.stomp.client.Disconnect(); err != nil {
			c.Println("Error: Could not disconnect, err:", err)
			return
		}
		c.Println("..Success")
		cli.stomp.client = nil
	}
}

const showStompHelp = "show stomp"

func (cli *Cli) showStomp(c *ishell.Context) {
	c.Printf("%-25s\n", "STOMP Connection Status")
	if cli.stomp.client != nil {
		c.Printf(" %-24s : %-12s\n", "Address", cli.cfg.stompAddr)
		c.Printf(" %-24s : %-12s\n", "Connection", "Connected")
	} else {
		c.Printf(" %-24s : %-12s\n", "Connection", "Not Connected")
	}
	c.Println("-------------------------------------------------")
}

const connectStompHelp = "connect stomp <addr:port>"

func (cli *Cli) connectStompCmd(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println("Wrong input.", connectStompHelp)
		return
	}
	if cli.stomp.client != nil {
		cli.stomp.client.Disconnect()
	}
	dbAddr := c.Args[0]
	if err := cli.connectStomp(dbAddr, cli.cfg.connTimeout); err != nil {
		c.Println("Error:", err)
		return
	}
	c.Println("Success")
}

func (cli *Cli) connectStomp(addr string, timeout time.Duration) error {
	log.Println("Connecting to STOMP Broker @", addr)
	netConn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
		stomp.ConnOpt.Login("guest", "guest"),
		stomp.ConnOpt.Host("/"),
	}

	stompConn, err := stomp.Connect(netConn, options...)
	if err != nil {
		return err
	}
	cli.stomp.client = stompConn
	log.Println("Connection to STOMP ...SUCCESS")

	return nil
}
