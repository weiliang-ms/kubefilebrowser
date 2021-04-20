package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func init() {
	rootCmd.AddCommand(touchCmd)
}

// touchCmd represents the touch command
var touchCmd = &cobra.Command{
	Use:     "touch",
	Aliases: []string{"echo"},
	Short:   "Get content from standard input and write it to file.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Create(args[0])
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(253)
		}
		_, err = io.Copy(f, os.Stdin)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(253)
		}
	},
}
