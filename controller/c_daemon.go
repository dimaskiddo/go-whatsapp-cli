package controller

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	hlp "github.com/dimaskiddo/go-whatsapp-cli/helper"
)

// Daemon Variable Structure
var Daemon = &cobra.Command{
	Use:   "daemon",
	Short: "Run as daemon service",
	Long:  "Daemon Service for WhatsApp Web",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		timeout := hlp.GetEnv("WA_TIMEOUT", "int", false)
		if timeout == nil {
			timeout, _ = cmd.Flags().GetInt("timeout")
		}

		reconnect := hlp.GetEnv("WA_RECONNECT", "int", false)
		if reconnect == nil {
			reconnect, _ = cmd.Flags().GetInt("reconnect")
		}

		test := hlp.GetEnv("WA_TEST", "bool", false)
		if test == nil {
			test, _ = cmd.Flags().GetBool("test")
		}

		hlp.CMDList, err = hlp.CMDParse("./data.json")
		if err != nil {
			log.Fatalln(strings.ToLower(err.Error()))
		}

		file := "./data.gob"

		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

		log.Println("daemon: starting communication with whatsapp")

		for {
			if hlp.WASessionExist(file) && hlp.WAConn == nil {
				hlp.WAConn, err = hlp.WASessionInit(timeout.(int))
				if err != nil {
					log.Println(strings.ToLower(err.Error()))
					return
				}

				err = hlp.WASessionRestore(hlp.WAConn, file)
				if err != nil {
					log.Println(strings.ToLower(err.Error()))
					return
				}

				msisdn := strings.SplitN(hlp.WAConn.Info.Wid, "@", 2)[0]
				masked := msisdn[0:len(msisdn)-3] + "xxx"
				jid := msisdn + "@s.whatsapp.net"

				log.Println("daemon: logged in to whatsapp as " + masked)

				if test.(bool) {
					log.Println("daemon: sending test message to " + masked)

					err = hlp.WAMessageText(hlp.WAConn, jid, "Welcome to Go WhatsApp CLI\nPlease Test Any Handler Here!", 0)
					if err != nil {
						log.Println(strings.ToLower(err.Error()))
					}
				}

				<-time.After(time.Second)
				hlp.WAConn.AddHandler(&hlp.WAHandler{
					SessionConn:   hlp.WAConn,
					SessionJid:    jid,
					SessionFile:   file,
					SessionStart:  uint64(time.Now().Unix()),
					ReconnectTime: reconnect.(int),
					IsTest:        test.(bool),
				})
			} else if !hlp.WASessionExist(file) && hlp.WAConn != nil {
				_, _ = hlp.WAConn.Disconnect()
				hlp.WAConn = nil

				log.Println("daemon: disconnected from whatsapp, missing session file")
			} else if !hlp.WASessionExist(file) && hlp.WAConn == nil {
				log.Println("daemon: trying to login, waiting for session file...")
			}

			select {
			case <-sigchan:
				fmt.Println("")

				if hlp.WAConn != nil {
					_, _ = hlp.WAConn.Disconnect()
				}
				hlp.WAConn = nil

				log.Println("daemon: terminating process")
				os.Exit(0)
			case <-time.After(5 * time.Second):
			}
		}
	},
}

func init() {
	Daemon.Flags().Int("timeout", 10, "Timeout connection in second(s), can be override using environment variable WA_TIMEOUT")
	Daemon.Flags().Int("reconnect", 30, "Reconnection time when connection closed in second(s), can be override using environment variable WA_RECONNECT")
	Daemon.Flags().Bool("test", false, "Test mode (only allow from the same id), can be override using environment variable WA_TEST")
}
