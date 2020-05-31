package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/abyss.club/uexky/lib/uid"
)

var (
	uidParse   string
	uidDisplay int64
)

func init() {
	uidCmd.PersistentFlags().StringVar(&uidParse, "parse", "", "parse string to int64")
	uidCmd.PersistentFlags().Int64Var(&uidDisplay, "disp", 0, "uid display")
}

var uidCmd = &cobra.Command{
	Use:   "uid",
	Short: "uid utils",
	Run: func(cmd *cobra.Command, args []string) {
		if uidParse != "" {
			uid, err := uid.ParseUID(uidParse)
			if err != nil {
				log.Fatal(err)
			}
			log.Infof("%s is uid: %v", uidParse, uid)
		}
		if uidDisplay != 0 {
			str := uid.UID(uidDisplay).ToBase64String()
			log.Infof("uid %v's string: %s", uidDisplay, str)
		}
	},
}
