/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/peyzor/todo-cli/storage"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete <ID>",
	Short: "A brief description of your command",
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

		err = storage.UpdateCSVRecord(f, ID)
		if err != nil {
			fmt.Printf("couldn't update record: %v", err)
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
	rootCmd.AddCommand(completeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
