package devtools

import (
	"context"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/uexky"
)

var devFlags struct {
	email        string
	method       string
	userCount    int
	threadCount  int
	maxPostCount int
	minPostCount int
}

func init() {
	signInUserCmd.PersistentFlags().StringVar(&devFlags.email, "email", "", "user email")
	signInUserCmd.PersistentFlags().StringVar(&devFlags.method, "method", "url", "user email")
	Command.AddCommand(signInUserCmd, SetMainTagsCmd, rebuildCmd, mockDataCmd)
}

func mapArgs(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

var Command = &cobra.Command{
	Use:   "dev",
	Short: "dev utils",
	Long:  "utils for develop",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if config.GetEnv() != config.DevEnv {
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
		if devFlags.email == "" {
			log.Fatalf("must specify user email")
		}
		service, err := uexky.InitDevService()
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		ctx = service.TxAdapter.AttachDB(ctx)
		code, err := service.TrySignInByEmail(ctx, devFlags.email)
		if err != nil {
			log.Fatal(err)
		}
		if devFlags.method == "url" {
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

var SetMainTagsCmd = &cobra.Command{
	Use:   "settags anime,game",
	Short: "set main tags, separated by comma",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tags := strings.Split(args[0], ",")
		tagsTrimmed := mapArgs(tags, strings.TrimSpace)
		service, err := uexky.InitDevService()
		if err != nil {
			log.Fatal(err)
		}
		ctx := service.TxAdapter.AttachDB(context.Background())
		if err := service.SetMainTags(ctx, tagsTrimmed); err != nil {
			log.Fatal(err)
		}
	},
}

var rebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "rebuild dev database",
	Run: func(cmd *cobra.Command, args []string) {
		if err := postgres.RebuildDB(); err != nil {
			log.Fatal(err)
		}
	},
}
