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
	"github.com/metrico/cloki-config/config/writer"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Note: ClokiConfig
type ClokiConfig struct {
	configPaths []string
	logName     string
	logPath     string
	typeCloki   CLOKI_TYPE
	Setting     *config.ClokiBaseSettingServer
}

type CLOKI_TYPE int

const (
	CLOKI_WRITER CLOKI_TYPE = iota
	CLOKI_READER
)

func New(typeCloki CLOKI_TYPE, configPath string, mergePaths []string, logName, logPath string) *ClokiConfig {

	c := new(ClokiConfig)

	c.Setting = new(config.ClokiBaseSettingServer)
	defaults.SetDefaults(c.Setting) //<-- This set the defaults values

	c.configPaths = []string{configPath}
	for _, p := range mergePaths {
		c.configPaths = append(c.configPaths, p)
	}
	c.logPath = logPath

	//Type
	c.typeCloki = typeCloki

	return c
}

func (c *ClokiConfig) readConfigFromFS() {
	cnt := 0
	for _, p := range c.configPaths {
		if len(p) < 5 || p[len(p)-5:] != ".json" {
			fmt.Printf("only json extension allowed: %s\n", p)
			continue
		}
		viper.SetConfigFile(p)
		if cnt == 0 {
			err := viper.ReadInConfig()
			if err != nil {
				fmt.Println(err)
				continue
			}
			cnt++
			continue
		}
		err := viper.MergeInConfig()
		if err != nil {
			fmt.Println(err)
		}
		cnt++
	}
}

func (c *ClokiConfig) ReadConfig() {
	// Getting constant values

	c.readConfigFromFS()

	//Env variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "[", "_", "]", ""))
	viper.SetEnvPrefix(c.Setting.EnvPrefix)
	c.setEnvironDataBase()

	//Bind Env from Config
	c.bindEnvs(*c.Setting)

	err := viper.Unmarshal(c.Setting, func(config *mapstructure.DecoderConfig) {
		config.TagName = "json"
	})
	if err != nil {
		fmt.Println("couldn't unmarshal viper.")
	}

	var re = regexp.MustCompile(`database_data\[(\d)\]`)
	envParams := []int{}
	allSettings := viper.AllSettings()
	for key := range allSettings {
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
		c.Setting.DATABASE_DATA = nil
		dataConfig := viper.Get("database_data")
		dataVal := dataConfig.([]interface{})
		for idx := range dataVal {
			val := dataVal[idx].(map[string]interface{})
			data := config.ClokiBaseDataBase{}
			defaults.SetDefaults(&data) //<-- This set the defaults values
			err := mapstructure.Decode(val, &data)
			if err != nil {
				fmt.Println("ERROR during mapstructure decode[1]:", err)
			}
			c.Setting.DATABASE_DATA = append(c.Setting.DATABASE_DATA, data)
		}
	}

	//We should do extraction and after sorting 0,1,2,3
	sort.Ints(envParams[:])
	//Here we do ENV check
	for _, idx := range envParams {
		value := allSettings[fmt.Sprintf("database_data[%d]", idx)]
		val := value.(map[string]interface{})
		//If the configuration already exists - we replace only existing params
		if len(c.Setting.DATABASE_DATA) > idx {
			err := mapstructure.WeakDecode(val, &c.Setting.DATABASE_DATA[idx])
			if err != nil {
				fmt.Println("ERROR during mapstructure decode[0] sort:", err)
			}
		} else {
			data := config.ClokiBaseDataBase{}
			defaults.SetDefaults(&data) //<-- This set the defaults values
			err := mapstructure.WeakDecode(val, &data)
			if err != nil {
				fmt.Println("ERROR during mapstructure decode[1] sort:", err)
			}
			c.Setting.DATABASE_DATA = append(c.Setting.DATABASE_DATA, data)
		}
	}

	//Prometheus scrapper
	re = regexp.MustCompile(`prometheus_scrape\[(\d)\]`)
	envParams = []int{}
	for key := range allSettings {
		if strings.HasPrefix(key, "prometheus_scrape[") {
			key = re.ReplaceAllString(key, "$1")
			i, err := strconv.Atoi(key)
			if err == nil {
				envParams = append(envParams, i)
			}
		}
	}

	//Read the data prometheus_scrape from config. This is fix because defaults doesn't know the size of array
	if viper.IsSet("prometheus_scrape") {
		c.Setting.ClokiWriter.PROMETHEUS_SCRAPE = nil
		scraperConfig := viper.Get("prometheus_scrape")
		dataVal := scraperConfig.([]interface{})
		for idx := range dataVal {
			val := dataVal[idx].(map[string]interface{})
			data := writer.PrometheusScrape{}
			defaults.SetDefaults(&data) //<-- This set the defaults values
			err := mapstructure.Decode(val, &data)
			if err != nil {
				fmt.Println("ERROR during mapstructure decode[1]:", err)
			}
			c.Setting.ClokiWriter.PROMETHEUS_SCRAPE = append(c.Setting.ClokiWriter.PROMETHEUS_SCRAPE, data)
		}
	}

	//We should do extraction and after sorting 0,1,2,3
	sort.Ints(envParams[:])
	//Here we do ENV check
	for _, idx := range envParams {
		value := allSettings[fmt.Sprintf("prometheus_scrape[%d]", idx)]
		val := value.(map[string]interface{})
		//If the configuration already exists - we replace only existing params
		if len(c.Setting.ClokiWriter.PROMETHEUS_SCRAPE) > idx {
			err := mapstructure.WeakDecode(val, &c.Setting.ClokiWriter.PROMETHEUS_SCRAPE[idx])
			if err != nil {
				fmt.Println("ERROR during mapstructure scraper decode[0]:", err)
			}
		} else {
			data := writer.PrometheusScrape{}
			defaults.SetDefaults(&data) //<-- This set the defaults values
			err := mapstructure.WeakDecode(val, &data)
			if err != nil {
				fmt.Println("ERROR during mapstructure scraper decode[1]:", err)
			}
			c.Setting.ClokiWriter.PROMETHEUS_SCRAPE = append(c.Setting.ClokiWriter.PROMETHEUS_SCRAPE, data)
		}
	}

	//Check the command line
	if c.logName != "" {
		c.Setting.LOG_SETTINGS.Name = c.logName
	}

	if c.logPath != "" {
		c.Setting.LOG_SETTINGS.Path = c.logPath
	}

	//viper.Debug()
	c.setFastConfigSettings()

	// default table names setting
	for i, node := range c.Setting.DATABASE_DATA {
		if node.TableSeries == "" {
			if node.ClusterName == "" {
				c.Setting.DATABASE_DATA[i].TableSeries = "time_series_v2"
			} else {
				c.Setting.DATABASE_DATA[i].TableSeries = "time_series_v2_dist"
			}
		}
		if node.TableSamples == "" {
			if node.ClusterName == "" {
				c.Setting.DATABASE_DATA[i].TableSamples = "samples_v4"
			} else {
				c.Setting.DATABASE_DATA[i].TableSamples = "samples_v4_dist"
			}
		}
		if node.TableMetrics == "" {
			if node.ClusterName == "" {
				c.Setting.DATABASE_DATA[i].TableMetrics = "samples_v4"
			} else {
				c.Setting.DATABASE_DATA[i].TableMetrics = "samples_v4_dist"
			}
		}
	}
	//c.Setting = &c.Setting
}

