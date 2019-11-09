package compute

import "github.com/mitchellh/cli"

type Base struct{}

//New returns our base compute command that only holds help
func New() *Base {
	return &Base{}
}

func (b *Base) Run(args []string) int {
	return cli.RunResultHelp
}

func (b *Base) Synopsis() string {
	return "compute commands"
}

func (b *Base) Help() string {
	return ""
}
