package ctl

import (
	"github.com/spf13/cobra"

	"github.com/dimaskiddo/go-whatsapp-cli/hlp"
)

// Login Variable Structure
var Login = &cobra.Command{
	Use:   "login",
	Short: "Login to WhatsApp Web",
	Long:  "Login to WhatsApp Web",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		timeout, err := cmd.Flags().GetInt("timeout")
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		file := "./share/session.gob"

		conn, err := hlp.WASessionInit(timeout)
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		err = hlp.WASessionLogin(conn, file)
		if err != nil {
			hlp.LogPrintln(hlp.LogLevelFatal, err.Error())
		}

		hlp.LogPrintln(hlp.LogLevelInfo, "successfully login to whatsapp web")
	},
}

func init() {
	Login.Flags().Int("timeout", 10, "Timeout connection in second(s), can be override using environment variable WA_TIMEOUT")
}
