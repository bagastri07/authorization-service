package console

import (
	"os"

	"github.com/bagastri07/boilerplate-service/internal/config"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "boilerplate-service",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func setupLogger() {
	formatter := runtime.Formatter{
		ChildFormatter: &logrus.JSONFormatter{},
		Line:           true,
		File:           true,
	}

	if config.Env() == "development" {
		formatter = runtime.Formatter{
			ChildFormatter: &logrus.TextFormatter{
				ForceColors:   true,
				FullTimestamp: true,
			},
			Line: true,
			File: true,
		}
	}

	logrus.SetFormatter(&formatter)
	logrus.SetOutput(os.Stdout)

	logLevel, err := logrus.ParseLevel(config.LogLevel())
	if err != nil {
		logLevel = logrus.DebugLevel
	}
	logrus.SetLevel(logLevel)
}

func init() {
	config.GetConf()
	setupLogger()
}
