package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dimaskiddo/go-whatsapp-cli/pkg/env"
	"github.com/dimaskiddo/go-whatsapp-cli/pkg/log"
	"github.com/dimaskiddo/go-whatsapp-cli/pkg/whatsapp"
)

// Login Variable Structure
var Login = &cobra.Command{
	Use:   "login",
	Short: "Login to WhatsApp Web",
	Long:  "Login to WhatsApp Web",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

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

		file := "./config/stores/session.gob"

		conn, info, err := whatsapp.WASessionInit(clientVersionMajor, clientVersionMinor, clientVersionBuild, timeout)
		if err != nil {
			log.Println(log.LogLevelFatal, err.Error())
		}
		log.Println(log.LogLevelInfo, info)

		err = whatsapp.WASessionLogin(conn, file)
		if err != nil {
			log.Println(log.LogLevelFatal, err.Error())
		}

		log.Println(log.LogLevelInfo, "successfully login to whatsapp web")
	},
}

func init() {
	Login.Flags().Int("client-version-major", 2, "WhatsApp Client major version. Can be override using WHATSAPP_CLIENT_VERSION_MAJOR environment variable")
	Login.Flags().Int("client-version-minor", 2035, "WhatsApp Client minor version. Can be override using WHATSAPP_CLIENT_VERSION_MINOR environment variable")
	Login.Flags().Int("client-version-build", 15, "WhatsApp Client build version. Can be override using WHATSAPP_CLIENT_VERSION_BUILD environment variable")

	Login.Flags().Int("timeout", 5, "Timeout connection in second(s). Can be override using WHATSAPP_TIMEOUT environment variable")
}
