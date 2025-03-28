package cmd

import (
	"fmt"
	"github.com/peyzor/todo-cli/storage"
	"github.com/spf13/cobra"
	"time"
)

var addCmd = &cobra.Command{
	Use:   "add <todo>",
	Short: "Adds todo to the list of todos",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := storage.GetOrCreateCSVStorage()
		if err != nil {
			fmt.Printf("couldn't get storage: %v", err)
			return
		}
		defer f.Close()

		todo := args[0]
		err = storage.AddNewCSVRecord(f, []string{"1", todo, time.Now().UTC().String()})
		if err != nil {
			fmt.Printf("couldn't add new record: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