//system params for replications, groups
func (c *ClokiConfig) setFastConfigSettings() {

	/***********************************/
	switch c.Setting.SYSTEM_SETTINGS.HashType {
	case "cityhash":
		c.Setting.FingerPrintType = writer.FINGERPRINT_CityHash
	case "bernstein":
	case "default":
		c.Setting.FingerPrintType = writer.FINGERPRINT_Bernstein
	}

	minVersion := c.Setting.HTTPS_SETTINGS.MinTLSVersionString

	if minVersion == "TLS1.0" {
		c.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS10
	} else if minVersion == "TLS1.1" {
		c.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS11
	} else if minVersion == "TLS1.2" {
		c.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS12
	} else if minVersion == "TLS1.3" {
		c.Setting.HTTPS_SETTINGS.MinTLSVersion = tls.VersionTLS13
	}

	maxVersion := c.Setting.HTTPS_SETTINGS.MaxTLSVersionString

	if maxVersion == "TLS1.0" {
		c.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS10
	} else if maxVersion == "TLS1.1" {
		c.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS11
	} else if maxVersion == "TLS1.2" {
		c.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS12
	} else if maxVersion == "TLS1.3" {
		c.Setting.HTTPS_SETTINGS.MaxTLSVersion = tls.VersionTLS13
	}
}

//this function will check CLOKI_DATABASE_DATA and set internal bind for viper
//i.e. CLOKI_DATABASE_DATA_0_HOSTNAME -> database_data[0].hostname
func (c *ClokiConfig) setEnvironDataBase() bool {

	var re = regexp.MustCompile(`_(\d)_`)
	for _, s := range os.Environ() {
		if strings.HasPrefix(s, c.Setting.EnvPrefix+"_DATABASE_DATA") {
			a := strings.Split(s, "=")
			key := strings.TrimPrefix(a[0], c.Setting.EnvPrefix+"_")
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
