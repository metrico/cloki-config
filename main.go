package main

import (
	"crypto/tls"

	//"database/sql"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/metrico/cloki-config/config"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Note: Vipers are not safe for concurrent Get() and Set() operations.
type ClokiConfig struct {
	// Delimiter that separates a list of keys
	// used to access a nested value in one go
	keyDelim string

	// A set of paths to look for the config file in
	configPaths []string
}

func New() *ClokiConfig {

	c := new(ClokiConfig)

	return c
}

func (c *ClokiConfig) checkHelpVersionFlags() {
	if *appFlags.ShowHelpMessage {
		flag.Usage()
		os.Exit(0)
	}

	if *appFlags.ShowVersion {
		fmt.Printf("VERSION: %s\r\n", VERSION_APPLICATION)
		os.Exit(0)
	}

	if *appFlags.GenerateKey {
		strLic, err := license.GenerateKey()
		if err == nil {
			if _, err := os.Stdout.WriteString(strLic + "\r\n"); err != nil {
				fmt.Print("couldn't generate a key: ", err.Error())
			}
		} else {
			fmt.Print("couldn't generate a key: ", err.Error())
		}

		os.Exit(0)
	}

}

//https://github.com/atreugo/examples/blob/master/basic/main.go
//https://github.com/jackwhelpton/fasthttp-routing
func (c *ClokiConfig) readConfig() {
	// Getting constant values
	if configEnv := os.Getenv("ClOKIWRITERAPPENV"); configEnv != "" {
		viper.SetConfigName("cloki-writer_" + configEnv)
	} else {
		viper.SetConfigName("cloki-writer")
	}
	viper.SetConfigType("json")

	if configPath := os.Getenv("ClOKIWRITERAPPPATH"); configPath != "" {
		viper.AddConfigPath(configPath)
	} else {
		viper.AddConfigPath(*appFlags.ConfigPath)
	}

	viper.AddConfigPath(".")

	//Default value
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - checking env: ", err)
		logger.Error("No configuration file loaded - using defaults - checking env")
	}

	viper.SetConfigName("cloki-writer_custom")
	err = viper.MergeInConfig()
	if err != nil {
		logger.Debug("No custom configuration file loaded.")
	}

	//Env variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "[", "_", "]", ""))
	viper.SetEnvPrefix(config.Setting.EnvPrefix)
	SetEnvironDataBase()

	//Bind Env from Config
	BindEnvs(config.Setting)

	err = viper.Unmarshal(&config.Setting, func(config *mapstructure.DecoderConfig) {
		config.TagName = "json"
	})
	if err != nil {
		logger.Debug("couldn't unmarshal viper.")
	}

	var re = regexp.MustCompile(`database_data\[(\d)\]`)
	envParams := []int{}
	allSettings := viper.AllSettings()
	for key, _ := range allSettings {
		if strings.HasPrefix(key, "database_data[") {
			key = re.ReplaceAllString(key, "$1")
			i, err := strconv.Atoi(key)
			if err == nil {
				envParams = append(envParams, i)
			}
		}
	}

	//Read the data db connection from config. This is fix because defaults doesn't know the size of array
	if viper.IsSet("database_data") {
		config.Setting.DATABASE_DATA = nil
		dataConfig := viper.Get("database_data")
		dataVal := dataConfig.([]interface{})
		for idx := range dataVal {
			val := dataVal[idx].(map[string]interface{})
			data := config.ClokiWriterDataBase{}
			defaults.SetDefaults(&data) //<-- This set the defaults values
			err := mapstructure.Decode(val, &data)
			if err != nil {
				logger.Error("ERROR during mapstructure decode[1]:", err)
			}
			config.Setting.DATABASE_DATA = append(config.Setting.DATABASE_DATA, data)
		}
	}

	//We should do extraction and after sorting 0,1,2,3
	sort.Ints(envParams[:])
	//Here we do ENV check
	for _, idx := range envParams {
		value := allSettings[fmt.Sprintf("database_data[%d]", idx)]
		val := value.(map[string]interface{})
		//If the configuration already exists - we replace only existing params
		if len(config.Setting.DATABASE_DATA) > idx {
			err := mapstructure.Decode(val, &config.Setting.DATABASE_DATA[idx])
			if err != nil {
				logger.Error("ERROR during mapstructure decode[0]:", err)
			}
		} else {
			data := config.ClokiWriterDataBase{}
			defaults.SetDefaults(&data) //<-- This set the defaults values
			err := mapstructure.Decode(val, &data)
			if err != nil {
				logger.Error("ERROR during mapstructure decode[1]:", err)
			}
			config.Setting.DATABASE_DATA = append(config.Setting.DATABASE_DATA, data)
		}
	}

	//Check the command line
	if *appFlags.LogName != "" {
		config.Setting.LOG_SETTINGS.Name = *appFlags.LogName
	}

	if *appFlags.LogPath != "" {
		config.Setting.LOG_SETTINGS.Path = *appFlags.LogPath
	}

	//viper.Debug()
}

//system params for replications, groups
func (c *ClokiConfig) setFastConfigSettings() {

	/***********************************/
	switch config.Setting.SYSTEM_SETTINGS.HashType {
	case "cityhash":
		config.Setting.FingerPrintType = config.FINGERPRINT_CityHash
	case "bernstein":
	case "default":
		config.Setting.FingerPrintType = config.FINGERPRINT_Bernstein
	}

	minVersion := config.Setting.HTTPS_SETTINGS.MinTLSVersionString

	if minVersion == "TLS1.0" {
		config.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS10
	} else if minVersion == "TLS1.1" {
		config.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS11
	} else if minVersion == "TLS1.2" {
		config.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS12
	} else if minVersion == "TLS1.3" {
		config.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS13
	}

	maxVersion := config.Setting.HTTPS_SETTINGS.MaxTLSVersionString

	if maxVersion == "TLS1.0" {
		config.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS10
	} else if maxVersion == "TLS1.1" {
		config.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS11
	} else if maxVersion == "TLS1.2" {
		config.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS12
	} else if maxVersion == "TLS1.3" {
		config.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS13
	}
}

//this function will check CLOKI_DATABASE_DATA and set internal bind for viper
//i.e. CLOKI_DATABASE_DATA_0_HOSTNAME -> database_data[0].hostname
func (c *ClokiConfig) SetEnvironDataBase() bool {

	var re = regexp.MustCompile(`_(\d)_`)
	for _, s := range os.Environ() {
		if strings.HasPrefix(s, config.Setting.EnvPrefix+"_DATABASE_DATA") {
			a := strings.Split(s, "=")
			key := strings.TrimPrefix(a[0], config.Setting.EnvPrefix+"_")
			key = re.ReplaceAllString(key, "[$1].")
			viper.BindEnv(key)
		}
	}
	return true
}

//Now we should bind the ENV params
func (c *ClokiConfig) BindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			c.BindEnvs(v.Interface(), append(parts, tv)...)
		case reflect.Slice:
			continue
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}
