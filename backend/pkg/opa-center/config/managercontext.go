package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

// Main configuration folder path.
var mainConfigFolderPath = "conf/"

// TemplateErrLoadingEnvCredentialEmpty Template Error when Loading Environment variable Credentials.
var TemplateErrLoadingEnvCredentialEmpty = "error loading credentials, environment variable %s is empty"

var validate = validator.New()

type managercontext struct {
	cfg                       *Config
	configs                   []*viper.Viper
	onChangeHooks             []func()
	logger                    log.Logger
	internalFileWatchChannels []chan bool
}

func (ctx *managercontext) AddOnChangeHook(hook func()) {
	ctx.onChangeHooks = append(ctx.onChangeHooks, hook)
}

func (ctx *managercontext) Load() error {
	// List files
	files, err := ioutil.ReadDir(mainConfigFolderPath)
	if err != nil {
		return err
	}

	// Generate viper instances for static configs
	ctx.configs = generateViperInstances(files)

	// Load configuration
	err = ctx.loadConfiguration()
	if err != nil {
		return err
	}

	// Loop over config files
	funk.ForEach(ctx.configs, func(vv interface{}) {
		// Cast viper object
		vip := vv.(*viper.Viper)

		// Add hooks for on change events
		vip.OnConfigChange(func(in fsnotify.Event) {
			ctx.logger.Infof("Reload configuration detected for file %s", vip.ConfigFileUsed())

			// Reload config
			err2 := ctx.loadConfiguration()
			if err2 != nil {
				ctx.logger.Error(err2)
				// Stop here and do not call hooks => configuration is unstable
				return
			}
			// Call all hooks in sequence in order to manage correctly reload database and after
			// services that depends on it
			funk.ForEach(ctx.onChangeHooks, func(hook func()) { hook() })
		})
		// Watch for configuration changes
		vip.WatchConfig()
	})

	return nil
}

// Imported and modified from viper v1.7.0.
func (ctx *managercontext) watchInternalFile(filePath string, forceStop chan bool, onChange func()) {
	initWG := sync.WaitGroup{}
	initWG.Add(1)

	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			ctx.logger.Fatal(err)
		}
		defer watcher.Close()

		configFile := filepath.Clean(filePath)
		configDir, _ := filepath.Split(configFile)
		realConfigFile, _ := filepath.EvalSymlinks(filePath)

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)

		go func() {
			for {
				select {
				case <-forceStop:
					eventsWG.Done()

					return
				case event, ok := <-watcher.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()

						return
					}

					currentConfigFile, _ := filepath.EvalSymlinks(filePath)
					// we only care about the config file with the following cases:
					// 1 - if the config file was modified or created
					// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if (filepath.Clean(event.Name) == configFile &&
						event.Op&writeOrCreateMask != 0) ||
						(currentConfigFile != "" && currentConfigFile != realConfigFile) {
						realConfigFile = currentConfigFile

						// Call on change
						onChange()
					} else if filepath.Clean(event.Name) == configFile && event.Op&fsnotify.Remove&fsnotify.Remove != 0 {
						eventsWG.Done()

						return
					}

				case err, ok := <-watcher.Errors:
					if ok { // 'Errors' channel is not closed
						ctx.logger.Errorf("watcher error: %v\n", err)
					}

					eventsWG.Done()

					return
				}
			}
		}()

		_ = watcher.Add(configDir)

		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
}

func (ctx *managercontext) loadDefaultConfigurationValues(vip *viper.Viper) {
	// Load default configuration
	vip.SetDefault("log.level", DefaultLogLevel)
	vip.SetDefault("log.format", DefaultLogFormat)
	vip.SetDefault("server.port", DefaultPort)
	vip.SetDefault("opaPublisherServer.port", DefaultOPAPublisherPort)
	vip.SetDefault("internalServer.port", DefaultInternalPort)
	vip.SetDefault("lockDistributor.tableName", DefaultLockDistributorTableName)
	vip.SetDefault("lockDistributor.leaseDuration", DefaultLockDistributorLeaseDuration)
	vip.SetDefault("lockDistributor.heartbeatFrequency", DefaultLockDistributionHeartbeatFrequency)
}

func generateViperInstances(files []os.FileInfo) []*viper.Viper {
	list := make([]*viper.Viper, 0)
	// Loop over static files to create viper instance for them
	funk.ForEach(files, func(file os.FileInfo) {
		filename := file.Name()
		// Create config file name
		cfgFileName := strings.TrimSuffix(filename, path.Ext(filename))
		// Test if config file name is compliant (ignore hidden files like .keep or directory)
		if !strings.HasPrefix(filename, ".") && cfgFileName != "" && !file.IsDir() {
			// Create new viper instance
			vip := viper.New()
			// Set config name
			vip.SetConfigName(cfgFileName)
			// Add configuration path
			vip.AddConfigPath(mainConfigFolderPath)
			// Append it
			list = append(list, vip)
		}
	})

	return list
}

