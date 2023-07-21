/*
 */
package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/senzing/go-common/engineconfigurationjsonparser"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-common/option"
	"github.com/senzing/go-common/option/optiontype"
	"github.com/senzing/init-database/initializer"
	"github.com/senzing/senzing-tools/cmdhelper"
	"github.com/senzing/senzing-tools/constant"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envarEngineConfigurationFile        = "SENZING_TOOLS_ENGINE_CONFIGURATION_FILE"
	envarSqlFile                 string = "SENZING_TOOLS_SQL_FILE"
	Short                        string = "Initialize a database with the Senzing schema and configuration"
	Use                          string = "init-database"
)

var (
	Long string = getLong()
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

var OptionSqlFile = option.ContextVariable{
	Arg:     "sql-file",
	Default: getSqlFileDefault(),
	Envar:   envarSqlFile,
	Help:    "Path to file of SQL used to create Senzing database schema [%s]",
	Type:    optiontype.String,
}

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	option.Configuration,
	option.DatabaseUrl,
	option.Datasources,
	option.EngineConfigurationJson,
	option.EngineLogLevel,
	option.EngineModuleName,
	option.LicenseStringBase64,
	option.LogLevel,
	option.ObserverOrigin,
	option.ObserverUrl,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Construct the JSON string for the Senzing engine configuration.
func buildSenzingEngineConfigurationJson(ctx context.Context, aViper *viper.Viper) (string, error) {
	var err error = nil
	var result string = ""
	result = aViper.GetString(option.EngineConfigurationJson.Arg)
	if len(result) == 0 {
		options := map[string]string{
			"configPath":          aViper.GetString(option.ConfigPath.Arg),
			"databaseUrl":         aViper.GetString(option.DatabaseUrl.Arg),
			"licenseStringBase64": aViper.GetString(option.LicenseStringBase64.Arg),
			"resourcePath":        aViper.GetString(option.ResourcePath.Arg),
			"senzingDirectory":    aViper.GetString(option.SenzingDirectory.Arg),
			"supportPath":         aViper.GetString(option.SupportPath.Arg),
		}
		result, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(options)
		if err != nil {
			return result, err
		}
	}
	err = g2engineconfigurationjson.VerifySenzingEngineConfigurationJson(ctx, result)
	if err != nil {
		return result, err
	}
	return result, err
}

// Construct the path to the "g2config.json" file.
func getEngineConfigurationFileDefault() string {
	var result string = ""
	ctx := context.Background()

	// Early exit.  Environment variable is set.

	result, isSet := os.LookupEnv(envarEngineConfigurationFile)
	if isSet {
		return result
	}

	// Find information from SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.

	parsedSenzingEngineConfigurationJson, err := getParsedEngineConfigurationJson()
	if err != nil {
		return result
	}
	resourcePath, err := parsedSenzingEngineConfigurationJson.GetResourcePath(ctx)
	if err != nil {
		return result
	}
	result = resourcePath + "/templates/g2config.json"
	return result
}

// Create the value for the "Long" variable.
func getLong() string {
	var result string = `
Initialize a database with the Senzing schema and configuration.
For more information, visit https://github.com/Senzing/init-database
	`

	sqlFileDefault := getSqlFileDefault()
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
func getParsedEngineConfigurationJson() (engineconfigurationjsonparser.EngineConfigurationJsonParser, error) {
	var result engineconfigurationjsonparser.EngineConfigurationJsonParser = nil
	ctx := context.Background()

	// Early exit.  Environment variable is set.

	senzingEngineConfigurationJson, isSet := os.LookupEnv(option.EngineConfigurationJson.Arg)
	if isSet {
		return engineconfigurationjsonparser.New(senzingEngineConfigurationJson)
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

	senzingEngineConfigurationJson, err := buildSenzingEngineConfigurationJson(ctx, myViper)
	if err != nil {
		return result, err
	}
	return engineconfigurationjsonparser.New(senzingEngineConfigurationJson)
}

// Get the path to the SQL file used to create the Senzing database schema.
func getSqlFileDefault() string {
	var result string = ""
	ctx := context.Background()

	// Early exit.  Environment variable is set.

	result, isSet := os.LookupEnv(envarSqlFile)
	if isSet {
		return result
	}

	// Find information from SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.

	parsedSenzingEngineConfigurationJson, err := getParsedEngineConfigurationJson()
	if err != nil {
		return result
	}
	resourcePath, err := parsedSenzingEngineConfigurationJson.GetResourcePath(ctx)
	if err != nil {
		return result
	}
	databaseUrls, err := parsedSenzingEngineConfigurationJson.GetDatabaseUrls(ctx)
	if err != nil {
		return result
	}
	if len(databaseUrls) == 0 {
		return result
	}
	databaseUrl := databaseUrls[0]

	// Parse database URL to find which type of database is used.

	parsedUrl, err := url.Parse(databaseUrl)
	if err != nil {
		if strings.HasPrefix(databaseUrl, "postgresql") {
			index := strings.LastIndex(databaseUrl, ":")
			newDatabaseUrl := databaseUrl[:index] + "/" + databaseUrl[index+1:]
			parsedUrl, err = url.Parse(newDatabaseUrl)
		}
		if err != nil {
			return result
		}
	}

	// Based on database type, choose SQL file.

	switch parsedUrl.Scheme {
	case "sqlite3":
		result = resourcePath + "/schema/g2core-schema-sqlite-create.sql"
	case "postgresql":
		result = resourcePath + "/schema/g2core-schema-postgresql-create.sql"
	case "mysql":
		result = resourcePath + "/schema/g2core-schema-mysql-create.sql"
	case "mssql":
		result = resourcePath + "/schema/g2core-schema-mssql-create.sql"
	}

	return result
}

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, append(ContextVariables, OptionSqlFile, OptionEngineConfigurationFile))
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
	cmdhelper.PreRun(cobraCommand, args, Use, append(ContextVariables, OptionSqlFile, OptionEngineConfigurationFile))
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error = nil
	ctx := context.Background()

	senzingEngineConfigurationJson, err := buildSenzingEngineConfigurationJson(ctx, viper.GetViper())
	if err != nil {
		return err
	}

	initializer := &initializer.InitializerImpl{
		DataSources:                    viper.GetStringSlice(option.Datasources.Arg),
		ObserverOrigin:                 viper.GetString(option.ObserverOrigin.Arg),
		ObserverUrl:                    viper.GetString(option.ObserverUrl.Arg),
		SenzingEngineConfigurationFile: viper.GetString(OptionEngineConfigurationFile.Arg),
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		SenzingLogLevel:                viper.GetString(option.LogLevel.Arg),
		SenzingModuleName:              viper.GetString(option.EngineModuleName.Arg),
		SenzingVerboseLogging:          viper.GetInt(option.EngineLogLevel.Arg),
		SqlFile:                        viper.GetString(OptionSqlFile.Arg),
	}
	return initializer.Initialize(ctx)
}

// Used in construction of cobra.Command
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
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
