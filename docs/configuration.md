# Configuration

The configuration must be set in multiple YAML files located in `conf/` folder from the current working directory.

You can create multiple files containing different part of the configuration. A global merge will be done across all data in all files.

Moreover, the configuration files will be watched for modifications.

You can see a full example in the [Example section](#example)

## Main structure

| Key                    | Type                                                                        | Required | Default                                       | Description                                                                                                                                 |
| ---------------------- | --------------------------------------------------------------------------- | -------- | --------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| log                    | [LogConfiguration](#logconfiguration)                                       | No       | None                                          | Log configurations                                                                                                                          |
| tracing                | None                                                                        | No       | [TracingConfiguration](#tracingconfiguration) | Tracing configurations (Jaeger compatible)                                                                                                  |
| server                 | [ServerConfiguration](#serverconfiguration)                                 | No       | None                                          | Public Server configurations (used for API and UI)                                                                                          |
| internalServer         | [ServerConfiguration](#serverconfiguration)                                 | No       | None                                          | Internal Server configurations                                                                                                              |
| opaPublisherServer     | [ServerConfiguration](#serverconfiguration)                                 | No       | None                                          | OPA Publisher Server configurations (will be used for OPA servers publish)                                                                  |
| database               | [DatabaseConfiguration](#databaseconfiguration)                             | Yes      | None                                          | Database configurations                                                                                                                     |
| oidcAuthentication     | [OIDCAuthenticationConfiguration](#oidcauthenticationconfiguration)         | No       | None                                          | OIDC Authentication system configurations (Without this, no authentication will be done)                                                    |
| opaServerAuthorization | [OPAServerAuthorizationConfiguration](#opaserverauthorizationconfiguration) | No       | None                                          | OPA Authorization Server used for authorizations after authentication (through OIDC, without authentication, no authorization will be done) |
| center                 | [CenterConfiguration](#centerconfiguration)                                 | Yes      | None                                          | OPA Center specific configurations                                                                                                          |

## LogConfiguration

| Key      | Type   | Required | Default | Description                                         |
| -------- | ------ | -------- | ------- | --------------------------------------------------- |
| level    | String | No       | `info`  | Log level                                           |
| format   | String | No       | `json`  | Log format (available values are: `json` or `text`) |
| filePath | String | No       | `""`    | Log file path                                       |

## TracingConfiguration

| Key           | Type              | Required       | Default | Description                           |
| ------------- | ----------------- | -------------- | ------- | ------------------------------------- |
| enabled       | Boolean           | No             | `false` | Enable tracing (Jaeger)               |
| logSpan       | Boolean           | No             | `false` | Should the logger log the span sent ? |
| flushInterval | String            | No             | `""`    | Flush interval                        |
| udpHost       | String            | Yes if enabled | `""`    | UDP Host to send span trace           |
| queueSize     | Integer           | No             | `nil`   | Queue size                            |
| fixedTags     | Map[String]String | No             | `nil`   | Custom tags to be added on spans      |

## ServerConfiguration

| Key        | Type                                  | Required | Default | Description        |
| ---------- | ------------------------------------- | -------- | ------- | ------------------ |
| listenAddr | String                                | No       | `""`    | Listen Address     |
| port       | Integer                               | No       | `8080`  | Listening Port     |
| cors       | [ServerCorsConfig](#servercorsconfig) | No       | `nil`   | CORS configuration |

## ServerCorsConfig

This feature is powered by [gin-contrib/cors](https://github.com/gin-contrib/cors/). You can read more documentation about all field there.

| Key                     | Type     | Required | Default | Description                                                                                                                                                                        |
| ----------------------- | -------- | -------- | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| allowOrigins            | [String] | No       | `nil`   | Allow origins array. Example: https://fake.com. This support stars in origins.                                                                                                     |
| allowMethods            | [String] | No       | `nil`   | Allow HTTP Methods                                                                                                                                                                 |
| allowHeaders            | [String] | No       | `nil`   | Allow headers                                                                                                                                                                      |
| exposeHeaders           | [String] | No       | `nil`   | Expose headers                                                                                                                                                                     |
| maxAgeDuration          | String   | No       | `""`    | Max age.                                                                                                                                                                           |
| allowCredentials        | Boolean  | No       | `nil`   | Allow credentials                                                                                                                                                                  |
| allowWildcard           | Boolean  | No       | `nil`   | Allow wildcard                                                                                                                                                                     |
| allowBrowserExtensions  | Boolean  | No       | `nil`   |  Allow Browser Extensions                                                                                                                                                          |
| allowWebSockets         | Boolean  | No       | `nil`   | Allow websockets                                                                                                                                                                   |
| allowFiles              | Boolean  | No       | `nil`   | Allow files                                                                                                                                                                        |
| allowAllOrigins         | Boolean  | No       | `nil`   | Allow all origins                                                                                                                                                                  |
| useDefaultConfiguration | Boolean  | No       | `nil`   | This flag is used to use gin-contrib/cors default configuration and set all custom configuration on top of it. If only this is set, only the default configuration will be applied |

## DatabaseConfiguration

Database connection and management are based on [Gorm](https://gorm.io/). Only PostgreSQL connection is supported.

| Key                              | Type                                                | Required | Default | Description                                                                                                                                                                                                   |
| -------------------------------- | --------------------------------------------------- | -------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| connectionUrl                    | [CredentialConfiguration](#credentialconfiguration) | Yes      | None    | PostgreSQL URL like this one: `postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable` or this one `host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable` |
| disableForeignKeyWhenMigrating   | Boolean                                             | No       | `false` | Disable foreign keys for auto migration (Gorm option)                                                                                                                                                         |
| allowGlobalUpdate                | Boolean                                             | No       | `false` | Allow global update (Gorm option)                                                                                                                                                                             |
| prepareStatement                 | Boolean                                             | No       | `false` | Prepare statement (Gorm option)                                                                                                                                                                               |
| sqlMaxIdleConnections            | Integer                                             | No       | None    | SQL maximum idle connections for the pool                                                                                                                                                                     |
| sqlMaxOpenConnections            | Integer                                             | No       | None    | SQL maximum open connections for the pool                                                                                                                                                                     |
| sqlConnectionMaxLifetimeDuration | String                                              | No       | `""`    | SQL connection max lifetime duration                                                                                                                                                                          |

## OIDCAuthenticationConfiguration

| Key           | Type                                                | Required | Default                          | Description                                                              |
| ------------- | --------------------------------------------------- | -------- | -------------------------------- | ------------------------------------------------------------------------ |
| clientId      | String                                              | Yes      | None                             | Client ID                                                                |
| clientSecret  | [CredentialConfiguration](#credentialconfiguration) | No       | None                             | Client Secret                                                            |
| issuerUrl     | String                                              | Yes      | None                             | Issuer URL (example: https://fake.com/realm/fake-realm                   |
| redirectUrl   | String                                              | Yes      | None                             | Redirect URL (this is the service url)                                   |
| scopes        | [String]                                            | No       | `["openid", "profile", "email"]` | Scopes                                                                   |
| state         | String                                              | Yes      | None                             | Random string to have a secure connection with oidc provider             |
| emailVerified | Boolean                                             | No       | `false`                          | Check that user email is verified in user token (field `email_verified`) |
| cookieName    | String                                              | No       | `oidc`                           | Cookie generated name                                                    |
| cookieSecure  | Boolean                                             | No       | `false`                          | Is the cookie generated secure ?                                         |

## OPAServerAuthorizationConfiguration

You can see the input and output format of requests made to OPA server [here](opa-formats.md).

| Key  | Type              | Required | Default | Description                                                                                                         |
| ---- | ----------------- | -------- | ------- | ------------------------------------------------------------------------------------------------------------------- |
| url  | String            | Yes      | None    | OPA server url for authorizations checks. This URL mustn't be the url to the default decision but the complete one. |
| tags | Map[String]String | No       | `nil`   | Tags that will be added to each requests done to OPA                                                                |

## CenterConfiguration

| Key                               | Type    | Required | Default                                                                                                                                                                                                                                   | Description                                                                        |
| --------------------------------- | ------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| baseUrl                           | String  | Yes      | None                                                                                                                                                                                                                                      | OPA Center url for generated configuration or others things                        |
| cronRetentionProcess              | String  | Yes      | Cron to start retention process. This will start the retention process to remove data following maximum time declared for status data and decision logs. The cron input must be accepted by [robfig/cron](https://github.com/robfig/cron) |
| skipCronRetentionProcessAtStartup | Boolean | No       | `false`                                                                                                                                                                                                                                   | Retention process will be started at startup without this being filled with `true` |

## Example

This example will show all possible configurations in only 1 file. As said before, you can split it in all needed files.

```yaml
# Log configuration
log:
  # Log level
  level: info
  # Log format
  format: text
  # Log file path
  # filePath:

# Tracing configuration
tracing:
  # Enabled tracing
  enabled: false
  # Log span
  # logSpan: false
  # Flush interval
  # flushInterval:
  # UDP Host
  # udpHost: localhost:6831
  # Queue size
  # queueSize:
  # Fixed tags
  # fixedTags:
  #   tag1: value1

# Public server configuration
# server:
#   # Listen address
#   listenAddr:
#   # Port
#   port: 8080
#   # CORS configurations
#   # cors:
#   #   # Allow origins
#   #   allowOrigins: []
#   #   # Allow methods
#   #   allowMethods: []
#   #   # Allow headers
#   #   allowHeaders: []
#   #   # Expose headers
#   #   exposeHeaders: []
#   #   # Max age duration
#   #   maxAgeDuration: 300
#   #   # Allow credentials
#   #   allowCredentials: false
#   #   # Allow wildcard
#   #   allowWildcard: false
#   #   # Allow browser extensions
#   #   allowBrowserExtensions: false
#   #   # Allow websockets
#   #   allowWebSockets: false
#   #   # Allow origins
#   #   allowAllOrigins: false
#   #   # Use default configuration
#   #   useDefaultConfiguration: true

# Internal server configuration
# internalServer:
#   # Listen address
#   listenAddr:
#   # Port
#   port: 9090

# OPA Publisher server configuration
# server:
#   # Listen address
#   listenAddr:
#   # Port
#   port: 8081
#   # CORS configurations
#   # cors:
#   #   # Allow origins
#   #   allowOrigins: []
#   #   # Allow methods
#   #   allowMethods: []
#   #   # Allow headers
#   #   allowHeaders: []
#   #   # Expose headers
#   #   exposeHeaders: []
#   #   # Max age duration
#   #   maxAgeDuration: 300
#   #   # Allow credentials
#   #   allowCredentials: false
#   #   # Allow wildcard
#   #   allowWildcard: false
#   #   # Allow browser extensions
#   #   allowBrowserExtensions: false
#   #   # Allow websockets
#   #   allowWebSockets: false
#   #   # Allow origins
#   #   allowAllOrigins: false
#   #   # Use default configuration
#   #   useDefaultConfiguration: true

# Database configurations
database:
  # Connection URL (PostgreSQL)
  # Can come from an environment variable, a file or a value directly
  connectionUrl:
    # env:
    # path:
    value: host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable
  # Disable foreign key when migrating
  # disableForeignKeyWhenMigrating: false
  # Allow global update
  # allowGlobalUpdate: false
  # Prepare statement
  # prepareStatement: false
  # SQL Max idle connections
  # sqlMaxIdleConnections: 10
  # SQL Max open connections
  # sqlMaxOpenConnections: 10
  # SQL connection max lifetime duration
  # sqlConnectionMaxLifetimeDuration: 5m

# OIDC Authentication
# oidcAuthentication:
#   # Client id
#   clientId: client-id
#   # Client secret
#   clientSecret:
#     env:
#     path:
#     value:
#   # Issuer URL
#   issuerUrl: http://localhost:8088/auth/realms/integration
#   # Redirect URL
#   redirectUrl: http://localhost:3000/ # /auth/oidc/callback will be added by the application
#   # Scopes
#   scopes: [openid, email, profile]
#   # State
#   state: random-string
#   # Cookie name
#   cookieName: oidc
#   # Cookie secure
#   cookieSecure: true
#   # Email must be marked as verified in the token ?
#   emailVerified: true

# OPA Server authorization
# opaServerAuthorization:
#   # OPA url
#   url: http://localhost:8181/v1/data/example/authz/allowed
#   # Tags
#   tags:
#     tag1: value1

# OPA Center configurations
center:
  # OPA Center url
  baseUrl: http://localhost:8080
  # Cron for starting retention process
  cronRetentionProcess: "@every 30s"
  # Skip retention process at startup
  skipRetentionProcessAtStartup: false
```
