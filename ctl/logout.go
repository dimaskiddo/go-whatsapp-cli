package ctl

import (
	"github.com/spf13/cobra"

	"github.com/dimaskiddo/go-whatsapp-cli/hlp"
	"github.com/dimaskiddo/go-whatsapp-cli/hlp/libs"
)

// Logout Variable Structure
var Logout = &cobra.Command{
	Use:   "logout",
	Short: "Logout from WhatsApp Web",
	Long:  "Logout from WhatsApp Web",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

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

		file := "./share/session.gob"

		conn, info, err := libs.WASessionInit(clientVersionMajor, clientVersionMinor, clientVersionBuild, timeout)
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}
		hlp.LogPrintln(hlp.LogLevelInfo, info)

		err = libs.WASessionRestore(conn, file)
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		err = libs.WASessionLogout(conn, file)
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		hlp.LogPrintln(hlp.LogLevelInfo, "successfully logout from whatsapp web")
	},
}

func init() {
	Logout.Flags().Int("client-version-major", 2, "WhatsApp Client major version. Can be override using WHATSAPP_CLIENT_VERSION_MAJOR environment variable")
	Logout.Flags().Int("client-version-minor", 2035, "WhatsApp Client minor version. Can be override using WHATSAPP_CLIENT_VERSION_MINOR environment variable")
	Logout.Flags().Int("client-version-build", 15, "WhatsApp Client build version. Can be override using WHATSAPP_CLIENT_VERSION_BUILD environment variable")

	Logout.Flags().Int("timeout", 5, "Timeout connection in second(s). Can be override using WHATSAPP_TIMEOUT environment variable")
}
