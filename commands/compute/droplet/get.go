package droplet

import (
	"strconv"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/doctl/commands/displayers"
	"github.com/digitalocean/doctl/do"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewGet returns a new Droplet command thing.
func NewGet(v *viper.Viper, client do.DropletsService, displayer displayers.Displayer) *cobra.Command {
	dropletGet := &DropletGet{
		Service:   client,
		Displayer: displayer,
	}

	getCmd := &cobra.Command{
		Use:     "get <droplet-id>",
		Short:   "get droplet",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return dropletGet.PreRun(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return dropletGet.Run()
		},
	}

	set := getCmd.Flags()
	set.StringVarP(&dropletGet.Template, doctl.ArgTemplate, "", "", "Go template format. Few sample values:{{.ID}} {{.Name}} {{.Memory}} {{.Region.Name}} {{.Image}} {{.Tags}}")

	return getCmd
}

type DropletGet struct {
	ID        int
	Template  string
	Displayer displayers.Displayer
	Service   do.DropletsService
}

func (dg *DropletGet) PreRun(args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	dg.ID = id
	return nil
}

func (dg *DropletGet) Run() error {
	droplet, err := dg.Service.Get(dg.ID)
	if err != nil {
		return err
	}

	return dg.Displayer.DisplayBetter(&displayers.Droplet{Droplets: do.Droplets{*droplet}})
}
