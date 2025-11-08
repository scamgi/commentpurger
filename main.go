package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "commentpurger [paths...]",
	Short: "CommentPurger removes comments from specified files.",
	Long: `CommentPurger is a command-line tool that removes comments from various file types,
including HTML, CSS, JavaScript, TypeScript, and Vue files.
You can specify one or more files or directories as arguments.`,
	Args: cobra.MinimumNArgs(1),
	Run:  run,
}

func run(_ *cobra.Command, args []string) {
	processPaths(args)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
