package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "go-ingest-cli",
		Short: "Ingest records into a Databricks Unity Catalog table via Zerobus",
		Long: `Ingest records into a Databricks Unity Catalog table via Zerobus.

Send a single message inline:
  go-ingest-cli --table catalog.schema.table --message "hello"

Or pipe lines from stdin:
  cat messages.txt | go-ingest-cli --table catalog.schema.table
  echo "single line" | go-ingest-cli --table catalog.schema.table`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ingest(IngestConfig{
				Endpoint:     viper.GetString("endpoint"),
				Host:         viper.GetString("host"),
				ClientID:     viper.GetString("client_id"),
				ClientSecret: viper.GetString("client_secret"),
				Table:        viper.GetString("table"),
				RequestID:    viper.GetString("request_id"),
				Message:      viper.GetString("message"),
			})
		},
	}

	// Define flags.
	rootCmd.Flags().String("endpoint", "", "Zerobus gRPC endpoint")
	rootCmd.Flags().String("host", "", "Databricks workspace URL")
	rootCmd.Flags().String("client-id", "", "Databricks service principal client ID")
	rootCmd.Flags().String("client-secret", "", "Databricks service principal client secret")
	rootCmd.Flags().String("table", "", "Unity Catalog table name")
	rootCmd.Flags().String("request-id", "", "Request ID for ingested messages (default: random UUID)")
	rootCmd.Flags().String("message", "", "Single message to ingest (omit to read from stdin)")

	// Bind flags to Viper.
	viper.BindPFlag("endpoint", rootCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("host", rootCmd.Flags().Lookup("host"))
	viper.BindPFlag("client_id", rootCmd.Flags().Lookup("client-id"))
	viper.BindPFlag("client_secret", rootCmd.Flags().Lookup("client-secret"))
	viper.BindPFlag("table", rootCmd.Flags().Lookup("table"))
	viper.BindPFlag("request_id", rootCmd.Flags().Lookup("request-id"))
	viper.BindPFlag("message", rootCmd.Flags().Lookup("message"))

	// Bind environment variables.
	viper.BindEnv("endpoint", "ZEROBUS_ENDPOINT")
	viper.BindEnv("host", "DATABRICKS_HOST")
	viper.BindEnv("client_id", "DATABRICKS_CLIENT_ID")
	viper.BindEnv("client_secret", "DATABRICKS_CLIENT_SECRET")
	viper.BindEnv("table", "TABLE_NAME")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
