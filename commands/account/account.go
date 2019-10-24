package account

import (
	"github.com/digitalocean/doctl/do"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type accountGet struct{}

type accountRateLimit struct{}

//NewAccountCmd ghasdfhjasdhfj
func NewAccountCmd(v *viper.Viper, client do.AccountService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "account commands",
		Long:  "account is used to access account commands",
	}

	accountGet := &accountGet{}
	cmd.AddCommand(&cobra.Command{
		Use:     "get",
		Short:   "get account",
		Aliases: []string{"get"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return accountGet.run()
		},
	})

	accountRateLimit := &accountRateLimit{}
	cmd.AddCommand(&cobra.Command{
		Use:     "ratelimit",
		Short:   "get API rate limits",
		Aliases: []string{"rl"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return accountRateLimit.run()
		},
	})

	return cmd
}

func (ac *accountGet) run() error {
	return nil
}

func (ac *accountRateLimit) run() error {
	return nil
}
