package cmd

import (
	"context"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	brutecat "github.com/vitorfhc/brutecat/pkg"
	"github.com/vitorfhc/brutecat/pkg/authenticators"
)

var ftpAuthenticator = authenticators.FTP{}

var ftpCmd = &cobra.Command{
	Use:   "ftp",
	Short: "Brute force FTP servers",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Info("Using FTP mode")

		logrus.Info("Starting engine")
		engine, err := brutecat.NewEngineWithFiles(
			bruteCatOptions.UsersFile,
			bruteCatOptions.PasswordsFile,
			bruteCatOptions.Threads,
			&ftpAuthenticator,
		)
		if err != nil {
			return err
		}
		cliOptions.Engine = engine

		bgCtx := context.Background()
		ctx, cancel := context.WithCancel(bgCtx)
		defer cancel()

		engine.OnSuccessCallback = func(creds brutecat.Credentials) {
			logrus.Infof("Found credentials: %s:%s", creds.Username, creds.Password)
			if !cliOptions.ContinueOnSuccess {
				logrus.Info("Stopping engine, please wait")
				cancel()
			}
		}

		engine.OnErrorCallback = func(err error) {
			if strings.Contains(err.Error(), "operation was canceled") || strings.Contains(err.Error(), "RunWithContext") {
				return
			}
			logrus.Error(err)
			cancel()
		}

		engine.OnFailCallback = func(creds brutecat.Credentials) {
			logrus.Debugf("Failed to authenticate with %s:%s", creds.Username, creds.Password)
		}

		logrus.Info("Running engine")
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer logrus.Info("Engine finished")
			engine.Run(ctx)
		}()

		wg.Wait()
		cliOptions.Cancel()
		cliOptions.Wg.Wait()

		return nil
	},
}

func init() {
	ftpCmd.Flags().StringVarP(&ftpAuthenticator.Host, "host", "H", "127.0.0.1", "FTP server host")
	ftpCmd.Flags().Uint16VarP(&ftpAuthenticator.Port, "port", "P", 21, "FTP server port")

	rootCmd.AddCommand(ftpCmd)
}
