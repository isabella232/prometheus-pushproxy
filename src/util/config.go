package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"unicode"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
)

// DefaultConfigFile - default config file
// it can be overwritten by env variable CONFIG_FILE
const DefaultConfigFile = "../config/default_config.yml"

// Configuration has a set of parameters to configure this server.
// The same name can be used in environment variable to override yml or json values.
type Configuration struct {
	// PORT is the http port
	Port string `json:"port"`

	Name string `json:"name"`

	// LogLevel is used to set the application log level
	LogLevel string `json:"logLevel"`

	// DbPassword is database password
	DbPassword string `json:"dbPassword"`

	// DbConnectionStr is the database connection string
	DbConnectionStr string `json:"dbConnectionStr"`

	// PbDbType is the database type
	PbDbType string `json:"dbDbType"`

	// DefaultAPIKey is the API key, an empty default key will disable validation
	DefaultAPIKey string `json:"defaultAPIKey"`

	// HTTPs certificate set up
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

// Config is global configuration store
var Config Configuration

// Init initializes configuration
func Init() {
	configFile := AssignString(os.Getenv("CONFIG_FILE"), DefaultConfigFile)
	ReadConfigFile(configFile)

	log.SetLevel(logLevel(Config.LogLevel))

	log.Warnf("Configuration built from file - %s", configFile)
}

// ReadConfigFile reads configuration file.
func ReadConfigFile(configFile string) {
	fileBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("failed to load configuration file %s", configFile)
		panic(err)
	}

	if hasJSONPrefix(fileBytes) {
		err = json.Unmarshal(fileBytes, &Config)
		if err != nil {
			panic(err)
		}
	} else {
		err = yaml.Unmarshal(fileBytes, &Config)
		if err != nil {
			panic(err)
		}
	}

	// Next section allows env variable overwrites config file value
	fields := reflect.TypeOf(Config)
	// pointer to struct
	values := reflect.ValueOf(&Config)
	// struct
	st := values.Elem()
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i).Name
		f := st.FieldByName(field)

		if f.Kind() == reflect.String {
			envV := os.Getenv(field)
			if len(envV) > 0 && f.IsValid() && f.CanSet() {
				f.SetString(strings.TrimSuffix(envV, "\n")) // ensure no \n at the end of line that was introduced by loading k8s secrete file
			}
			os.Setenv(field, f.String())
		}
	}

	fmt.Printf("Configuration %v\n", Config)
}

//GetConfig returns a reference to the Configuration
func GetConfig() *Configuration {
	return &Config
}

func logLevel(level string) log.Level {
	switch strings.TrimSpace(strings.ToLower(level)) {
	case "debug":
		return log.DebugLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	default:
		return log.InfoLevel
	}
}

var jsonPrefix = []byte("{")

func hasJSONPrefix(buf []byte) bool {
	return hasPrefix(buf, jsonPrefix)
}

// Return true if the first non-whitespace bytes in buf is prefix.
func hasPrefix(buf []byte, prefix []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, prefix)
}
