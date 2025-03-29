package cmd

import (
	"fmt"
	"github.com/peyzor/todo-cli/storage"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		rows, err := storage.GetRows(f)
		if err != nil {
			fmt.Fprintf(os.Stdout, "couldn't get rows: %v", err)
		}

		for _, row := range rows {
			fmt.Fprintln(w, strings.Join(row, "\t"))
		}

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
