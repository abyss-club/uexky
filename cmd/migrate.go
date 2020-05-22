package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab.com/abyss.club/uexky/config"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres migrate
	_ "github.com/golang-migrate/migrate/v4/source/github"     // migrate file source
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
	gotoCmd.LocalFlags().UintVarP(&version, "version", "v", 0, "version")
	upCmd.LocalFlags().IntVarP(&nsteps, "nsteps", "n", 0, "nsteps")
	downCmd.LocalFlags().IntVarP(&nsteps, "nsteps", "n", 0, "nsteps")
	forceCmd.LocalFlags().UintVarP(&version, "version", "v", 0, "version")
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
	Args:  cobra.MinimumNArgs(1),
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
		if len(args) == 0 {
			if err := mgr.Down(); err != nil {
				log.Fatal(fmt.Errorf("migrate up failed: %w", err))
			}
		} else {
			version, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(fmt.Errorf("invalid version %v: %w", args[0], err))
			}
			if err := mgr.Steps(-version); err != nil {
				log.Fatal(fmt.Errorf("migrate up failed: %w", err))
			}
		}
	},
}
