package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	ctl "github.com/dimaskiddo/go-whatsapp-cli/controller"
)

// Root Variable Structure
var r = &cobra.Command{
	Use:   "gowa",
	Short: "Go WhatsApp CLI",
	Long:  "Go WhatsApp CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Init Function
func init() {
	// Initialize Command
	r.AddCommand(ctl.Version)
	r.AddCommand(ctl.Login)
	r.AddCommand(ctl.Daemon)
	r.AddCommand(ctl.Logout)
}

// Main Function
func main() {
	err := r.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
