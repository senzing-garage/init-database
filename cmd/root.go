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
	configurationFile string
	buildVersion      string = "0.0.0"
	buildIteration    string = "0"
)

func makeVersion(version string, iteration string) string {
	var result string = ""
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
	Long:  `For more information, visit https://github.com/Senzing/initdatabase`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error = nil
		ctx := context.TODO()

		logLevel, ok := logger.TextToLevelMap[viper.GetString("log-level")]
		if !ok {
			logLevel = logger.LevelInfo
		}

		senzingEngineConfigurationJson := viper.GetString("engine-configuration-json")
		if len(senzingEngineConfigurationJson) == 0 {
			senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJson(viper.GetString("engine-configuration-json"))
			if err != nil {
				return err
			}
		}

		initializer := &initializer.InitializerImpl{
			DataSources:                    viper.GetStringSlice("datasources"),
			SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
			SenzingModuleName:              viper.GetString("senzing-module-name"),
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
	now := time.Now()
	cobra.OnInitialize(initConfig)

	// Define default values for input parameters.

	defaultDatabaseUrl := ""
	defaultDatasources := []string{}
	defaultEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		defaultEngineConfigurationJson = err.Error()
	}
	defaultEngineLogLevel := 0
	defaultEngineModuleName := fmt.Sprintf("initdatabase-%s", now.UTC())
	defaultLogLevel := "INFO"

	// Define flags for command.

	RootCmd.Flags().String("database-url", defaultDatabaseUrl, "URL of database to initialize [SENZING_TOOLS_DATABASE_URL]")
	RootCmd.Flags().StringSlice("datasources", defaultDatasources, "datasources to be added to initial Senzing configuration [SENZING_TOOLS_DATASOURCES]")
	RootCmd.Flags().String("engine-configuration-json", defaultEngineConfigurationJson, "JSON string sent to Senzing's init() function [SENZING_TOOLS_ENGINE_CONFIGURATION_JSON]")
	RootCmd.Flags().Int("engine-log-level", defaultEngineLogLevel, "log level for Senzing Engine [SENZING_TOOLS_ENGINE_LOG_LEVEL]")
	RootCmd.Flags().String("engine-module-name", defaultEngineModuleName, "the identifier given to the Senzing engine [SENZING_TOOLS_ENGINE_MODULE_NAME]")
	RootCmd.Flags().String("log-level", defaultLogLevel, "log level of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, or PANIC [SENZING_TOOLS_LOG_LEVEL]")

	// Integrate with Viper.

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("SENZING_TOOLS")

	// Define flags in Viper.

	viper.SetDefault("database-url", defaultDatabaseUrl)
	viper.BindPFlag("database-url", RootCmd.Flags().Lookup("database-url"))

	viper.SetDefault("datasources", defaultDatasources)
	viper.BindPFlag("datasources", RootCmd.Flags().Lookup("datasources"))

	viper.SetDefault("engine-configuration-json", defaultEngineConfigurationJson)
	viper.BindPFlag("engine-configuration-json", RootCmd.Flags().Lookup("engine-configuration-json"))

	viper.SetDefault("engine-log-level", defaultEngineLogLevel)
	viper.BindPFlag("engine-log-level", RootCmd.Flags().Lookup("engine-log-level"))

	viper.SetDefault("engine-module-name", defaultEngineModuleName)
	viper.BindPFlag("engine-module-name", RootCmd.Flags().Lookup("engine-module-name"))

	viper.SetDefault("log-level", defaultLogLevel)
	viper.BindPFlag("log-level", RootCmd.Flags().Lookup("log-level"))

	// Set version template.

	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	RootCmd.SetVersionTemplate(versionTemplate)
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
