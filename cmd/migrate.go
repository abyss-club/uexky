package cmd

import (
	"fmt"
	"path/filepath"

	_ "github.com/go-pg/pg/v9" // postgres driver
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // migrate file source
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/abyss.club/uexky/lib/config"
)

var migrateFilesPath string

var (
	mgr     *migrate.Migrate
	version uint
	nsteps  int
)

func init() {
	migrateCmd.PersistentFlags().StringVarP(&migrateFilesPath, "source", "s", "", "migration files path")
	migrateCmd.AddCommand(gotoCmd, upCmd, downCmd, forceCmd, versionCmd)
	gotoCmd.PersistentFlags().UintVarP(&version, "version", "v", 0, "version")
	upCmd.PersistentFlags().IntVarP(&nsteps, "nsteps", "n", 0, "nsteps")
	downCmd.PersistentFlags().IntVarP(&nsteps, "nsteps", "n", 0, "nsteps")
	forceCmd.PersistentFlags().UintVarP(&version, "version", "v", 0, "version")
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate tools",
	Long:  "migrate tools for uexky postgres",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if migrateFilesPath == "" {
			migrateFilesPath = "./migrates"
		}
		path, err := filepath.Abs(migrateFilesPath)
		if err != nil {
			log.Fatal(fmt.Errorf("parse migration file path: %w", err))
		}
		source := fmt.Sprintf("file://%s", path)
		log.Infof("migrates source is: %s", source)
		m, err := migrate.New(source, config.Get().PostgresURI)
		if err != nil {
			log.Fatal(err)
		}
		mgr = m
	},
}

var gotoCmd = &cobra.Command{
	Use:   "goto -v version",
	Short: "migrate to version",
	Run: func(cmd *cobra.Command, args []string) {
		if version == 0 {
			log.Fatal(fmt.Errorf("invalid version %v", args[0]))
		}
		if err := mgr.Migrate(version); err != nil {
			log.Fatal(fmt.Errorf("migrate goto failed: %w", err))
		}
	},
}

var upCmd = &cobra.Command{
	Use:   "up [-n N]",
	Short: "apply all or N up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if nsteps == 0 {
			if err := mgr.Up(); err != nil {
				log.Fatal(fmt.Errorf("migrate up failed: %w", err))
			}
		} else {
			if err := mgr.Steps(nsteps); err != nil {
				log.Fatal(fmt.Errorf("migrate up failed: %w", err))
			}
		}
	},
}

var downCmd = &cobra.Command{
	Use:   "down [-n N]",
	Short: "apply all or N down migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if nsteps == 0 {
			if err := mgr.Down(); err != nil {
				log.Fatal(fmt.Errorf("migrate down failed: %w", err))
			}
		} else {
			if err := mgr.Steps(-nsteps); err != nil {
				log.Fatal(fmt.Errorf("migrate down failed: %w", err))
			}
		}
	},
}

var forceCmd = &cobra.Command{
	Use:   "force -v version",
	Short: "set version but don't run migration (ignores dirty state)",
	Run: func(cmd *cobra.Command, args []string) {
		if version == 0 {
			log.Fatal(fmt.Errorf("invalid version %v", args[0]))
		}
		if err := mgr.Force(int(version)); err != nil {
			log.Fatal(fmt.Errorf("migrate force version failed: %w", err))
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print current migration version",
	Run: func(cmd *cobra.Command, args []string) {
		ver, dirty, err := mgr.Version()
		if err != nil {
			log.Fatal(fmt.Errorf("migrate display version failed: %w", err))
		}
		fmt.Printf("version: %v\nis dirty: %v\n", ver, dirty)
	},
}
