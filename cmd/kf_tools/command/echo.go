package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func init() {
	rootCmd.AddCommand(echoCmd)
}

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:     "echo",
	Aliases: []string{"touch"},
	Short:   "Get content from standard input and write it to file.",
	Args: cobra.ExactArgs(1),
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