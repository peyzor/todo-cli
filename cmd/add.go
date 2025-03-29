package cmd

import (
	"fmt"
	"github.com/peyzor/todo-cli/storage"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
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
		nextID, err := storage.GetNextID(f)
		if err != nil {
			fmt.Printf("couldn't get next id: %v", err)
			return
		}

		err = storage.AddNewCSVRecord(f, []string{strconv.Itoa(nextID), todo, time.Now().UTC().String(), storage.IsDoneNo})
		if err != nil {
			fmt.Printf("couldn't add new record: %v", err)
			return
		}

		f.Seek(0, io.SeekStart)

		err = storage.GetRowsTabular(f, os.Stdout)
		if err != nil {
			fmt.Printf("couldn't get rows: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
