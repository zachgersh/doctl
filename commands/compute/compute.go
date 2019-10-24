package compute

import (
	"github.com/spf13/cobra"
)

// NewCommand returns a new wrapper or whatever we decide.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "compute",
		Short: "what ever",
	}
}
