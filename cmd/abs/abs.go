/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package abs

import (
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var AbsCmd = &cobra.Command{
	Use:   "abs",
	Short: "Commands for interacting with Audiobookshelf",
}

func init() {
	AbsCmd.AddCommand(libraryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
