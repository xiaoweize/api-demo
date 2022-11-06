/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "demo",
	Short: "demo host-api",
	Long:  `demo host-api impl`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
