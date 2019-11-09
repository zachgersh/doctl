package volume

import "github.com/mitchellh/cli"

type Base struct{}

//New returns our base volume command that only holds help
func New() *Base {
	return &Base{}
}

func (b *Base) Run(args []string) int {
	return cli.RunResultHelp
}

func (b *Base) Synopsis() string {
	return ""
}

func (b *Base) Help() string {
	return ""
}
