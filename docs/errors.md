# init-database errors

## Error prefixes

Error identifiers are in the format `senzing-PPPPnnnn` where:

`P` is a prefix used to identify the package.
`n` is a location within the package.

Prefixes:

1. `6501` - initializer
1. `6502` - senzingconfig
1. `6503` - senzingschema

## Common errors

### Postgresql

1. "Error: pq: SSL is not enabled on the server"
    1. The database URL needs the `sslmode` parameter.
       Example:

        ```console
        postgresql://username:password@postgres.example.com:5432/G2/?sslmode=disable
        ```

    1. [Connection String Parameters](https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters)

## Errors by ID

### senzing-65010010

- Trace the entering of initializer.Initialize().
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010011

- Trace the exiting of initializer.Initialize(); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010012

- Trace the exiting of initializer.Initialize(); initializerImpl.InitializeSpecificDatabase failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010013

- Trace the exiting of initializer.Initialize(); senzingSchema.SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010014

- Trace the exiting of initializer.Initialize(); senzingSchema.InitializeSenzing failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010015

- Trace the exiting of initializer.Initialize(); senzingConfig.SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010016

- Trace the exiting of initializer.Initialize(); senzingConfig.InitializeSenzing; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010017

- Trace the exiting of initializer.Initialize(); initializerImpl.observers.RegisterObserver; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010018

- Trace the exiting of initializer.Initialize(); initializerImpl.createGrpcObserver; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010019

- Trace the exiting of initializer.Initialize(); initializerImpl.registerObserverSenzingSchema; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010020

- Trace the exiting of initializer.Initialize(); initializerImpl.registerObserverSenzingConfig; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010029

- Trace the exiting of initializer.Initialize() returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010040

- Trace the entering of initializer.InitializeSpecificDatabase().
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010041

- Trace the exiting of initializer.InitializeSpecificDatabase(); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010042

- Trace the exiting of initializer.InitializeSpecificDatabase(); engineconfigurationjsonparser.New failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010043

- Trace the exiting of initializer.InitializeSpecificDatabase(); parser.GetDatabaseUrls failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010044

- Trace the exiting of initializer.InitializeSpecificDatabase(); url.Parse failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010045

- Trace the exiting of initializer.InitializeSpecificDatabase(); initializerImpl.initializeSpecificDatabaseSqlite; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010049

- Trace the exiting of initializer.InitializeSpecificDatabase() returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010050

- Trace the entering of initializer.RegisterObserver(%s).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010051

