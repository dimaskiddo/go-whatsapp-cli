package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dimaskiddo/go-whatsapp-cli/pkg/env"
	"github.com/dimaskiddo/go-whatsapp-cli/pkg/log"
	"github.com/dimaskiddo/go-whatsapp-cli/pkg/whatsapp"
)

// Logout Variable Structure
var Logout = &cobra.Command{
	Use:   "logout",
	Short: "Logout from WhatsApp Web",
	Long:  "Logout from WhatsApp Web",
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

		err = whatsapp.WASessionRestore(conn, file)
		if err != nil {
			log.Println(log.LogLevelFatal, err.Error())
		}

		err = whatsapp.WASessionLogout(conn, file)
		if err != nil {
			log.Println(log.LogLevelFatal, err.Error())
		}

		log.Println(log.LogLevelInfo, "successfully logout from whatsapp web")
	},
}

func init() {
	Logout.Flags().Int("client-version-major", WhatsAppVerMajor, "WhatsApp Client major version. Can be override using WHATSAPP_CLIENT_VERSION_MAJOR environment variable")
	Logout.Flags().Int("client-version-minor", WhatsAppVerMinor, "WhatsApp Client minor version. Can be override using WHATSAPP_CLIENT_VERSION_MINOR environment variable")
	Logout.Flags().Int("client-version-build", WhatsAppVerBuild, "WhatsApp Client build version. Can be override using WHATSAPP_CLIENT_VERSION_BUILD environment variable")

	Logout.Flags().Int("timeout", 5, "Timeout connection in second(s). Can be override using WHATSAPP_TIMEOUT environment variable")
}
