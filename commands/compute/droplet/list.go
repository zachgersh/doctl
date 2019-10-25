package droplet

import (
	"fmt"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/doctl/commands/displayers"
	"github.com/digitalocean/doctl/do"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewList returns a new Droplet command thing.
func NewList(v *viper.Viper, client do.DropletsService, displayer displayers.Displayer) *cobra.Command {
	dropletList := &DropletList{
		Service:   client,
		Displayer: displayer,
	}

	listCmd := &cobra.Command{
		Use:     "list [GLOB]",
		Short:   "list droplets",
		Aliases: []string{"ls"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return dropletList.PreRun(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return dropletList.Run()
		},
	}

	set := listCmd.Flags()
	set.StringVarP(&dropletList.Region, doctl.ArgRegionSlug, "", "", "Droplet Region")
	set.StringVarP(&dropletList.TagName, doctl.ArgTagName, "", "", "Tag Name")

	return listCmd
}

type DropletList struct {
	Region    string
	TagName   string
	Globs     []string
	Displayer displayers.Displayer
	Service   do.DropletsService
}

func (dl *DropletList) PreRun(args []string) error {
	dl.Globs = args
	return nil
}

func (dl *DropletList) Run() error {
	var matches []glob.Glob

	for _, globStr := range dl.Globs {
		g, err := glob.Compile(globStr)
		if err != nil {
			return fmt.Errorf("unknown glob %q", globStr)
		}

		matches = append(matches, g)
	}

	var (
		matchedList do.Droplets
		list        do.Droplets
		err         error
	)

	if dl.TagName == "" {
		list, err = dl.Service.List()
	} else {
		list, err = dl.Service.ListByTag(dl.TagName)
	}
	if err != nil {
		return err
	}

	for _, droplet := range list {
		var skip = true
		if len(matches) == 0 {
			skip = false
		} else {
			for _, m := range matches {
				if m.Match(droplet.Name) {
					skip = false
				}
			}
		}

		if !skip && dl.Region != "" {
			if dl.Region != droplet.Region.Slug {
				skip = true
			}
		}

		if !skip {
			matchedList = append(matchedList, droplet)
		}
	}

	item := &displayers.Droplet{Droplets: matchedList}
	return dl.Displayer.DisplayBetter(item)
}
