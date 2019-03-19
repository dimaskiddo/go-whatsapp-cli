package controller

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	hlp "github.com/dimaskiddo/go-whatsapp-cli/helper"
)

// Login Variable Structure
var Login = &cobra.Command{
	Use:   "login",
	Short: "Login to WhatsApp Web",
	Long:  "Login to WhatsApp Web",
	Run: func(cmd *cobra.Command, args []string) {
		timeout, err := cmd.Flags().GetInt("timeout")
		if err != nil {
			fmt.Println(strings.ToLower(err.Error()))
			return
		}

		file := "./data.gob"

		conn, err := hlp.WAInitConn(timeout)
		if err != nil {
			fmt.Println(strings.ToLower(err.Error()))
			return
		}

		err = hlp.WASessionLogin(conn, file)
		if err != nil {
			fmt.Println(strings.ToLower(err.Error()))
			return
		}

		fmt.Println("successfully login to whatsapp web")
	},
}

func init() {
	Login.Flags().Int("timeout", 10, "Timeout connection in second(s)")
}
