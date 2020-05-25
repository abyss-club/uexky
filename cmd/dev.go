package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/wire"
)

var (
	email  string
	method string
)

func init() {
	signInUserCmd.PersistentFlags().StringVar(&email, "email", "", "user email")
	signInUserCmd.PersistentFlags().StringVar(&method, "method", "url", "user email")
	devCmd.AddCommand(signInUserCmd)
}

// develop utils
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "dev utils",
	Long:  "utils for develop",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if !strings.Contains(cfgFile, "dev") && !strings.Contains(config.Get().PostgresURI, "dev") {
			fmt.Println("Ensure you are in dev environment!")
			fmt.Println("Let your config files or postgres uri contains 'dev'")
			os.Exit(1)
		}
	},
}

var signInUserCmd = &cobra.Command{
	Use:   "signin",
	Short: "sign in user",
	Run: func(cmd *cobra.Command, args []string) {
		if email == "" {
			log.Fatalf("must specify user email")
		}
		service, err := wire.InitDevService()
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		ctx = service.TxAdapter.AttachDB(ctx)
		code, err := service.GenSignInCodeByEmail(ctx, email)
		if err != nil {
			log.Fatal(err)
		}
		if method == "url" {
			fmt.Println("Sign In URL: ", code.SignInURL())
			return
		}
		token, err := service.SignInByCode(ctx, string(code))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Token cookie: ", token.Cookie().String())
	},
}
