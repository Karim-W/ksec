package cmd

import (
	"os"

	"github.com/karim-w/ksec/models"
	"github.com/karim-w/ksec/service"
	"github.com/spf13/cobra"
	cobracompletefig "github.com/withfig/autocomplete-tools/integrations/cobra"
)

const VERSION = "0.0.1"

var RootCmd = &cobra.Command{
	Use:   "ksec",
	Short: "ksec is a tool for managing secrets in Kubernetes",
	Long:  "ksec is a tool for managing secrets in Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		c := &models.Secrets{
			Namespace:  cmd.Flag("namespace").Value.String(),
			Secret:     cmd.Flag("secret").Value.String(),
			Set:        cmd.Flag("set").Value.String() == "true",
			Key:        cmd.Flag("key").Value.String(),
			Value:      cmd.Flag("value").Value.String(),
			Get:        cmd.Flag("get").Value.String() == "true",
			Delete:     cmd.Flag("delete").Value.String() == "true",
			List:       cmd.Flag("list").Value.String() == "true",
			All:        cmd.Flag("all").Value.String() == "true",
			EnvPath:    cmd.Flag("env").Value.String(),
			FillPath:   cmd.Flag("fill").Value.String(),
			Modify:     cmd.Flag("modify").Value.String() == "true",
			FileFormat: cmd.Flag("file-format").Value.String(),
		}
		service.KubectlSecretsSvc(c)
	},
}

var (
	Verbose bool
	Source  string
)

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(cobracompletefig.CreateCompletionSpecCommand())
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringP("namespace", "n", "default", "namespace")
	RootCmd.PersistentFlags().StringP("secret", "s", "", "secret name")
	RootCmd.PersistentFlags().BoolP("set", "w", false, "set secret value")
	RootCmd.PersistentFlags().StringP("key", "k", "", "secret key")
	RootCmd.PersistentFlags().StringP("value", "V", "", "secret value")
	RootCmd.PersistentFlags().BoolP("get", "g", false, "get secret value")
	RootCmd.PersistentFlags().BoolP("delete", "d", false, "delete secret")
	RootCmd.PersistentFlags().BoolP("list", "l", false, "list secrets")
	RootCmd.PersistentFlags().BoolP("all", "a", false, "list all secrets")
	RootCmd.PersistentFlags().StringP("env", "e", "", "Create from a .env file")
	RootCmd.PersistentFlags().StringP("fill", "f", "", "Fill a file with secrets")
	RootCmd.PersistentFlags().BoolP("modify", "m", false, "Modify a secret in an interactive mode")
	RootCmd.PersistentFlags().StringP("file-format", "F", "yaml", "File format")
}
