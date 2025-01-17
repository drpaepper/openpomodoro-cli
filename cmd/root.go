package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/drpaepper/go-openpomodoro"
	"github.com/drpaepper/openpomodoro-cli/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pomodoro",
	Short: "Pomodoro CLI",
	Long:  "A simple Pomodoro command-line client for the Open Pomodoro format",
}

var (
	client   *openpomodoro.Client
	settings *openpomodoro.Settings

	directoryFlag string
	formatFlag    string
	waitFlag      bool
)

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(
		&directoryFlag, "directory", "", ``,
		"directory to read/write Open Pomodoro data (default is ~/.config/pomodoro/)")

	RootCmd.PersistentFlags().StringVarP(
		&formatFlag, "format", "f", format.DefaultFormat,
		"format to display Pomodoros in")

	RootCmd.PersistentFlags().BoolVarP(
		&waitFlag, "wait", "w", false,
		"wait for the Pomodoro to end before exiting")

	viper.AutomaticEnv()

	var err error

	client, err = openpomodoro.NewClient(directoryFlag)
	if err != nil {
		log.Fatalf("Could not create client: %v", err)
	}

	settings, err = client.Settings()
	if err != nil {
		log.Fatalf("Could not retrieve settings: %v", err)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
