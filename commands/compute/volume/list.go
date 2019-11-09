package volume

import (
	"fmt"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/doctl/commands/displayers"
	"github.com/digitalocean/doctl/do"
	"github.com/gobwas/glob"
	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

//List is in charge holds the region / flag set data required to retrieve
//volume data.
type List struct {
	region    string
	displayer displayers.Displayer
	service   do.VolumesService
	v         *viper.Viper
	ui        cli.Ui
	fs        *flag.FlagSet
}

// NewList returns a new volume command thing.
func NewList(v *viper.Viper, client do.VolumesService, displayer displayers.Displayer, ui cli.Ui) *List {
	list := &List{
		displayer: displayer,
		service:   client,
		v:         v,
		ui:        ui,
	}

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVarP(&list.region, doctl.ArgRegionSlug, "", "", "Volume Region")

	list.fs = fs

	return list
}

//Run parses the provided flags before returning the rest of the non
//flag arguments to be used as globs when searching for a matching
//volume.
func (l *List) Run(args []string) int {
	if err := l.fs.Parse(args); err != nil {
		return 1
	}

	globs := l.fs.Args()

	var matches []glob.Glob
	for _, globStr := range globs {
		g, err := glob.Compile(globStr)
		if err != nil {
			l.ui.Error(fmt.Sprintf("unknown glob %q", globStr))
			return 1
		}

		matches = append(matches, g)
	}

	list, err := l.service.List()
	if err != nil {
		l.ui.Error(fmt.Sprintf("failed to list volumes: %s", err))
		return 1
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

		if !skip && l.region != "" {
			if l.region != volume.Region.Slug {
				skip = true
			}
		}

		if !skip {
			matchedList = append(matchedList, volume)
		}
	}

	item := &displayers.Volume{Volumes: matchedList}
	content, err := l.displayer.DisplayBetter(item)
	if err != nil {
		l.ui.Error(fmt.Sprintf("failed to display volumes: %s", err))
		return 1
	}

	l.ui.Output(content)

	return 0
}
