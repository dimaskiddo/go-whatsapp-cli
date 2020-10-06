package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/dimaskiddo/go-whatsapp-cli/pkg/env"
	"github.com/dimaskiddo/go-whatsapp-cli/pkg/log"
	"github.com/dimaskiddo/go-whatsapp-cli/pkg/parser"
	"github.com/dimaskiddo/go-whatsapp-cli/pkg/whatsapp"
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

		clientVersionMajor, err := env.GetEnvInt("WHATSAPP_CLIENT_VERSION_MAJOR")
		if err != nil {
			clientVersionMajor, err = cmd.Flags().GetInt("client-version-major")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		clientVersionMinor, err := env.GetEnvInt("WHATSAPP_CLIENT_VERSION_MINOR")
		if err != nil {
			clientVersionMinor, err = cmd.Flags().GetInt("client-version-minor")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		clientVersionBuild, err := env.GetEnvInt("WHATSAPP_CLIENT_VERSION_BUILD")
		if err != nil {
			clientVersionBuild, err = cmd.Flags().GetInt("client-version-build")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		timeout, err := env.GetEnvInt("WHATSAPP_TIMEOUT")
		if err != nil {
			timeout, err = cmd.Flags().GetInt("timeout")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		reconnect, err := env.GetEnvInt("WHATSAPP_RECONNECT")
		if err != nil {
			reconnect, err = cmd.Flags().GetInt("reconnect")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		test, err := env.GetEnvBool("WHATSAPP_TEST")
		if err != nil {
			test, err = cmd.Flags().GetBool("test")
			if err != nil {
				log.Println(log.LogLevelFatal, err.Error())
			}
		}

		parser.JSONList, err = parser.JSONParse("./config/json/cmds.json")
		if err != nil {
			log.Println(log.LogLevelFatal, err.Error())
		}

		file := "./config/stores/session.gob"

		for {
			if whatsapp.WASessionExist(file) && whatsapp.WAConn == nil {
				var info string

				whatsapp.WAConn, info, err = whatsapp.WASessionInit(clientVersionMajor, clientVersionMinor, clientVersionBuild, timeout)
				if err != nil {
					log.Println(log.LogLevelFatal, err.Error())
				}
				log.Println(log.LogLevelInfo, info)

				log.Println(log.LogLevelInfo, "starting communication with whatsapp")

				err = whatsapp.WASessionRestore(whatsapp.WAConn, file)
				if err != nil {
					log.Println(log.LogLevelFatal, err.Error())
				}

				msisdn := strings.SplitN(whatsapp.WAConn.Info.Wid, "@", 2)[0]
				masked := msisdn[0:len(msisdn)-3] + "xxx"

				jid := msisdn + "@s.whatsapp.net"
				tag := fmt.Sprintf("@%s", msisdn)

				log.Println(log.LogLevelInfo, "logged in to whatsapp as "+masked)

				if test {
					log.Println(log.LogLevelInfo, "sending test message to "+masked)

					err = whatsapp.WAMessageText(whatsapp.WAConn, jid, "Welcome to Go WhatsApp CLI\nPlease Test Any Handler Here!", "", "")
					if err != nil {
						log.Println(log.LogLevelError, err.Error())
					}
				}

				<-time.After(time.Second)
				whatsapp.WAConn.AddHandler(&whatsapp.WAHandler{
					SessionConn:   whatsapp.WAConn,
					SessionJid:    jid,
					SessionTag:    tag,
					SessionFile:   file,
					SessionStart:  uint64(time.Now().Unix()),
					ReconnectTime: reconnect,
					IsTest:        test,
				})
			} else if !whatsapp.WASessionExist(file) && whatsapp.WAConn != nil {
				_, _ = whatsapp.WAConn.Disconnect()
				whatsapp.WAConn = nil

				log.Println(log.LogLevelWarn, "disconnected from whatsapp, missing session file")
			} else if !whatsapp.WASessionExist(file) && whatsapp.WAConn == nil {
				log.Println(log.LogLevelWarn, "trying to login, waiting for session file")
			}

			select {
			case <-sig:
				fmt.Println("")

				if whatsapp.WAConn != nil {
					_, _ = whatsapp.WAConn.Disconnect()
				}
				whatsapp.WAConn = nil

				log.Println(log.LogLevelInfo, "terminating process")
				os.Exit(0)
			case <-time.After(5 * time.Second):
			}
		}
	},
}

func init() {
	Daemon.Flags().Int("client-version-major", 2, "WhatsApp Client major version. Can be override using WHATSAPP_CLIENT_VERSION_MAJOR environment variable")
	Daemon.Flags().Int("client-version-minor", 2035, "WhatsApp Client minor version. Can be override using WHATSAPP_CLIENT_VERSION_MINOR environment variable")
	Daemon.Flags().Int("client-version-build", 15, "WhatsApp Client build version. Can be override using WHATSAPP_CLIENT_VERSION_BUILD environment variable")

	Daemon.Flags().Int("timeout", 5, "Timeout connection in second(s). Can be override using WHATSAPP_TIMEOUT environment variable")
	Daemon.Flags().Int("reconnect", 30, "Reconnection time when connection closed in second(s). Can be override using WHATSAPP_RECONNECT environment variable")
	Daemon.Flags().Bool("test", false, "Test mode (only allow from the same ID). Can be override using WHATSAPP_TEST environment variable")
}
