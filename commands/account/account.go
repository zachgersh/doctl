package account

import (
	"github.com/digitalocean/doctl/commands/displayers"
	"github.com/digitalocean/doctl/do"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AccountGet struct {
	Displayer displayers.Displayer
	Service   do.AccountService
}

type AccountRateLimit struct {
	Displayer displayers.Displayer
	Service   do.AccountService
}

//NewAccountCmd ghasdfhjasdhfj
func NewAccountCmd(v *viper.Viper, client do.AccountService, displayer displayers.Displayer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "account commands",
		Long:  "account is used to access account commands",
	}

	accountGet := &AccountGet{
		Service:   client,
		Displayer: displayer,
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "get",
		Short:   "get account",
		Aliases: []string{"get"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return accountGet.Run()
		},
	})

	accountRateLimit := &AccountRateLimit{
		Service:   client,
		Displayer: displayer,
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "ratelimit",
		Short:   "get API rate limits",
		Aliases: []string{"rl"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return accountRateLimit.Run()
		},
	})

	return cmd
}

func (ac *AccountGet) Run() error {
	account, err := ac.Service.Get()
	if err != nil {
		return err
	}
	return ac.Displayer.DisplayBetter(&displayers.Account{Account: account})
}

func (ac *AccountRateLimit) Run() error {
	rateLimit, err := ac.Service.RateLimit()
	if err != nil {
		return err
	}
	return ac.Displayer.DisplayBetter(&displayers.RateLimit{RateLimit: rateLimit})
}
