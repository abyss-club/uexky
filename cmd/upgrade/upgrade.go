package upgrade

import (
	"fmt"
	"os"

	"github.com/go-pg/pg/v9"
	"github.com/spf13/cobra"
	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/lib/errors"
)

var prevDB string

func init() {
	Command.PersistentFlags().StringVar(&prevDB, "prev", "", "db address of prev version")
}

var Command = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade database from v1",
	Long:  "import old version data and save to new database",
	Run: func(cmd *cobra.Command, args []string) {
		if err := upgrade(); err != nil {
			os.Exit(1)
		}
	},
}

func upgrade() error {
	newDB := config.Get().PostgresURI
	if newDB == "" || prevDB == "" {
		return errors.BadParams.New("you should specified database")
	}
	fmt.Println("prev db: ", prevDB, "new db:", newDB)
	var migrator Migrator
	var err error
	migrator.PrevDB, err = connectDB(prevDB)
	if err != nil {
		return errors.Wrap(err, "connect to prev database")
	}
	migrator.NewDB, err = connectDB(newDB)
	if err != nil {
		return errors.Wrap(err, "connect to new database")
	}
	if err := migrator.DoMigrate(); err != nil {
		return err
	}
	fmt.Println("Version Migrate Succeed!")
	return nil
}

func connectDB(url string) (*pg.DB, error) {
	opt, err := pg.ParseURL(url)
	opt.PoolSize = 16
	if err != nil {
		return nil, errors.BadParams.Handle(err, "parse postgres uri")
	}
	return pg.Connect(opt), nil
}
