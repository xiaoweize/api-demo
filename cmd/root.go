/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xiaoweize/api-demo/version"
)

var vers bool

var rootCmd = &cobra.Command{
	Use:   "demo",
	Short: "demo host-api",
	Long:  `demo host-api impl`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			//打印版本信息
			fmt.Println(version.FullVersion())
		} else {
			//打印帮助信息
			return cmd.Usage()
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, "print version")
}
