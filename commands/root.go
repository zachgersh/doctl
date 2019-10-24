package commands

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/digitalocean/doctl"
)

//NewRootCmd makes the base of our doctl command
func NewRootCmd(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctl",
		Short: "doctl is a command line interface for the DigitalOcean API.",
	}

	flagSet := cmd.PersistentFlags()
	flagSet.StringP("config", "c", filepath.Join(configHome(), defaultConfigName), "config file")
	v.BindPFlag("config", flagSet.Lookup("config"))

	flagSet.StringP("api-url", "u", "", "Override default API V2 endpoint")
	v.BindPFlag("api-url", flagSet.Lookup("api-url"))

	flagSet.StringP(doctl.ArgAccessToken, "t", "", "API V2 Access Token")
	v.BindPFlag(doctl.ArgAccessToken, flagSet.Lookup(doctl.ArgAccessToken))

	flagSet.StringP(doctl.ArgOutput, "o", "text", "output format [text|json]")
	v.BindPFlag("output", flagSet.Lookup(doctl.ArgOutput))

	flagSet.StringP(doctl.ArgContext, "", "", "authentication context")
	v.BindPFlag("context", flagSet.Lookup(doctl.ArgContext))

	flagSet.BoolP("trace", "", false, "trace api access")
	flagSet.BoolP(doctl.ArgVerbose, "v", false, "verbose output")

	return cmd
}
