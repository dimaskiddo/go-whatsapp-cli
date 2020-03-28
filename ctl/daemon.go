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
	"github.com/dimaskiddo/go-whatsapp-cli/hlp/libs"
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

		clientVersionMajor, err := hlp.GetEnvInt("WHATSAPP_CLIENT_VERSION_MAJOR")
		if err != nil {
			clientVersionMajor, err = cmd.Flags().GetInt("client-version-major")
			if err != nil {
				hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
			}
		}

		clientVersionMinor, err := hlp.GetEnvInt("WHATSAPP_CLIENT_VERSION_MINOR")
		if err != nil {
			clientVersionMinor, err = cmd.Flags().GetInt("client-version-minor")
			if err != nil {
				hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
			}
		}

		clientVersionBuild, err := hlp.GetEnvInt("WHATSAPP_CLIENT_VERSION_BUILD")
		if err != nil {
			clientVersionBuild, err = cmd.Flags().GetInt("client-version-build")
			if err != nil {
				hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
			}
		}

		timeout, err := hlp.GetEnvInt("WHATSAPP_TIMEOUT")
		if err != nil {
			timeout, err = cmd.Flags().GetInt("timeout")
			if err != nil {
				hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
			}
		}

		reconnect, err := hlp.GetEnvInt("WHATSAPP_RECONNECT")
		if err != nil {
			reconnect, err = cmd.Flags().GetInt("reconnect")
			if err != nil {
				hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
			}
		}

		test, err := hlp.GetEnvBool("WHATSAPP_TEST")
		if err != nil {
			test, err = cmd.Flags().GetBool("test")
			if err != nil {
				hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
			}
		}

		hlp.CMDList, err = hlp.CMDParse("./share/commands.json")
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		file := "./share/session.gob"

		for {
			if libs.WASessionExist(file) && libs.WAConn == nil {
				var info string

				libs.WAConn, info, err = libs.WASessionInit(clientVersionMajor, clientVersionMinor, clientVersionBuild, timeout)
				if err != nil {
					hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
				}
				hlp.LogPrintln(hlp.LogLevelInfo, info)

				hlp.LogPrintln(hlp.LogLevelInfo, "starting communication with whatsapp")

				err = libs.WASessionRestore(libs.WAConn, file)
				if err != nil {
					hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
				}

				msisdn := strings.SplitN(libs.WAConn.Info.Wid, "@", 2)[0]
				masked := msisdn[0:len(msisdn)-3] + "xxx"

				jid := msisdn + "@s.whatsapp.net"
				tag := fmt.Sprintf("@%s", msisdn)

				hlp.LogPrintln(hlp.LogLevelInfo, "logged in to whatsapp as "+masked)

				if test {
					hlp.LogPrintln(hlp.LogLevelInfo, "sending test message to "+masked)

					err = libs.WAMessageText(libs.WAConn, jid, "Welcome to Go WhatsApp CLI\nPlease Test Any Handler Here!", "", "", 0)
					if err != nil {
						hlp.LogPrintln(hlp.LogLevelError, err.Error())
					}
				}

				<-time.After(time.Second)
				libs.WAConn.AddHandler(&libs.WAHandler{
					SessionConn:   libs.WAConn,
					SessionJid:    jid,
					SessionTag:    tag,
					SessionFile:   file,
					SessionStart:  uint64(time.Now().Unix()),
					ReconnectTime: reconnect,
					IsTest:        test,
				})
			} else if !libs.WASessionExist(file) && libs.WAConn != nil {
				_, _ = libs.WAConn.Disconnect()
				libs.WAConn = nil

				hlp.LogPrintln(hlp.LogLevelWarn, "disconnected from whatsapp, missing session file")
			} else if !libs.WASessionExist(file) && libs.WAConn == nil {
				hlp.LogPrintln(hlp.LogLevelWarn, "trying to login, waiting for session file")
			}

			select {
			case <-sig:
				fmt.Println("")

				if libs.WAConn != nil {
					_, _ = libs.WAConn.Disconnect()
				}
				libs.WAConn = nil

				hlp.LogPrintln(hlp.LogLevelInfo, "terminating process")
				os.Exit(0)
			case <-time.After(5 * time.Second):
			}
		}
	},
}

func init() {
	Daemon.Flags().Int("client-version-major", 0, "WhatsApp Client major version. Can be override using WHATSAPP_CLIENT_VERSION_MAJOR environment variable")
	Daemon.Flags().Int("client-version-minor", 4, "WhatsApp Client minor version. Can be override using WHATSAPP_CLIENT_VERSION_MINOR environment variable")
	Daemon.Flags().Int("client-version-build", 1300, "WhatsApp Client build version. Can be override using WHATSAPP_CLIENT_VERSION_BUILD environment variable")

	Daemon.Flags().Int("timeout", 5, "Timeout connection in second(s). Can be override using WHATSAPP_TIMEOUT environment variable")
	Daemon.Flags().Int("reconnect", 30, "Reconnection time when connection closed in second(s). Can be override using WHATSAPP_RECONNECT environment variable")
	Daemon.Flags().Bool("test", false, "Test mode (only allow from the same ID). Can be override using WHATSAPP_TEST environment variable")
}