- Trace the exiting of initializer.RegisterObserver(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010052

- Trace the exiting of initializer.RegisterObserver(%s); initializerImpl.observers.RegisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010053

- Trace the exiting of initializer.RegisterObserver(%s); initializerImpl.getSenzingConfig().RegisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010054

- Trace the exiting of initializer.RegisterObserver(%s); initializerImpl.getSenzingSchema().RegisterObserver; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010059

- Trace the exiting of initializer.RegisterObserver(%s) returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010060

- Trace the entering of initializer.SetLogLevel(%s).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010061

- Trace the exiting of initializer.SetLogLevel(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010062

- Trace the exiting of initializer.SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010063

- Trace the exiting of initializer.SetLogLevel(%s); initializerImpl.getLogger().SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010064

- Trace the exiting of initializer.SetLogLevel(%s); initializerImpl.senzingConfigSingleton.SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010065

- Trace the exiting of initializer.SetLogLevel(%s); initializerImpl.getSenzingSchema().SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010069

- Trace the exiting of initializer.SetLogLevel(%s) returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010070

- Trace the entering of initializer.UnregisterObserver(%s).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010071

- Trace the exiting of initializer.UnregisterObserver(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010072

- Trace the exiting of initializer.UnregisterObserver(%s); initializerImpl.getSenzingConfig().UnregisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010073

- Trace the exiting of initializer.UnregisterObserver(%s); initializerImpl.getSenzingSchema().UnregisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010074

- Trace the exiting of initializer.UnregisterObserver(%s); initializerImpl.observers.UnregisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010079

- Trace the exiting of initializer.UnregisterObserver(%s) returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010080

- Trace the entering of initializer.SetObserverOrigin(%s).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010081

- Trace the exiting of initializer.SetObserverOrigin(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010089

- Trace the exiting of initializer.SetObserverOrigin(%s).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010100

- Trace the entering of initializer.initializeSpecificDatabaseSqlite(%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010101

- Trace the exiting of initializer.initializeSpecificDatabaseSqlite(%v); os.Stat failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010102

- Trace the exiting of initializer.initializeSpecificDatabaseSqlite(%v); os.MkdirAll failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010103

- Trace the exiting of initializer.initializeSpecificDatabaseSqlite(%v); os.Create failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65010109

- Trace the exiting of initializer.initializeSpecificDatabaseSqlite(%v) returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011000

- initializer.Initialize parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011001

- initializer.InitializeSpecificDatabase parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011002

- initializer.RegisterObserver parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011003

- initializer.SetLogLevel parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011004

- initializer.SetObserverOrigin parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011005

- initializer.UnregisterObserver parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011011

- initializer.Initialize(); json.Marshal failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011012

- initializer.Initialize(); initializerImpl.InitializeSpecificDatabase failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011013

- initializer.Initialize(); initializerImpl.getSenzingSchema failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011014

- initializer.Initialize(); senzingSchema.InitializeSenzing failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011015

- initializer.Initialize(); initializerImpl.getSenzingConfig failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011016

- initializer.Initialize(); senzingConfig.InitializeSenzing; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011017

- initializer.Initialize(); initializerImpl.observers.RegisterObserver; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011018

- initializer.Initialize(); initializerImpl.createGrpcObserver; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011041

- initializer.InitializeSpecificDatabase(); json.Marshal failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011042

- initializer.InitializeSpecificDatabase(); engineconfigurationjsonparser.New failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011043

- initializer.InitializeSpecificDatabase(); parser.GetDatabaseUrls failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011044

- initializer.InitializeSpecificDatabase(); url.Parse failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011045

- initializer.InitializeSpecificDatabase(); initializerImpl.initializeSpecificDatabaseSqlite; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011051

- initializer.RegisterObserver(%s); json.Marshal failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011052

- initializer.RegisterObserver(%s); initializerImpl.observers.RegisterObserver failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011053

- initializer.RegisterObserver(%s); initializerImpl.getSenzingConfig().RegisterObserver failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011054

- initializer.RegisterObserver(%s); initializerImpl.getSenzingSchema().RegisterObserver; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011061

- initializer.SetLogLevel(%s); json.Marshal failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011062

- initializer.SetLogLevel(%s); logging.IsValidLogLevelName failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011063

- initializer.SetLogLevel(%s); initializerImpl.getLogger().SetLogLevel failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011064

- initializer.SetLogLevel(%s); initializerImpl.senzingConfigSingleton.SetLogLevel failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011065

- initializer.SetLogLevel(%s); initializerImpl.getSenzingSchema().SetLogLevel failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011071

- initializer.UnregisterObserver(%s); json.Marshal failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011072

- initializer.UnregisterObserver(%s); initializerImpl.getSenzingConfig().UnregisterObserver failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011073

- initializer.UnregisterObserver(%s); initializerImpl.getSenzingSchema().UnregisterObserver failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011074

- initializer.UnregisterObserver(%s); initializerImpl.observers.UnregisterObserver failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011081

- initializer.SetObserverOrigin(%s); json.Marshal failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011101

- initializer.initializeSpecificDatabaseSqlite(%v); os.Stat failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011102

- initializer.initializeSpecificDatabaseSqlite(%v); os.MkdirAll failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65011103

- initializer.initializeSpecificDatabaseSqlite(%v); os.Create failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65012001

- Created file: %s
- See <https://github.com/Senzing/init-database/blob/main/initializer/initializer.go>

### senzing-65020010

- Trace the entering of senzingconfig.InitializeSenzing().
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020011

- Trace the exiting of senzingconfig.InitializeSenzing(); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020012

- Trace the exiting of senzingconfig.InitializeSenzing(); senzingConfig.getDependentServices failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020013

- Trace the exiting of senzingconfig.InitializeSenzing(); g2Configmgr.GetDefaultConfigID failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020014

- Trace the exiting of senzingconfig.InitializeSenzing(); Senzing configuration already exists; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020015

- Trace the exiting of senzingconfig.InitializeSenzing(); g2Config.Create failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020016

- Trace the exiting of senzingconfig.InitializeSenzing(); senzingConfig.addDatasources failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020017

- Trace the exiting of senzingconfig.InitializeSenzing(); g2Config.Save failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020018

- Trace the exiting of senzingconfig.InitializeSenzing(); g2Configmgr.AddConfig failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020019

- Trace the exiting of senzingconfig.InitializeSenzing(); g2Configmgr.SetDefaultConfigID failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020029

- Trace the exiting of senzingconfig.InitializeSenzing() returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020030

- Trace the entering of senzingconfig.RegisterObserver(%s).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020031

- Trace the exiting of senzingconfig.RegisterObserver(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020032

- Trace the exiting of senzingconfig.RegisterObserver(%s); senzingConfig.observers.RegisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020033

- Trace the exiting of senzingconfig.RegisterObserver(%s); senzingConfig.getDependentServices failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020034

- Trace the exiting of senzingconfig.RegisterObserver(%s); g2Config.RegisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020035

- Trace the exiting of senzingconfig.RegisterObserver(%s); g2Configmgr.RegisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020039

- Trace the exiting of senzingconfig.RegisterObserver(%s) returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020040

- Trace the entering of senzingconfig.SetLogLevel(%s).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020041

- Trace the exiting of senzingconfig.SetLogLevel(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020042

- Trace the exiting of senzingconfig.SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020043

- Trace the exiting of senzingconfig.SetLogLevel(%s); senzingConfig.getLogger().SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020044

- Trace the exiting of senzingconfig.SetLogLevel(%s); senzingConfig.getDependentServices failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020045

- Trace the exiting of senzingconfig.SetLogLevel(%s); g2Config.SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020046

- Trace the exiting of senzingconfig.SetLogLevel(%s); g2Configmgr.SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020049

- Trace the exiting of senzingconfig.SetLogLevel(%s) returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020050

- Trace the entering of senzingconfig.UnregisterObserver(%s).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020051

- Trace the exiting of senzingconfig.UnregisterObserver(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020052

- Trace the exiting of senzingconfig.UnregisterObserver(%s); g2Config.UnregisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020053

- Trace the exiting of senzingconfig.UnregisterObserver(%s); g2Configmgr.UnregisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020054

- Trace the exiting of senzingconfig.UnregisterObserver(%s); senzingConfig.observers.UnregisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020059

- Trace the exiting of senzingconfig.UnregisterObserver(%s) returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020060

- Trace the entering of senzingconfig.SetObserverOrigin(%s).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020061

- Trace the exiting of senzingconfig.SetObserverOrigin(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65020069

- Trace the exiting of senzingconfig.SetObserverOrigin(%s).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021001

- senzingconfig.InitializeSenzing parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021002

- senzingconfig.RegisterObserver parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021003

- senzingconfig.SetLogLevel parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021004

- senzingconfig.SetObserverOrigin parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021005

- senzingconfig.UnregisterObserver parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021011

- senzingconfig.Initialize(); json.Marshal failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021012

- senzingconfig.Initialize(); senzingConfig.getDependentServices failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021013

- senzingconfig.Initialize(); g2Configmgr.GetDefaultConfigID failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021014

- senzingconfig.Initialize(); senzingSchema.InitializeSenzing failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021015

- senzingconfig.Initialize(); g2Config.Create failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021016

- senzingconfig.Initialize(); senzingConfig.addDatasources failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021017

- senzingconfig.Initialize(); g2Config.Save failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021018

- senzingconfig.Initialize(); g2Configmgr.AddConfig failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021019

- senzingconfig.Initialize(); g2Configmgr.SetDefaultConfigID failed; Error: %v.
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021031

- senzingconfig.RegisterObserver(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021032

- senzingconfig.RegisterObserver(%s); senzingConfig.observers.RegisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021033

- senzingconfig.RegisterObserver(%s); senzingConfig.getDependentServices failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021034

- senzingconfig.RegisterObserver(%s); g2Config.RegisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021035

- senzingconfig.RegisterObserver(%s); g2Configmgr.RegisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021041

- senzingconfig.SetLogLevel(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021042

- senzingconfig.SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021043

- senzingconfig.SetLogLevel(%s); senzingConfig.getLogger().SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021044

- senzingconfig.SetLogLevel(%s); senzingConfig.getDependentServices failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021045

- senzingconfig.SetLogLevel(%s); g2Config.SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021046

- senzingconfig.SetLogLevel(%s); g2Configmgr.SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021051

- senzingconfig.UnregisterObserver(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021052

- senzingconfig.UnregisterObserver(%s); g2Config.UnregisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021053

- senzingconfig.UnregisterObserver(%s); g2Configmgr.UnregisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65021054

- senzingconfig.UnregisterObserver(%s); senzingConfig.observers.UnregisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65022001

- "Added Datasource: %s"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65022002

- "No new Senzing configuration created.  One already exists (%d).
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65022003

- "Created Senzing configuration: %d named: %s"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65028001

- senzingconfig.InitializeSenzing - config exists"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65028002

- senzingconfig.InitializeSenzing"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65028003

- senzingconfig.RegisterObserver"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65028004

- senzingconfig.SetLogLevel"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65028005

- senzingconfig.SetObserverOrigin"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65028006

- senzingconfig.UnregisterObserver"
- See <https://github.com/Senzing/init-database/blob/main/senzingconfig/senzingconfig.go>

### senzing-65030010

- Trace the entering of senzingschema.InitializeSenzing().
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030011

- Trace the exiting of senzingschema.InitializeSenzing(); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030012

- Trace the exiting of senzingschema.InitializeSenzing(); engineconfigurationjsonparser.New failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030013

- Trace the exiting of senzingschema.InitializeSenzing(); parser.GetResourcePath failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030014

- Trace the exiting of senzingschema.InitializeSenzing(); parser.GetDatabaseUrls failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030015

- Trace the exiting of senzingschema.InitializeSenzing(); senzingSchema.processDatabase failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030019

- Trace the exiting of senzingschema.InitializeSenzing() returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030020

- Trace the entering of senzingschema.RegisterObserver(%s).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030021

- Trace the exiting of senzingschema.RegisterObserver(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030022

- Trace the exiting of senzingschema.RegisterObserver(%s); senzingSchema.observers.RegisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030029

- Trace the exiting of senzingschema.RegisterObserver(%s) returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030030

- Trace the entering of senzingschema.SetLogLevel(%s).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030031

- Trace the exiting of senzingschema.SetLogLevel(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030032

- Trace the exiting of senzingschema.SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030033

- Trace the exiting of senzingschema.senzingSchema.getLogger().SetLogLevel(%s) failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030039

- Trace the exiting of senzingschema.SetLogLevel(%s) returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030040

- Trace the entering of senzingschema.UnregisterObserver(%s).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030041

- Trace the exiting of senzingschema.UnregisterObserver(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030042

- Trace the exiting of senzingschema.UnregisterObserver(%s); senzingSchema.observers.UnregisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030049

- Trace the exiting of senzingschema.UnregisterObserver(%s) returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030050

- Trace the entering of senzingschema.SetObserverOrigin(%s).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030051

- Trace the exiting of senzingschema.SetObserverOrigin(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030059

- Trace the exiting of senzingschema.SetObserverOrigin(%s).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030100

- Trace the entering of senzingschema.processDatabase(%s, %s).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030101

- Trace the exiting of senzingschema.processDatabase(%s, %s); url.Parse failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030102

- Trace the exiting of senzingschema.processDatabase(%s, %s); connector.NewConnector failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030103

- Trace the exiting of senzingschema.processDatabase(%s, %s); sqlExecutor.SetLogLevel failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030104

- Trace the exiting of senzingschema.processDatabase(%s, %s); sqlExecutor.RegisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030105

- Trace the exiting of senzingschema.processDatabase(%s, %s); sqlExecutor.ProcessFileName failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65030109

- Trace the exiting of senzingschema.processDatabase(%s, %s) returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031001

- senzingschema.InitializeSenzing parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031002

- senzingschema.RegisterObserver parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031003

- senzingschema.SetLogLevel parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031004

- senzingschema.SetObserverOrigin parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031005

- senzingschema.UnregisterObserver parameters: %+v"
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031011

- senzingschema.InitializeSenzing(); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031012

- senzingschema.InitializeSenzing(); engineconfigurationjsonparser.New failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031013

- senzingschema.InitializeSenzing(); parser.GetResourcePath failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031014

- senzingschema.InitializeSenzing(); parser.GetDatabaseUrls failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031015

- senzingschema.InitializeSenzing(); senzingSchema.processDatabase failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031021

- senzingschema.RegisterObserver(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031022

- senzingschema.RegisterObserver(%s); senzingSchema.observers.RegisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031031

- senzingschema.SetLogLevel(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031032

- senzingschema.SetLogLevel(%s); logging.IsValidLogLevelName failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031033

- senzingschema.senzingSchema.getLogger().SetLogLevel(%s) failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031041

- senzingschema.UnregisterObserver(%s); json.Marshal failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65031042

- senzingschema.UnregisterObserver(%s); senzingSchema.observers.UnregisterObserver failed; returned (%v).
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>

### senzing-65032001

- Sent SQL in %s to database %s"
- See <https://github.com/Senzing/init-database/blob/main/senzingschema/senzingschema.go>
