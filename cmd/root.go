/*
 */
package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/senzing-garage/go-cmdhelping/cmdhelper"
	"github.com/senzing-garage/go-cmdhelping/constant"
	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-cmdhelping/option/optiontype"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-helpers/settingsparser"
	"github.com/senzing-garage/init-database/initializer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envarEngineConfigurationFile        = "SENZING_TOOLS_ENGINE_CONFIGURATION_FILE"
	envarSQLFile                 string = "SENZING_TOOLS_SQL_FILE"
	Short                        string = "Initialize a database with the Senzing schema and configuration"
	Use                          string = "init-database"
)

var (
	Long = getLong()
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var OptionEngineConfigurationFile = option.ContextVariable{
	Arg:     "engine-configuration-file",
	Default: getEngineConfigurationFileDefault(),
	Envar:   envarEngineConfigurationFile,
	Help:    "Path to file of JSON used to configure Senzing engine [%s]",
	Type:    optiontype.String,
}

var OptionSQLFile = option.ContextVariable{
	Arg:     "sql-file",
	Default: getSQLFileDefault(),
	Envar:   envarSQLFile,
	Help:    "Path to file of SQL used to create Senzing database schema [%s]",
	Type:    optiontype.String,
}

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	option.Configuration,
	option.DatabaseURL,
	option.Datasources,
	option.EngineSettings,
	option.EngineLogLevel,
	option.EngineInstanceName,
	option.LicenseStringBase64,
	option.LogLevel,
	option.ObserverOrigin,
	option.ObserverURL,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

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
	cmdhelper.PreRun(cobraCommand, args, Use, append(ContextVariables, OptionSQLFile, OptionEngineConfigurationFile))
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error
	ctx := context.Background()

	senzingSettings, err := buildSenzingEngineConfigurationJSON(ctx, viper.GetViper())
	if err != nil {
		return err
	}

	initializer := &initializer.BasicInitializer{
		DataSources:           viper.GetStringSlice(option.Datasources.Arg),
		ObserverOrigin:        viper.GetString(option.ObserverOrigin.Arg),
		ObserverURL:           viper.GetString(option.ObserverURL.Arg),
		SenzingInstanceName:   viper.GetString(option.EngineInstanceName.Arg),
		SenzingLogLevel:       viper.GetString(option.LogLevel.Arg),
		SenzingSettings:       senzingSettings,
		SenzingSettingsFile:   viper.GetString(OptionEngineConfigurationFile.Arg),
		SenzingVerboseLogging: viper.GetInt64(option.EngineLogLevel.Arg),
		SQLFile:               viper.GetString(OptionSQLFile.Arg),
	}
	return initializer.Initialize(ctx)
}

// Used in construction of cobra.Command
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Construct the JSON string for the Senzing engine configuration.
func buildSenzingEngineConfigurationJSON(ctx context.Context, aViper *viper.Viper) (string, error) {
	var err error
	var result string
	result = aViper.GetString(option.EngineSettings.Arg)
	if len(result) == 0 {
		options := map[string]string{
			"configPath":          aViper.GetString(option.ConfigPath.Arg),
			"databaseURL":         aViper.GetString(option.DatabaseURL.Arg),
			"licenseStringBase64": aViper.GetString(option.LicenseStringBase64.Arg),
			"resourcePath":        aViper.GetString(option.ResourcePath.Arg),
			"senzingDirectory":    aViper.GetString(option.SenzingDirectory.Arg),
			"supportPath":         aViper.GetString(option.SupportPath.Arg),
		}
		result, err = settings.BuildSimpleSettingsUsingMap(options)
		if err != nil {
			return result, err
		}
	}

	err = settings.VerifySettings(ctx, result)
	if err != nil {
		return result, err
	}
	return result, err
}

