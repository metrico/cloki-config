package writer

const (
	FINGERPRINT_Bernstein = iota
	FINGERPRINT_CityHash
)

// ============================ WRITER ONLY ================================= //
type MQTTTopicConf struct {
	Name           string            `json:"name" mapstructure:"name" default:""`
	Tags           map[string]string `json:"tags" mapstructure:"tags" default:""`
	IncludeTopic   bool              `json:"include_topic" mapstructure:"include_topic" default:"true"`
	Format         string            `json:"format" mapstructure:"format" default:""`
	Extract        [][]string        `json:"extract" mapstructure:"extract" default:"[]"`
	TimestampField []string          `json:"timestamp" mapstructure:"timestamp" default:"[]"`
	// ns, us, ms, s, m, h
	TimeUnit string `json:"time_unit" mapstructure:"time_unit" default:""`
}

//
type ClokiWriterSettingServer struct {
	MQTT_CLIENT struct {
		SessID     string        `json:"session" mapstructure:"session" default:"cloki-client"`
		ServerHost string        `json:"host" mapstructure:"host" default:""`
		ServerPort uint32        `json:"port" mapstructure:"port" default:""`
		User       string        `json:"user" mapstructure:"user" default:""`
		Password   string        `json:"password" mapstructure:"password" default:""`
		Topic      MQTTTopicConf `json:"topic" mapstructure:"topic" default:""`
	} `json:"mqtt" mapstructure:"mqtt"`
}