func (ctx *managercontext) loadConfiguration() error {
	// Load must start by flushing all existing watcher on internal files
	for i := 0; i < len(ctx.internalFileWatchChannels); i++ {
		ch := ctx.internalFileWatchChannels[i]
		// Send the force stop
		ch <- true
	}

	// Create a viper instance for default value and merging
	globalViper := viper.New()

	// Put default values
	ctx.loadDefaultConfigurationValues(globalViper)

	// Loop over configs
	for _, vip := range ctx.configs {
		err := vip.ReadInConfig()
		if err != nil {
			return err
		}

		err = globalViper.MergeConfigMap(vip.AllSettings())
		if err != nil {
			return err
		}
	}

	// Prepare configuration object
	var out Config
	// Quick unmarshal.
	err := globalViper.Unmarshal(&out)
	if err != nil {
		return err
	}

	// Load default values
	err = loadBusinessDefaultValues(&out)
	if err != nil {
		return err
	}

	// Configuration validation
	err = validate.Struct(out)
	if err != nil {
		return err
	}

	// Load all credentials
	credentials, err := loadAllCredentials(&out)
	if err != nil {
		return err
	}
	// Initialize or flush watch internal file channels
	internalFileWatchChannels := make([]chan bool, 0)
	ctx.internalFileWatchChannels = internalFileWatchChannels
	// Loop over all credentials in order to watch file change
	funk.ForEach(credentials, func(item interface{}) {
		cred := item.(*CredentialConfig)
		// Check if credential is about a path
		if cred.Path != "" {
			// Create channel
			ch := make(chan bool)
			// Run the watch file
			ctx.watchInternalFile(cred.Path, ch, func() {
				// File change detected
				ctx.logger.Infof("Reload credential file detected for path %s", cred.Path)

				// Reload config
				err2 := loadCredential(cred)
				if err2 != nil {
					ctx.logger.Error(err2)
					// Stop here and do not call hooks => configuration is unstable
					return
				}
				// Call all hooks in sequence in order to manage correctly reload database and after
				// services that depends on it
				funk.ForEach(ctx.onChangeHooks, func(hook func()) { hook() })
			})
			// Add channel to list of channels
			ctx.internalFileWatchChannels = append(ctx.internalFileWatchChannels, ch)
		}
	})

	err = validateBusinessConfig(&out)
	if err != nil {
		return err
	}

	ctx.cfg = &out

	return nil
}

// Load default values based on business rules.
func loadBusinessDefaultValues(out *Config) error {
	// Load default oidc configurations
	if out.OIDCAuthentication != nil {
		// Add default scopes
		if out.OIDCAuthentication.Scopes == nil {
			out.OIDCAuthentication.Scopes = DefaultOIDCScopes
		}
		// Add default cookie name
		if out.OIDCAuthentication.CookieName == "" {
			out.OIDCAuthentication.CookieName = DefaultCookieName
		}
	}

	// Load default tags for opa authorization
	if out.OPAServerAuthorization != nil && out.OPAServerAuthorization.Tags == nil {
		out.OPAServerAuthorization.Tags = map[string]string{}
	}

	// Load default tracing configuration
	if out.Tracing == nil {
		out.Tracing = &TracingConfig{Enabled: false}
	}

	// TODO Load default values based on business rules
	return nil
}

// Load credential configs here.
func loadAllCredentials(out *Config) ([]*CredentialConfig, error) {
	// Initialize answer
	result := make([]*CredentialConfig, 0)

	// Load database credential
	err := loadCredential(out.Database.ConnectionURL)
	if err != nil {
		return nil, err
	}
	// Append result
	result = append(result, out.Database.ConnectionURL)

	// Load credential for OIDC configuration
	if out.OIDCAuthentication != nil && out.OIDCAuthentication.ClientSecret != nil {
		err := loadCredential(out.OIDCAuthentication.ClientSecret)
		if err != nil {
			return nil, err
		}
		// Append result
		result = append(result, out.OIDCAuthentication.ClientSecret)
	}

	// TODO Load credential configs here

	return result, nil
}

func loadCredential(credCfg *CredentialConfig) error {
	if credCfg.Path != "" {
		// Secret file
		databytes, err := ioutil.ReadFile(credCfg.Path)
		if err != nil {
			return err
		}
		// Store value
		credCfg.Value = string(databytes)
	} else if credCfg.Env != "" {
		// Environment variable
		envValue := os.Getenv(credCfg.Env)
		if envValue == "" {
			return fmt.Errorf(TemplateErrLoadingEnvCredentialEmpty, credCfg.Env)
		}
		// Store value
		credCfg.Value = envValue
	}
	// Default value
	return nil
}

// GetConfig allow to get configuration object.
func (ctx *managercontext) GetConfig() *Config {
	return ctx.cfg
}
