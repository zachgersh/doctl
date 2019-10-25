package compute

import (
	"github.com/digitalocean/doctl/commands/compute/droplet"
	"github.com/digitalocean/doctl/commands/compute/volume"
	"github.com/digitalocean/doctl/commands/displayers"
	"github.com/digitalocean/doctl/do"
	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	base        = &cobra.Command{Use: "compute"}
	subCommands = map[string]*cobra.Command{
		"droplet": &cobra.Command{
			Use:     "droplet",
			Aliases: []string{"d"},
			Short:   "droplet commands",
			Long:    "droplet is used to access droplet commands",
		},
		"volume": &cobra.Command{
			Use:   "volume",
			Short: "volume commands",
			Long:  "volume is used to access volume commands",
		},
	}
)

// NewCommand returns a new wrapper or whatever we decide.
func NewCommand(v *viper.Viper, client *godo.Client, displayer displayers.Displayer) *cobra.Command {
	for name, command := range subCommands {
		var attach []*cobra.Command
		base.AddCommand(command)

		switch name {
		case "droplet":
			get := droplet.NewGet(v, do.NewDropletsService(client), displayer)
			list := droplet.NewList(v, do.NewDropletsService(client), displayer)
			attach = append(attach, get, list)
		case "volume":
			list := volume.NewList(v, do.NewVolumesService(client), displayer)
			attach = append(attach, list)
		}

		command.AddCommand(attach...)
	}

	return base
}
