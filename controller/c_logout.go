package controller

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	hlp "github.com/dimaskiddo/go-whatsapp-cli/helper"
)

// Logout Variable Structure
var Logout = &cobra.Command{
	Use:   "logout",
	Short: "Logout from WhatsApp Web",
	Long:  "Logout from WhatsApp Web",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		timeout := hlp.GetEnv("WA_TIMEOUT", "int", false)
		if timeout == nil {
			timeout, _ = cmd.Flags().GetInt("timeout")
		}

		file := "./data.gob"

		conn, err := hlp.WASessionInit(timeout.(int))
		if err != nil {
			fmt.Println(strings.ToLower(err.Error()))
			return
		}

		err = hlp.WASessionRestore(conn, file)
		if err != nil {
			fmt.Println(strings.ToLower(err.Error()))
			return
		}

		err = hlp.WASessionLogout(conn, file)
		if err != nil {
			fmt.Println(strings.ToLower(err.Error()))
			return
		}

		fmt.Println("successfully logout from whatsapp web")
	},
}

func init() {
	Logout.Flags().Int("timeout", 10, "Timeout connection in second(s), can be override using environment variable WA_TIMEOUT")
}
