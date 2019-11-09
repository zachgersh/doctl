package ratelimit

import (
	"github.com/digitalocean/doctl/commands/displayers"
	"github.com/digitalocean/doctl/do"
	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

type Base struct {
	displayer displayers.Displayer
	service   do.AccountService
	v         *viper.Viper
	ui        cli.Ui
}

//New returns our base rate limit subcommand
func New(v *viper.Viper, client do.AccountService, displayer displayers.Displayer, ui cli.Ui) *Base {
	return &Base{
		displayer: displayer,
		service:   client,
		v:         v,
		ui:        ui,
	}
}

func (b *Base) Run(args []string) int {
	rateLimit, err := b.service.RateLimit()
	if err != nil {
		b.ui.Error(err.Error())
		return 1
	}

	content, err := b.displayer.DisplayBetter(&displayers.RateLimit{RateLimit: rateLimit})
	if err != nil {
		b.ui.Error(err.Error())
		return 1
	}

	b.ui.Output(content)

	return 0
}

func (b *Base) Synopsis() string {
	return "get API rate limits"
}

func (b *Base) Help() string {
	return ""
}
