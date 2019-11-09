package account

import "github.com/mitchellh/cli"

type Base struct{}

//New returns our base account command that only holds help
func New() *Base {
	return &Base{}
}

func (b *Base) Run(args []string) int {
	return cli.RunResultHelp
}

func (b *Base) Synopsis() string {
	return "account commands"
}

func (b *Base) Help() string {
	return ""
}
