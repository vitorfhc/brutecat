package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var banner = `|\__/,|     ('\
|_ _  |.----.) )
( T   )       /
(((^_(((/(((_/
B R U T E C A T
`

type BruteCatOptions struct {
	Threads       uint16
	UsersFile     string
	PasswordsFile string
}

type CLIOptions struct {
	ContinueOnSuccess bool
	StatsEvery        uint16
	Ctx               context.Context
}

var bruteCatOptions BruteCatOptions
var cliOptions CLIOptions

var rootCmd = &cobra.Command{
	Use:   "brutecat",
	Short: "Brute force services with ease",
	Long:  "",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}

		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		fmt.Println(banner)

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		logrus.Info("Done")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug output")

	rootCmd.PersistentFlags().Uint16VarP(&bruteCatOptions.Threads, "threads", "t", 4, "Number of threads to use")

	rootCmd.PersistentFlags().StringVarP(&bruteCatOptions.UsersFile, "users", "u", "users.txt", "File containing users to brute force")
	rootCmd.PersistentFlags().StringVarP(&bruteCatOptions.PasswordsFile, "passwords", "p", "passwords.txt", "File containing passwords to brute force")

	rootCmd.PersistentFlags().BoolVarP(&cliOptions.ContinueOnSuccess, "continue-on-success", "c", false, "Continue brute forcing even after finding a valid credential")
	rootCmd.PersistentFlags().Uint16Var(&cliOptions.StatsEvery, "stats-every", 10, "Print stats every N seconds")
}
