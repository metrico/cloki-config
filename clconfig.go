package clconfig

import (
	"crypto/tls"

	//"database/sql"

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

// Note: ClokiConfig
type ClokiConfig struct {
	configPath   string
	logName      string
	logPath      string
	typeCloki    CLOKI_TYPE
	ConfigWriter *config.ClokiWriterSettingServer
}

type CLOKI_TYPE int

const (
	CLOKI_WRITER CLOKI_TYPE = iota
	CLOKI_READER
)

func New(typeCloki CLOKI_TYPE, configPath, logName, logPath string) *ClokiConfig {

	c := new(ClokiConfig)

	c.ConfigWriter = new(config.ClokiWriterSettingServer)
	defaults.SetDefaults(c.ConfigWriter) //<-- This set the defaults values

	c.configPath = configPath
	c.logName = logName
	c.logPath = logPath

	//Type
	c.typeCloki = typeCloki

	return c
}

func (c *ClokiConfig) ReadConfig() {
	// Getting constant values

	envCloki := "ClOKIWRITERAPPENV"
	envPathCloki := "ClOKIWRITERAPPPATH"
	envConfig := "cloki-writer"

	if c.typeCloki == CLOKI_READER {
		envCloki = "CLOKIGOENV"
		envPathCloki = "CLOKIGOPATH"
		envConfig = "cloki_go"
	}

	if configEnv := os.Getenv(envCloki); configEnv != "" {
		viper.SetConfigName(envConfig + "_" + configEnv)
	} else {
		viper.SetConfigName(envConfig)
	}
	viper.SetConfigType("json")

	if configPath := os.Getenv(envPathCloki); configPath != "" {
		viper.AddConfigPath(configPath)
	} else {
		viper.AddConfigPath(c.configPath)
	}

	viper.AddConfigPath(".")

	//Default value
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - checking env: ", err)
	}

	viper.SetConfigName(envConfig + "_custom")
	err = viper.MergeInConfig()
	if err != nil {
		fmt.Println("No custom configuration file loaded.")
	}

	//Env variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "[", "_", "]", ""))
	viper.SetEnvPrefix(config.Setting.EnvPrefix)
	c.setEnvironDataBase()

	//Bind Env from Config
	c.bindEnvs(config.Setting)

	err = viper.Unmarshal(&config.Setting, func(config *mapstructure.DecoderConfig) {
		config.TagName = "json"
	})
	if err != nil {
		fmt.Println("couldn't unmarshal viper.")
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
				fmt.Println("ERROR during mapstructure decode[1]:", err)
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
				fmt.Println("ERROR during mapstructure decode[0]:", err)
			}
		} else {
			data := config.ClokiWriterDataBase{}
			defaults.SetDefaults(&data) //<-- This set the defaults values
			err := mapstructure.Decode(val, &data)
			if err != nil {
				fmt.Println("ERROR during mapstructure decode[1]:", err)
			}
			config.Setting.DATABASE_DATA = append(config.Setting.DATABASE_DATA, data)
		}
	}

	//Check the command line
	if c.logName != "" {
		config.Setting.LOG_SETTINGS.Name = c.logName
	}

	if c.logPath != "" {
		config.Setting.LOG_SETTINGS.Path = c.logPath
	}

	//viper.Debug()
	c.SetFastConfigSettings()
}

//system params for replications, groups
func (c *ClokiConfig) SetFastConfigSettings() {

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
func (c *ClokiConfig) setEnvironDataBase() bool {

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
func (c *ClokiConfig) bindEnvs(iface interface{}, parts ...string) {
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
			c.bindEnvs(v.Interface(), append(parts, tv)...)
		case reflect.Slice:
			continue
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}
