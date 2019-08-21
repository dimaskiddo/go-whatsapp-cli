package ctl

import (
	"github.com/spf13/cobra"

	"github.com/dimaskiddo/go-whatsapp-cli/hlp"
)

// Logout Variable Structure
var Logout = &cobra.Command{
	Use:   "logout",
	Short: "Logout from WhatsApp Web",
	Long:  "Logout from WhatsApp Web",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		timeout, err := hlp.GetEnvInt("WHATSAPP_TIMEOUT")
		if err != nil {
			timeout, err = cmd.Flags().GetInt("timeout")
			if err != nil {
				hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
			}
		}

		file := "./share/session.gob"

		conn, err := hlp.WASessionInit(timeout)
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		err = hlp.WASessionRestore(conn, file)
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		err = hlp.WASessionLogout(conn, file)
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		hlp.LogPrintln(hlp.LogLevelInfo, "successfully logout from whatsapp web")
	},
}

func init() {
	Logout.Flags().Int("timeout", 10, "Timeout connection in second(s). Can be override using WHATSAPP_TIMEOUT environment variable")
}
