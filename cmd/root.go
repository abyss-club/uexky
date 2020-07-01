package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/abyss.club/uexky/cmd/devtools"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/server"
)

func init() {
	cobra.OnInitialize(initLog, initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	rootCmd.AddCommand(migrateCmd, devtools.Command, uidCmd)
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func initConfig() {
	if err := config.Load(cfgFile); err != nil {
		log.Fatal(err)
	}
	log.Infof("run with config:\n%#v", config.Get())
}

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "uexky",
	Short: "Uexky is backend program of abyss",
	Long: `Abyss is an anoymouse-able and tagged-thread forum.
                uexky is backend program of abyss.`,
	Run: func(cmd *cobra.Command, args []string) {
		runService()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Info("Done.")
	},
}

func runService() {
	server, err := server.InitProdServer()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.Run())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
