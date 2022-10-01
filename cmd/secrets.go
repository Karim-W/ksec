package cmd

import (
	"github.com/karim-w/ksec/models"
	"github.com/karim-w/ksec/service"
	"github.com/spf13/cobra"
)

var SecCmd = &cobra.Command{
	Use:     "s",
	Aliases: []string{"sec", "secretv", "secrets"},
	Short:   "Manage secrets",
	Long:    "Manage you kubernetes cluster secrets",
	Run: func(cmd *cobra.Command, args []string) {
		c := &models.Secrets{
			Namespace: cmd.Flag("namespace").Value.String(),
			Secret:    cmd.Flag("secret").Value.String(),
			Set:       cmd.Flag("set").Value.String() == "true",
			Key:       cmd.Flag("key").Value.String(),
			Value:     cmd.Flag("value").Value.String(),
			Get:       cmd.Flag("get").Value.String() == "true",
			Delete:    cmd.Flag("delete").Value.String() == "true",
			List:      cmd.Flag("list").Value.String() == "true",
			All:       cmd.Flag("all").Value.String() == "true",
			EnvPath:   cmd.Flag("env").Value.String(),
		}
		service.KubectlSecretsSvc(c)
	},
}

func init() {
	RootCmd.AddCommand(SecCmd)
	SecCmd.PersistentFlags().StringP("namespace", "n", "default", "namespace")
	SecCmd.PersistentFlags().StringP("secret", "s", "", "secret name")
	SecCmd.PersistentFlags().BoolP("set", "w", false, "set secret value")
	SecCmd.PersistentFlags().StringP("key", "k", "", "secret key")
	SecCmd.PersistentFlags().StringP("value", "V", "", "secret value")
	SecCmd.PersistentFlags().BoolP("get", "g", false, "get secret value")
	SecCmd.PersistentFlags().BoolP("delete", "d", false, "delete secret")
	SecCmd.PersistentFlags().BoolP("list", "l", false, "list secrets")
	SecCmd.PersistentFlags().BoolP("all", "a", false, "list all secrets")
	SecCmd.PersistentFlags().StringP("env", "e", "", "environment")

}
