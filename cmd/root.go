/*
 */
package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/init-database/initializer"
	"github.com/senzing/senzing-tools/cmdhelper"
	"github.com/senzing/senzing-tools/envar"
	"github.com/senzing/senzing-tools/help"
	"github.com/senzing/senzing-tools/option"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Short string = "Initialize a database with the Senzing schema and configuration"
	Use   string = "init-database"
	Long  string = `
Initialize a database with the Senzing schema and configuration.
For more information, visit https://github.com/Senzing/init-database
    `
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var OptionSqlFile = cmdhelper.ContextString{
	Default: "bob",
	Envar:   "SENZING_TOOLS_SQL_FILE",
	Help:    "Path to file of SQL to process [%s]",
	Option:  "sql-file",
}

var OptionDatabaseUrl = cmdhelper.ContextString{
	Default: cmdhelper.OsLookupEnvString(envar.DatabaseUrl, "Jane"),
	Envar:   envar.DatabaseUrl,
	Help:    help.DatabaseUrl,
	Option:  option.DatabaseUrl,
}

var ContextInts = []cmdhelper.ContextInt{
	{
		Default: cmdhelper.OsLookupEnvInt(envar.EngineLogLevel, 0),
		Envar:   envar.EngineLogLevel,
		Help:    help.EngineLogLevel,
		Option:  option.EngineLogLevel,
	},
}

var ContextStrings = []cmdhelper.ContextString{
	{
		Default: cmdhelper.OsLookupEnvString(envar.Configuration, ""),
		Envar:   envar.Configuration,
		Help:    help.Configuration,
		Option:  option.Configuration,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.EngineConfigurationJson, ""),
		Envar:   envar.EngineConfigurationJson,
		Help:    help.EngineConfigurationJson,
		Option:  option.EngineConfigurationJson,
	},
	{
		Default: fmt.Sprintf("init-database-%d", time.Now().Unix()),
		Envar:   envar.EngineModuleName,
		Help:    help.EngineModuleName,
		Option:  option.EngineModuleName,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.LogLevel, "INFO"),
		Envar:   envar.LogLevel,
		Help:    help.LogLevel,
		Option:  option.LogLevel,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.ObserverOrigin, ""),
		Envar:   envar.ObserverOrigin,
		Help:    help.ObserverOrigin,
		Option:  option.ObserverOrigin,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.ObserverUrl, ""),
		Envar:   envar.ObserverUrl,
		Help:    help.ObserverUrl,
		Option:  option.ObserverUrl,
	},
}

var ContextStringSlices = []cmdhelper.ContextStringSlice{
	{
		Default: []string{},
		Envar:   envar.Datasources,
		Help:    help.Datasources,
		Option:  option.Datasources,
	},
}

var ContextVariables = &cmdhelper.ContextVariables{
	Ints:         ContextInts,
	Strings:      ContextStrings,
	StringSlices: ContextStringSlices,
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func viperizeString(cobraCommand *cobra.Command, option cmdhelper.ContextString) error {
	cobraCommand.Flags().String(option.Option, option.Default, fmt.Sprintf(option.Help, option.Envar))
	viper.SetDefault(option.Option, option.Default)
	err := viper.BindPFlag(option.Option, cobraCommand.Flags().Lookup(option.Option))
	return err
}

func findSqlFile(resourcePath string, databaseUrl string) string {
	var result string = ""
	// Determine which SQL file to process.

	parsedUrl, err := url.Parse(databaseUrl)
	if err != nil {
		if strings.HasPrefix(databaseUrl, "postgresql") {
			index := strings.LastIndex(databaseUrl, ":")
			newDatabaseUrl := databaseUrl[:index] + "/" + databaseUrl[index+1:]
			parsedUrl, err = url.Parse(newDatabaseUrl)
		}
		if err != nil {
			return ""
		}
	}

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
	cmdhelper.Init(RootCmd, *ContextVariables)

	fmt.Printf(">>>> 1.0 >>>> %s\n", viper.GetString(option.DatabaseUrl))

	// Tailor the "sql-file" option.

	err := viperizeString(RootCmd, OptionDatabaseUrl)
	if err != nil {
		panic(nil)
	}

	OptionSqlFile.Default = findSqlFile("bob", viper.GetString(OptionDatabaseUrl.Option))
	err = viperizeString(RootCmd, OptionSqlFile)
	if err != nil {
		panic(nil)
	}

}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	fmt.Printf(">>>> 2.0 >>>> %s\n", viper.GetString(option.DatabaseUrl))
	fmt.Printf(">>>> 2.1 >>>> %+v\n", viper.AllKeys())

	err := RootCmd.Execute()
	fmt.Printf(">>>> 2.2 >>>> %+v\n", viper.AllKeys())

	if err != nil {
		os.Exit(1)
	}
	fmt.Printf(">>>> 2.3 >>>> %+v\n", viper.AllKeys())

}

// Used in construction of cobra.Command
func PreRun(cobraCommand *cobra.Command, args []string) {
	fmt.Printf(">>>> 3.0 >>>> %s\n", viper.GetString(option.DatabaseUrl))

	cmdhelper.PreRun(cobraCommand, args, Use, *ContextVariables)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error = nil
	ctx := context.Background()

	fmt.Printf(">>>> 4.0 >>>> %s\n", viper.GetString(option.DatabaseUrl))

	// Build senzingEngineConfigurationJson.

	senzingEngineConfigurationJson := viper.GetString(option.EngineConfigurationJson)
	if len(senzingEngineConfigurationJson) == 0 {
		options := map[string]string{
			"configPath":          viper.GetString(option.ConfigPath),
			"databaseUrl":         viper.GetString(option.DatabaseUrl),
			"licenseStringBase64": viper.GetString(option.LicenseStringBase64),
			"resourcePath":        viper.GetString(option.ResourcePath),
			"senzingDirectory":    viper.GetString(option.SenzingDirectory),
			"supportPath":         viper.GetString(option.SupportPath),
		}
		senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(options)
		if err != nil {
			return err
		}
	}
	err = g2engineconfigurationjson.VerifySenzingEngineConfigurationJson(ctx, senzingEngineConfigurationJson)
	if err != nil {
		return err
	}

	initializer := &initializer.InitializerImpl{
		DataSources:                    viper.GetStringSlice(option.Datasources),
		ObserverOrigin:                 viper.GetString(option.ObserverOrigin),
		ObserverUrl:                    viper.GetString(option.ObserverUrl),
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		SenzingLogLevel:                viper.GetString(option.LogLevel),
		SenzingModuleName:              viper.GetString(option.EngineModuleName),
		SenzingVerboseLogging:          viper.GetInt(option.EngineLogLevel),
		SqlFile:                        viper.GetString(OptionSqlFile.Option),
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
