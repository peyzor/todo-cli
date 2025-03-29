package cmd

import (
	"fmt"
	"github.com/peyzor/todo-cli/storage"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "delete a todo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := storage.GetOrCreateCSVStorage()
		if err != nil {
			fmt.Printf("couldn't get storage: %v", err)
			return
		}
		defer f.Close()

		IDStr := args[0]
		ID, err := strconv.Atoi(IDStr)
		if err != nil {
			fmt.Printf("invalid ID: %v", err)
			return
		}

		err = storage.DeleteCSVRecord(f, ID)
		if err != nil {
			fmt.Printf("couldn't delete record: %v", err)
			return
		}

		f, err = storage.GetOrCreateCSVStorage()
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
	rootCmd.AddCommand(deleteCmd)
}
