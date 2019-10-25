package volume

import (
	"fmt"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/doctl/commands/displayers"
	"github.com/digitalocean/doctl/do"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewList returns a new volume command thing.
func NewList(v *viper.Viper, client do.VolumesService, displayer displayers.Displayer) *cobra.Command {
	volumeList := &VolumeList{
		Service:   client,
		Displayer: displayer,
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "list",
		Aliases: []string{"ls"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return volumeList.PreRun(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return volumeList.Run()
		},
	}

	set := listCmd.Flags()
	set.StringVarP(&volumeList.Region, doctl.ArgRegionSlug, "", "", "Volume Region")

	return listCmd
}

type VolumeList struct {
	Globs     []string
	Region    string
	Displayer displayers.Displayer
	Service   do.VolumesService
}

func (vl *VolumeList) PreRun(args []string) error {
	vl.Globs = args
	return nil
}

func (vl *VolumeList) Run() error {
	var matches []glob.Glob
	for _, globStr := range vl.Globs {
		g, err := glob.Compile(globStr)
		if err != nil {
			return fmt.Errorf("unknown glob %q", globStr)
		}

		matches = append(matches, g)
	}

	list, err := vl.Service.List()
	if err != nil {
		return err
	}
	var matchedList []do.Volume

	for _, volume := range list {
		var skip = true
		if len(matches) == 0 {
			skip = false
		} else {
			for _, m := range matches {
				if m.Match(volume.ID) {
					skip = false
				}
				if m.Match(volume.Name) {
					skip = false
				}
			}
		}

		if !skip && vl.Region != "" {
			if vl.Region != volume.Region.Slug {
				skip = true
			}
		}

		if !skip {
			matchedList = append(matchedList, volume)
		}
	}

	item := &displayers.Volume{Volumes: matchedList}
	return vl.Displayer.DisplayBetter(item)
}
