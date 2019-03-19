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
		timeout, err := cmd.Flags().GetInt("timeout")
		if err != nil {
			log.Println(strings.ToLower(err.Error()))
			return
		}

		reconnect, err := cmd.Flags().GetInt("reconnect")
		if err != nil {
			log.Println(strings.ToLower(err.Error()))
			return
		}

		test, err := cmd.Flags().GetBool("test")
		if err != nil {
			log.Println(strings.ToLower(err.Error()))
			return
		}

		hlp.WACmd, err = hlp.WAInitCmd("./data.json")
		if err != nil {
			log.Println(strings.ToLower(err.Error()))
			return
		}

		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

		for {
			file := "./data.gob"

			if hlp.WASessionExist(file) && hlp.WAConn == nil {
				hlp.WAConn, err = hlp.WAInitConn(timeout)
				if err != nil {
					log.Println(strings.ToLower(err.Error()))
					return
				}

				err = hlp.WASessionRestore(hlp.WAConn, file)
				if err != nil {
					log.Println(strings.ToLower(err.Error()))
					return
				}

				jid := strings.SplitN(hlp.WAConn.Info.Wid, "@", 2)[0] + "@s.whatsapp.net"
				log.Println("logged in as " + jid)

				if test {
					log.Println("sending text message to " + jid)

					err = hlp.WAMessageText(hlp.WAConn, jid, "Welcome to Go WhatsApp CLI (Test Mode)\nPlease Test Any Handler Here!", 0)
					if err != nil {
						log.Println(strings.ToLower(err.Error()))
					}
				}

				<-time.After(time.Second)
				var wah hlp.WAHandler

				wah.WAConn = hlp.WAConn
				wah.JID = jid
				wah.StartTime = uint64(time.Now().Unix())
				wah.ReconnectTime = reconnect
				wah.TestMode = test

				hlp.WAConn.AddHandler(&wah)
			} else if !hlp.WASessionExist(file) && hlp.WAConn != nil {
				_, _ = hlp.WAConn.Disconnect()
				hlp.WAConn = nil

				log.Println("logged out session not valid")
			} else if !hlp.WASessionExist(file) && hlp.WAConn == nil {
				log.Println("trying to login, waiting for session file ...")
			}

			select {
			case <-sigchan:
				fmt.Println("")
				log.Println("clossing connection")

				_, _ = hlp.WAConn.Disconnect()
				hlp.WAConn = nil

				return
			case <-time.After(5 * time.Second):
			}
		}
	},
}

func init() {
	Daemon.Flags().Int("timeout", 10, "Timeout connection in second(s)")
	Daemon.Flags().Int("reconnect", 30, "Reconnection time when connection closed in second(s)")
	Daemon.Flags().Bool("test", false, "Test mode (only allow from the same id)")
}
