package droplet

import (
	"strconv"

	"github.com/digitalocean/doctl/commands/displayers"
	"github.com/digitalocean/doctl/do"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DropletGet struct {
	ID        int
	Displayer displayers.Displayer
	Service   do.DropletsService
}

// NewCommand returns a new Droplet command thing.
func NewCommand(v *viper.Viper, client do.DropletsService, displayer displayers.Displayer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "droplet",
		Aliases: []string{"d"},
		Short:   "droplet commands",
		Long:    "droplet is used to access droplet commands",
	}

	dropletGet := &DropletGet{
		Service:   client,
		Displayer: displayer,
	}
	// AddStringFlag(cmdRunDropletGet, doctl.ArgTemplate, "", "", "Go template format. Few sample values:{{.ID}} {{.Name}} {{.Memory}} {{.Region.Name}} {{.Image}} {{.Tags}}")
	getCmd := &cobra.Command{
		Use:     "get",
		Short:   "get droplet",
		Aliases: []string{"g"},
		PreRunE: dropletGet.PreRun,
		RunE:    dropletGet.Run,
	}
	cmd.AddCommand(getCmd)

	return cmd
}

func (dg *DropletGet) PreRun(_ *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	dg.ID = id
	return nil
	// dropletGet.Template = v.Get(doctl.ArgTemplate)
}

func (dg *DropletGet) Run(_ *cobra.Command, _ []string) error {
	droplet, err := dg.Service.Get(dg.ID)
	if err != nil {
		return err
	}
	return dg.Displayer.DisplayBetter(&displayers.Droplet{Droplets: do.Droplets{*droplet}})
}
