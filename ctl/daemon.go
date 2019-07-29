package ctl

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/dimaskiddo/go-whatsapp-cli/hlp"
)

// Daemon Variable Structure
var Daemon = &cobra.Command{
	Use:   "daemon",
	Short: "Run as daemon service",
	Long:  "Daemon Service for WhatsApp Web",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

		timeout, err := cmd.Flags().GetInt("timeout")
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		reconnect, err := cmd.Flags().GetInt("reconnect")
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		test, err := cmd.Flags().GetBool("test")
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		hlp.CMDList, err = hlp.CMDParse("./share/commands.json")
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		file := "./share/session.gob"

		hlp.LogPrintln(hlp.LogLevelInfo, "starting communication with whatsapp")
		for {
			if hlp.WASessionExist(file) && hlp.WAConn == nil {
				hlp.WAConn, err = hlp.WASessionInit(timeout)
				if err != nil {
					hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
				}

				err = hlp.WASessionRestore(hlp.WAConn, file)
				if err != nil {
					hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
				}

				msisdn := strings.SplitN(hlp.WAConn.Info.Wid, "@", 2)[0]
				masked := msisdn[0:len(msisdn)-3] + "xxx"

				jid := msisdn + "@s.whatsapp.net"
				tag := fmt.Sprintf("@%s", msisdn)

				hlp.LogPrintln(hlp.LogLevelInfo, "logged in to whatsapp as "+masked)

				if test {
					hlp.LogPrintln(hlp.LogLevelInfo, "sending test message to "+masked)

					err = hlp.WAMessageText(hlp.WAConn, jid, "Welcome to Go WhatsApp CLI\nPlease Test Any Handler Here!", 0)
					if err != nil {
						hlp.LogPrintln(hlp.LogLevelError, err.Error())
					}
				}

				<-time.After(time.Second)
				hlp.WAConn.AddHandler(&hlp.WAHandler{
					SessionConn:   hlp.WAConn,
					SessionJid:    jid,
					SessionTag:    tag,
					SessionFile:   file,
					SessionStart:  uint64(time.Now().Unix()),
					ReconnectTime: reconnect,
					IsTest:        test,
				})
			} else if !hlp.WASessionExist(file) && hlp.WAConn != nil {
				_, _ = hlp.WAConn.Disconnect()
				hlp.WAConn = nil

				hlp.LogPrintln(hlp.LogLevelWarn, "disconnected from whatsapp, missing session file")
			} else if !hlp.WASessionExist(file) && hlp.WAConn == nil {
				hlp.LogPrintln(hlp.LogLevelWarn, "trying to login, waiting for session file")
			}

			select {
			case <-sig:
				fmt.Println("")

				if hlp.WAConn != nil {
					_, _ = hlp.WAConn.Disconnect()
				}
				hlp.WAConn = nil

				hlp.LogPrintln(hlp.LogLevelInfo, "terminating process")
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