// Construct the path to the "g2config.json" file.
func getEngineConfigurationFileDefault() string {
	var result string
	ctx := context.Background()

	// Early exit.  Environment variable is set.

	result, isSet := os.LookupEnv(envarEngineConfigurationFile)
	if isSet {
		return result
	}

	// Find information from SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.

	parsedSenzingEngineConfigurationJSON, err := getParsedEngineConfigurationJSON()
	if err != nil {
		return result
	}
	resourcePath, err := parsedSenzingEngineConfigurationJSON.GetResourcePath(ctx)
	if err != nil {
		return result
	}
	result = resourcePath + "/templates/g2config.json"
	return result
}

// Create the value for the "Long" variable.
func getLong() string {
	var result = `
Initialize a database with the Senzing schema and configuration.
For more information, visit https://github.com/senzing-garage/init-database
	`

	sqlFileDefault := getSQLFileDefault()
	if len(sqlFileDefault) > 0 {
		result = fmt.Sprintf("%s\nThe SQL file used to create the Senzing database schema will be %s", result, sqlFileDefault)
	}
	engineConfigurationFileDefault := getEngineConfigurationFileDefault()
	if len(engineConfigurationFileDefault) > 0 {
		result = fmt.Sprintf("%s\nThe JSON file used to create the Senzing configuration  will be %s", result, engineConfigurationFileDefault)
	}
	return result
}

// Create a temporary parsed Senzing engine configuration.
func getParsedEngineConfigurationJSON() (settingsparser.SettingsParser, error) {
	var result settingsparser.SettingsParser
	ctx := context.Background()

	// Early exit.  Environment variable is set.

	senzingSettings, isSet := os.LookupEnv(option.EngineSettings.Arg)
	if isSet {
		return settingsparser.New(senzingSettings)
	}

	// Create a local Viper.

	myViper := viper.New()
	myViper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	myViper.SetEnvKeyReplacer(replacer)
	myViper.SetEnvPrefix(constant.SetEnvPrefix)

	for _, contextVariable := range ContextVariables {
		if contextVariable.Type == optiontype.String {
			myViper.SetDefault(contextVariable.Arg, contextVariable.Default)
		}
	}

	// Build and parse Senzing engine configuration JSON.

	senzingSettings, err := buildSenzingEngineConfigurationJSON(ctx, myViper)
	if err != nil {
		return result, err
	}
	return settingsparser.New(senzingSettings)
}

// Get the path to the SQL file used to create the Senzing database schema.
func getSQLFileDefault() string {
	var result string
	ctx := context.Background()

	// Early exit.  Environment variable is set.

	result, isSet := os.LookupEnv(envarSQLFile)
	if isSet {
		return result
	}

	// Find information from SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.

	parsedSenzingEngineConfigurationJSON, err := getParsedEngineConfigurationJSON()
	if err != nil {
		return result
	}
	resourcePath, err := parsedSenzingEngineConfigurationJSON.GetResourcePath(ctx)
	if err != nil {
		return result
	}
	databaseURLs, err := parsedSenzingEngineConfigurationJSON.GetDatabaseURLs(ctx)
	if err != nil {
		return result
	}
	if len(databaseURLs) == 0 {
		return result
	}
	databaseURL := databaseURLs[0]

	// Parse database URL to find which type of database is used.

	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		if strings.HasPrefix(databaseURL, "postgresql") {
			index := strings.LastIndex(databaseURL, ":")
			newDatabaseURL := databaseURL[:index] + "/" + databaseURL[index+1:]
			parsedURL, err = url.Parse(newDatabaseURL)
		}
		if err != nil {
			return result
		}
	}

	// Based on database type, choose SQL file.

	switch parsedURL.Scheme {
	case "sqlite3":
		result = resourcePath + "/schema/szcore-schema-sqlite-create.sql"
	case "postgresql":
		result = resourcePath + "/schema/szcore-schema-postgresql-create.sql"
	case "mysql":
		result = resourcePath + "/schema/szcore-schema-mysql-create.sql"
	case "mssql":
		result = resourcePath + "/schema/szcore-schema-mssql-create.sql"
	}

	return result
}

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, append(ContextVariables, OptionSQLFile, OptionEngineConfigurationFile))
}
