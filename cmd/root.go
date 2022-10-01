package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

var RootCmd = &cobra.Command{
	Use:   "ksec",
	Short: "ksec is a tool for managing secrets in Kubernetes",
	Long:  "ksec is a tool for managing secrets in Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("KESC")
		fmt.Println("===============")
		fmt.Println("Welcome To Ksec")
		fmt.Println("How Can I Help")
		fmt.Println("Version: ", VERSION)

	},
}
var Verbose bool
var Source string

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}
