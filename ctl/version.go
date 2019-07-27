package ctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version Variable Structure
var Version = &cobra.Command{
	Use:   "version",
	Short: "Show current version",
	Long:  "Show Current Application Version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Go WhatsApp CLI Version 1.0")
	},
}
