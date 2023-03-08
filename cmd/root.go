/*
 */
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/initdatabase/initializer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	buildIteration                 string = "0"
	buildVersion                   string = "0.1.4"
	configurationFile              string
	defaultDatabaseUrl             string = ""
	defaultDatasources             []string
	defaultEngineConfigurationJson string = ""
	defaultEngineLogLevel          int    = 0
	defaultEngineModuleName        string = fmt.Sprintf("initdatabase-%d", time.Now().Unix())
	defaultLogLevel                string = "INFO"
)

func makeVersion(version string, iteration string) string {
	result := ""
	if buildIteration == "0" {
		result = version
	} else {
		result = fmt.Sprintf("%s-%s", buildVersion, buildIteration)
	}
	return result
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "initdatabase",
	Short: "Initialize a database with the Senzing schema and configuration",
	Long: `
Initialize a database with the Senzing schema and configuration.
For more information, visit https://github.com/Senzing/initdatabase
	`,
	PreRun: func(cobraCommand *cobra.Command, args []string) {

		// Integrate with Viper.

		replacer := strings.NewReplacer("-", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.SetEnvPrefix("SENZING_TOOLS")

		// Define flags in Viper.

		viper.SetDefault("engine-log-level", defaultEngineLogLevel)
		viper.BindPFlag("engine-log-level", cobraCommand.Flags().Lookup("engine-log-level"))

		viper.SetDefault("database-url", defaultDatabaseUrl)
		viper.BindPFlag("database-url", cobraCommand.Flags().Lookup("database-url"))

		viper.SetDefault("engine-configuration-json", defaultEngineConfigurationJson)
		viper.BindPFlag("engine-configuration-json", cobraCommand.Flags().Lookup("engine-configuration-json"))

		viper.SetDefault("engine-module-name", defaultEngineModuleName)
		viper.BindPFlag("engine-module-name", cobraCommand.Flags().Lookup("engine-module-name"))

		viper.SetDefault("log-level", defaultLogLevel)
		viper.BindPFlag("log-level", cobraCommand.Flags().Lookup("log-level"))

		viper.SetDefault("datasources", defaultDatasources)
		viper.BindPFlag("datasources", cobraCommand.Flags().Lookup("datasources"))

		// Set version template.

		versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
		cobraCommand.SetVersionTemplate(versionTemplate)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error = nil
		ctx := context.TODO()

		logLevel, ok := logger.TextToLevelMap[viper.GetString("log-level")]
		if !ok {
			logLevel = logger.LevelInfo
		}

		senzingEngineConfigurationJson := viper.GetString("engine-configuration-json")
		if len(senzingEngineConfigurationJson) == 0 {
			senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJson(viper.GetString("database-url"))
			if err != nil {
				return err
			}
		}

		initializer := &initializer.InitializerImpl{
			DataSources:                    viper.GetStringSlice("datasources"),
			SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
			SenzingModuleName:              viper.GetString("engine-module-name"),
			SenzingVerboseLogging:          viper.GetInt("engine-log-level"),
		}
		initializer.SetLogLevel(ctx, logLevel)
		initializer.Initialize(ctx)
		return err
	},
	Version: makeVersion(buildVersion, buildIteration),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Define flags for Cobra command.

	RootCmd.Flags().Int("engine-log-level", defaultEngineLogLevel, "Log level for Senzing Engine [SENZING_TOOLS_ENGINE_LOG_LEVEL]")
	RootCmd.Flags().String("database-url", defaultDatabaseUrl, "URL of database to initialize [SENZING_TOOLS_DATABASE_URL]")
	RootCmd.Flags().String("engine-configuration-json", defaultEngineConfigurationJson, "JSON string sent to Senzing's init() function [SENZING_TOOLS_ENGINE_CONFIGURATION_JSON]")
	RootCmd.Flags().String("engine-module-name", defaultEngineModuleName, "Identifier given to the Senzing engine [SENZING_TOOLS_ENGINE_MODULE_NAME]")
	RootCmd.Flags().String("log-level", defaultLogLevel, "Log level of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, or PANIC [SENZING_TOOLS_LOG_LEVEL]")
	RootCmd.Flags().StringSlice("datasources", defaultDatasources, "Datasources to be added to initial Senzing configuration [SENZING_TOOLS_DATASOURCES]")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configurationFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configurationFile)
	} else {

		// Find home directory.

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".senzing-tools" (without extension).

		viper.AddConfigPath(home + "/.senzing-tools")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/senzing-tools")
		viper.SetConfigType("yaml")
		viper.SetConfigName("initdatabase")
	}

	// Read in environment variables that match "SENZING_TOOLS_*" pattern.

	viper.AutomaticEnv()

	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
