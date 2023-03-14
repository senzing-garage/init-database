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
	"github.com/senzing/senzing-tools/constant"
	"github.com/senzing/senzing-tools/envar"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfiguration           string = ""
	defaultDatabaseUrl             string = ""
	defaultEngineConfigurationJson string = ""
	defaultEngineLogLevel          int    = 0
	defaultLogLevel                string = "INFO"
)

var (
	buildIteration          string = "0"
	buildVersion            string = "0.1.4"
	defaultDatasources      []string
	defaultEngineModuleName string = fmt.Sprintf("initdatabase-%d", time.Now().Unix())
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

func loadConfigurationFile(cobraCommand *cobra.Command) {
	configuration := cobraCommand.Flags().Lookup(constant.Configuration).Value.String()

	if configuration != "" { // Use configuration file specified as a command line option.
		viper.SetConfigFile(configuration)
	} else { // Search for a configuration file.

		// Determine home directory.

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Specify configuration file name.

		viper.SetConfigName("initdatabase")
		viper.SetConfigType("yaml")

		// Define search path order.

		viper.AddConfigPath(home + "/.senzing-tools")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/senzing-tools")
	}

	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func loadOptions(cobraCommand *cobra.Command) {

	// Integrate with Viper.

	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix(constant.SetEnvPrefix)

	// Define flags in Viper.

	viper.SetDefault(constant.DatabaseUrl, defaultDatabaseUrl)
	viper.BindPFlag(constant.DatabaseUrl, cobraCommand.Flags().Lookup(constant.DatabaseUrl))

	viper.SetDefault(constant.Datasources, defaultDatasources)
	viper.BindPFlag(constant.Datasources, cobraCommand.Flags().Lookup(constant.Datasources))

	viper.SetDefault(constant.EngineConfigurationJson, defaultEngineConfigurationJson)
	viper.BindPFlag(constant.EngineConfigurationJson, cobraCommand.Flags().Lookup(constant.EngineConfigurationJson))

	viper.SetDefault(constant.EngineLogLevel, defaultEngineLogLevel)
	viper.BindPFlag(constant.EngineLogLevel, cobraCommand.Flags().Lookup(constant.EngineLogLevel))

	viper.SetDefault(constant.EngineModuleName, defaultEngineModuleName)
	viper.BindPFlag(constant.EngineModuleName, cobraCommand.Flags().Lookup(constant.EngineModuleName))

	viper.SetDefault(constant.LogLevel, defaultLogLevel)
	viper.BindPFlag(constant.LogLevel, cobraCommand.Flags().Lookup(constant.LogLevel))
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
		fmt.Println(">>>>> initdatabase.PreRun")
		loadConfigurationFile(cobraCommand)
		loadOptions(cobraCommand)
		// Set version template.

		// versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
		cobraCommand.SetVersionTemplate(constant.VersionTemplate)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error = nil
		ctx := context.TODO()

		logLevel, ok := logger.TextToLevelMap[viper.GetString(constant.LogLevel)]
		if !ok {
			logLevel = logger.LevelInfo
		}

		senzingEngineConfigurationJson := viper.GetString(constant.EngineConfigurationJson)
		if len(senzingEngineConfigurationJson) == 0 {
			senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJson(viper.GetString(constant.DatabaseUrl))
			if err != nil {
				return err
			}
		}

		initializer := &initializer.InitializerImpl{
			DataSources:                    viper.GetStringSlice(constant.Datasources),
			SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
			SenzingModuleName:              viper.GetString(constant.EngineModuleName),
			SenzingVerboseLogging:          viper.GetInt(constant.EngineLogLevel),
		}
		err = initializer.SetLogLevel(ctx, logLevel)
		if err != nil {
			return err
		}
		err = initializer.Initialize(ctx)
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
	fmt.Println(">>>>> initdatabase.init()")

	RootCmd.Flags().Int(constant.EngineLogLevel, defaultEngineLogLevel, fmt.Sprintf("Log level for Senzing Engine [%s]", envar.EngineLogLevel))
	RootCmd.Flags().String(constant.Configuration, defaultConfiguration, fmt.Sprintf("Path to configuration file [%s]", envar.Configuration))
	RootCmd.Flags().String(constant.DatabaseUrl, defaultDatabaseUrl, fmt.Sprintf("URL of database to initialize [%s]", envar.DatabaseUrl))
	RootCmd.Flags().String(constant.EngineConfigurationJson, defaultEngineConfigurationJson, fmt.Sprintf("JSON string sent to Senzing's init() function [%s]", envar.EngineConfigurationJson))
	RootCmd.Flags().String(constant.EngineModuleName, defaultEngineModuleName, fmt.Sprintf("Identifier given to the Senzing engine [%s]", envar.EngineModuleName))
	RootCmd.Flags().String(constant.LogLevel, defaultLogLevel, fmt.Sprintf("Log level of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, or PANIC [%s]", envar.LogLevel))
	RootCmd.Flags().StringSlice(constant.Datasources, defaultDatasources, fmt.Sprintf("Datasources to be added to initial Senzing configuration [%s]", envar.Datasources))
}
