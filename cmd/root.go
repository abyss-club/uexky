package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var rootCmd = &cobra.Command{
	Use:   "uexky",
	Short: "Uexky is backend program of abyss",
	Long: `Abyss is an anoymouse-able and tagged-thread forum.
                uexky is backend program of abyss.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("hello, world")
		runService()
	},
}

func runService() {
	fmt.Println("hello, world")
	// port := "8000"
	// resolver := InitResolver()
	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}))
	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// http.Handle("/query", srv)
	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
