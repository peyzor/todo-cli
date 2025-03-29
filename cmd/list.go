package cmd

import (
	"fmt"
	"github.com/peyzor/todo-cli/storage"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all todos",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := storage.GetOrCreateCSVStorage()
		if err != nil {
			fmt.Printf("couldn't get storage: %v", err)
			return
		}
		defer f.Close()

		err = storage.GetRowsTabular(f, os.Stdout)
		if err != nil {
			fmt.Printf("couldn't get rows: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
