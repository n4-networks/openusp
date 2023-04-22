// Copyright 2023 N4-Networks.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/abiosoft/ishell"
)

func (cli *Cli) registerNounsLogging() {
	logCmds := []noun{
		{"show", "logging", showLoggingHelp, cli.showLogging},
		{"set", "logging", setLoggingHelp, cli.setLogging},
	}
	cli.registerNouns(logCmds)
}

func (cli *Cli) loggingInit() error {
	log.SetPrefix("OpenUsp: ")
	switch cli.cfg.logSetting {
	case "short":
		log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	case "long":
		log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	case "all":
		log.SetFlags(log.Lshortfile | log.Llongfile | log.Ldate | log.Ltime)
	default:
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	return nil
}

const showLoggingHelp = "show logging"

func (cli *Cli) showLogging(c *ishell.Context) {
	//c.Printf("%-25s : %-12s\n", "Cli logging Status", inst.path)
	f := log.Flags()
	if f&log.Ldate == 0 {
		c.Printf("  %-24s : %s\n", "Date Flag", "OFF")
	} else {
		c.Printf("  %-24s : %s\n", "Date Flag", "ON")
	}
	if f&log.Ltime == 0 {
		c.Printf("  %-24s : %s\n", "Time Flag", "OFF")
	} else {
		c.Printf("  %-24s : %s\n", "Time Flag", "ON")
	}
	if f&log.Lshortfile == 0 {
		c.Printf("  %-24s : %s\n", "Shortfile Flag", "OFF")
	} else {
		c.Printf("  %-24s : %s\n", "Shortfile Flag", "ON")
	}
	if f&log.Llongfile == 0 {
		c.Printf("  %-24s : %s\n", "Longfile Flag", "OFF")
	} else {
		c.Printf("  %-24s : %s\n", "Longfile Flag", "ON")
	}
	c.Println("-------------------------------------------------")

}

const setLoggingHelp = "set logging <on|off> <all|date|time|long|short>"

func (cli *Cli) setLogging(c *ishell.Context) {
	if len(c.Args) < 2 {
		c.Println(setLoggingHelp)
		return
	}
	cmd := c.Args[0]
	flag := c.Args[1]
	switch flag {
	case "all":
		if cmd == "off" {
			log.SetFlags(0)
			log.SetOutput(ioutil.Discard)
			c.Println("Switched off all logging flags")
			log.Println("If you are seeing this then something is broken...")
		} else if cmd == "on" {
			log.SetFlags(log.Lshortfile | log.Llongfile | log.Ldate | log.Ltime)
			log.SetOutput(os.Stdout)
			c.Println("Switched on all logging flags")
			log.Println("This message from logging engine, its has been switched on now")
		} else {
			c.Println("Invalid cmd. Syntax:", setLoggingHelp)
		}
	case "date":
		if cmd == "off" {
			f := log.Flags() &^ log.Ldate
			log.SetFlags(f)
			c.Println("Logging date flag has been switched off")
			log.Println("This msg through log should not have date")
		} else if cmd == "on" {
			f := log.Flags() | log.Ldate
			log.SetFlags(f)
			c.Println("Logging date flag has been switched on")
			log.Println("This msg through log should have date")
		} else {
			c.Println("Invalid cmd. Syntax:", setLoggingHelp)
		}
	case "time":
		if cmd == "off" {
			f := log.Flags() &^ log.Ltime
			log.SetFlags(f)
			c.Println("Logging time flag has been switched off")
			log.Println("This msg through log should not have time")
		} else if cmd == "on" {
			f := log.Flags() | log.Ltime
			log.SetFlags(f)
			c.Println("Logging time flag has been switched on")
			log.Println("This msg through log should have time")
		} else {
			c.Println("Invalid cmd. Syntax:", setLoggingHelp)
		}
	case "long":
		if cmd == "off" {
			f := log.Flags() &^ log.Llongfile
			log.SetFlags(f)
			c.Println("Logging long flag has been switched off")
			log.Println("This msg through log should not have long format")
		} else if cmd == "on" {
			f := log.Flags() &^ log.Lshortfile
			f = f | log.Llongfile
			log.SetFlags(f)
			c.Println("Logging long flag has been switched on")
			log.Println("This msg through log should have long format")
		} else {
			c.Println("Invalid cmd. Syntax:", setLoggingHelp)
		}
	case "short":
		if cmd == "off" {
			f := log.Flags() &^ log.Lshortfile
			log.SetFlags(f)
			c.Println("Logging short flag has been switched off")
			log.Println("This msg through log should not have short format")
		} else if cmd == "on" {
			f := log.Flags() | log.Lshortfile
			log.SetFlags(f)
			c.Println("Logging short flag has been switched on")
			log.Println("This msg through log should have short format")
		} else {
			c.Println("Invalid cmd. Syntax:", setLoggingHelp)
		}
	default:
		c.Println("Invalid cmd. Syntax:", setLoggingHelp)
	}
}
