package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/ibihim/banking-csv-cli/pkg/sql"
)

const (
	defaultDBPath = "./transactions.db"

	filenameFlag   = "filename"
	dbFlag         = "db"
	migrationsFlag = "migrations"
)

func BankingCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "banking",
		Short: "A tool to parse banking csv files",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			flag.CommandLine.VisitAll(func(flag *flag.Flag) {
				klog.V(4).Infof("Flag: --%s=%q", flag.Name, flag.Value)
			})
		},
	}

	// Init klog files
	fs := flag.NewFlagSet("", flag.PanicOnError)
	klog.InitFlags(fs)
	rootCmd.PersistentFlags().AddGoFlagSet(fs)

	appCmd := &cobra.Command{
		Use:   "app",
		Short: "Group transactions by purpose",
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath, err := cmd.Flags().GetString(dbFlag)
			if err != nil {
				return fmt.Errorf("failed to get dbFlag: %w", err)
			}

			// Init db
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			db := sql.NewDatabase(&sql.DatabaseOptions{
				URL: dbPath,
			})
			if err := db.Connect(ctx); err != nil {
				return fmt.Errorf("failed on db connect: %w", err)
			}
			defer db.Close()

			// Load transactions
			ts, err := db.GetTransactions()
			if err != nil {
				return fmt.Errorf("failed to load transactions: %w", err)
			}

			return RunApp(ts)
		},
	}
	appCmd.Flags().String(dbFlag, defaultDBPath, "Path to the database file")
	rootCmd.AddCommand(appCmd)

	// dbCmd represents the `db` subcommand
	dbCmd := &cobra.Command{
		Use:   "db",
		Short: "Database operations",
	}
	rootCmd.AddCommand(dbCmd)

	loadCmd := &cobra.Command{
		Use:   "load",
		Short: "Load transactions into the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename, err := cmd.Flags().GetString(filenameFlag)
			if err != nil {
				return err
			}

			reader, err := os.Open(filename)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer reader.Close()

			ts, err := ParseTransactions(reader)
			if err != nil {
				return fmt.Errorf("failed to parse transactions: %w", err)
			}

			// Init db
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			db := sql.NewDatabase(&sql.DatabaseOptions{
				URL: defaultDBPath,
			})
			if err := db.Connect(ctx); err != nil {
				return fmt.Errorf("failed on db connect: %w", err)
			}
			defer db.Close()

			for _, t := range ts {
				ok, err := db.HasTransaction(t)
				if err != nil {
					return fmt.Errorf("failed to check if transaction exists: %w", err)
				}
				if ok {
					return fmt.Errorf("transaction already exists: %+v", t)
				}

				_, err = db.AddTransaction(t)
				if err != nil {
					return fmt.Errorf("failed to add transaction (%+v): %w", t, err)
				}
			}

			return nil
		},
	}

	loadCmd.Flags().String(filenameFlag, "", "The path to the csv file")
	dbCmd.AddCommand(loadCmd)

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath, err := cmd.Flags().GetString(dbFlag)
			if err != nil {
				return fmt.Errorf("failed to get dbFlag: %w", err)
			}
			migrationsPath, err := cmd.Flags().GetString(migrationsFlag)
			if err != nil {
				return fmt.Errorf("failed to get migrationsFlag: %w", err)
			}

			if err := validateMigrateOptions(migrationsPath); err != nil {
				return fmt.Errorf("failed to validate migrationsPath: %w", err)
			}

			if err := completeMigrateOptions(dbPath); err != nil {
				return fmt.Errorf("failed to complate migrationOptions: %w", err)
			}

			return RunMigrate(dbPath, migrationsPath)
		},
	}

	migrateCmd.Flags().String(dbFlag, defaultDBPath, "Path to the database file")
	migrateCmd.Flags().String(migrationsFlag, "pkg/sql/migrations", "Path to the database file")
	dbCmd.AddCommand(migrateCmd)

	return rootCmd
}
