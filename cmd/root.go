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
	"github.com/senzing/init-database/initializer"
	"github.com/senzing/senzing-tools/constant"
	"github.com/senzing/senzing-tools/envar"
	"github.com/senzing/senzing-tools/helper"
	"github.com/senzing/senzing-tools/option"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfiguration           string = ""
	defaultDatabaseUrl             string = ""
	defaultEngineConfigurationJson string = ""
	defaultEngineLogLevel          int    = 0
	defaultLogLevel                string = "INFO"
	defaultObserverUrl             string = ""
	defaultObserverOrigin          string = ""
	Short                          string = "Initialize a database with the Senzing schema and configuration"
	Use                            string = "init-database"
	Long                           string = `
Initialize a database with the Senzing schema and configuration.
For more information, visit https://github.com/Senzing/init-database
	`
)

var (
	defaultDatasources      []string
	defaultEngineModuleName string = fmt.Sprintf("init-database-%d", time.Now().Unix())
)

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	RootCmd.Flags().Int(option.EngineLogLevel, defaultEngineLogLevel, fmt.Sprintf("Log level for Senzing Engine [%s]", envar.EngineLogLevel))
	RootCmd.Flags().String(option.Configuration, defaultConfiguration, fmt.Sprintf("Path to configuration file [%s]", envar.Configuration))
	RootCmd.Flags().String(option.DatabaseUrl, defaultDatabaseUrl, fmt.Sprintf("URL of database to initialize [%s]", envar.DatabaseUrl))
	RootCmd.Flags().String(option.EngineConfigurationJson, defaultEngineConfigurationJson, fmt.Sprintf("JSON string sent to Senzing's init() function [%s]", envar.EngineConfigurationJson))
	RootCmd.Flags().String(option.EngineModuleName, defaultEngineModuleName, fmt.Sprintf("Identifier given to the Senzing engine [%s]", envar.EngineModuleName))
	RootCmd.Flags().String(option.LogLevel, defaultLogLevel, fmt.Sprintf("Log level of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, or PANIC [%s]", envar.LogLevel))
	RootCmd.Flags().StringSlice(option.Datasources, defaultDatasources, fmt.Sprintf("Datasources to be added to initial Senzing configuration [%s]", envar.Datasources))
	RootCmd.Flags().String("observer-url", defaultObserverUrl, fmt.Sprintf("URL of Observer [%s]", "SENZING_TOOLS_OBSERVER_URL"))                                 // FIXME: use "option." and "envar." when available.
	RootCmd.Flags().String("observer-origin", defaultObserverOrigin, fmt.Sprintf("Identify this instance to the Observer [%s]", "SENZING_TOOLS_OBSERVER_ORIGIN")) // FIXME: use "option." and "envar." when available.

	// TODO: Update when available.
	// RootCmd.Flags().String(option.ObserverUrl, defaultObserverUrl, fmt.Sprintf("URL of Observer [%s]", envar.ObserverUrl"))
	// RootCmd.Flags().String(option.ObserverOrigin, defaultObserverOrigin, fmt.Sprintf("Identify this instance to the Observer [%s]", envar.ObserverOrigin))

}

// If a configuration file is present, load it.
func loadConfigurationFile(cobraCommand *cobra.Command) {
	configuration := ""
	configFlag := cobraCommand.Flags().Lookup(option.Configuration)
	if configFlag != nil {
		configuration = configFlag.Value.String()
	}
	if configuration != "" { // Use configuration file specified as a command line option.
		viper.SetConfigFile(configuration)
	} else { // Search for a configuration file.

		// Determine home directory.

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Specify configuration file name.

		viper.SetConfigName("init-database")
		viper.SetConfigType("yaml")

		// Define search path order.

		viper.AddConfigPath(home + "/.senzing-tools")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/senzing-tools")
	}

	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Applying configuration file:", viper.ConfigFileUsed())
	}
}

// Configure Viper with user-specified options.
func loadOptions(cobraCommand *cobra.Command) {
	var err error = nil
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix(constant.SetEnvPrefix)

	// Ints

	intOptions := map[string]int{
		option.EngineLogLevel: defaultEngineLogLevel,
	}
	for optionKey, optionValue := range intOptions {
		viper.SetDefault(optionKey, optionValue)
		err = viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
		if err != nil {
			panic(err)
		}
	}

	// Strings

	stringOptions := map[string]string{
		option.DatabaseUrl:             defaultDatabaseUrl,
		option.EngineConfigurationJson: defaultEngineConfigurationJson,
		option.EngineModuleName:        defaultEngineModuleName,
		option.LogLevel:                defaultLogLevel,
		"observer-url":                 defaultObserverUrl,
		"observer-origin":              defaultObserverOrigin,

		// TODO: Update when available.
		// option.ObserverUrl:                 defaultObserverUrl,
		// option.ObserverOrigin:              defaultObserverOrigin,
	}
	for optionKey, optionValue := range stringOptions {
		viper.SetDefault(optionKey, optionValue)
		err = viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
		if err != nil {
			panic(err)
		}
	}

	// StringSlice

	viper.SetDefault(option.Datasources, defaultDatasources)
	err = viper.BindPFlag(option.Datasources, cobraCommand.Flags().Lookup(option.Datasources))
	if err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Used in construction of cobra.Command
func PreRun(cobraCommand *cobra.Command, args []string) {
	loadConfigurationFile(cobraCommand)
	loadOptions(cobraCommand)
	cobraCommand.SetVersionTemplate(constant.VersionTemplate)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error = nil
	ctx := context.TODO()

	senzingEngineConfigurationJson := viper.GetString(option.EngineConfigurationJson)
	if len(senzingEngineConfigurationJson) == 0 {
		senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJson(viper.GetString(option.DatabaseUrl))
		if err != nil {
			return err
		}
	}

	initializer := &initializer.InitializerImpl{
		DataSources:                    viper.GetStringSlice(option.Datasources),
		ObserverOrigin:                 viper.GetString("observer-origin"),
		ObserverUrl:                    viper.GetString("observer-url"),
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		SenzingLogLevel:                viper.GetString(option.LogLevel),
		SenzingModuleName:              viper.GetString(option.EngineModuleName),
		SenzingVerboseLogging:          viper.GetInt(option.EngineLogLevel),
		// TODO: Update when available.
		// ObserverOrigin:                 viper.GetString(option.ObserverOrigin),
		// ObserverUrl:                    viper.GetString(option.ObserverUrl),
	}

	err = initializer.Initialize(ctx)
	return err
}

// Used in construction of cobra.Command
func Version() string {
	return helper.MakeVersion(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Command
// ----------------------------------------------------------------------------

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:     Use,
	Short:   Short,
	Long:    Long,
	PreRun:  PreRun,
	RunE:    RunE,
	Version: Version(),
}
