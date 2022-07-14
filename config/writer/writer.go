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

type KafkaTopicConf struct {
	Name           string            `json:"name" mapstructure:"name" default:""`
	Tags           map[string]string `json:"tags" mapstructure:"tags" default:""`
	IncludeTopic   bool              `json:"include_topic" mapstructure:"include_topic" default:"true"`
	IncludeHeaders bool              `json:"include_headers" mapstructure:"include_headers" default:"true"`
	IncludeKey     bool              `json:"include_key" mapstructure:"include_key" default:"true"`
	// json or empty
	Format       string     `json:"format" mapstructure:"format" default:""`
	Extract      [][]string `json:"extract" mapstructure:"extract" default:"[]"`
	ExtractValue []string   `json:"extract_value" mapstructure:"extract_value" default:"[]"`
	Partition    int        `json:"partition" mapstructure:"partition" default:"0"`
	Type         string     `json:"type" mapstructure:"type" default:""`
}

type KafkaConf struct {
	Host string `json:"host" mapstructure:"host" default:""`
	Port int    `json:"port" mapstructure:"port" default:"0"`
	// SASL auth type:  Plain/Scram
	AuthType string           `json:"auth_type" mapstructure:"auth_type" default:""`
	Username string           `json:"username" mapstructure:"username" default:""`
	Password string           `json:"password" mapstructure:"password" default:""`
	Topic    []KafkaTopicConf `json:"topic" mapstructure:"topic" default:""`
	GroupID  string           `json:"consumer_group" mapstructure:"consumer_group" default:""`
	// org id to write into DB
	OrgID string `json:"org_id" mapstructure:"org_id" default:"0"`
}

type PrometheusScrape struct {
	Endpoint    string `json:"endpoint" mapstructure:"endpoint" default:"https://user:password@127.0.0.1:9099/metrics"`
	InstanceTag string `json:"instance_tag" mapstructure:"instance_tag" default:"instance"`
	EndpointTag string `json:"endpoint_tag" mapstructure:"endpoint_tag" default:"endpoint"`
	Enable      bool   `json:"enable" mapstructure:"enable" default:"false"`
	// org id to write into DB
	OrgID string `json:"org_id" mapstructure:"org_id" default:"0"`
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
		// org id to write into DB
		OrgID string `json:"org_id" mapstructure:"org_id" default:"0"`
	} `json:"mqtt" mapstructure:"mqtt"`
	KAFKA_CLIENT []KafkaConf `json:"kafka" mapstructure:"kafka" default:""`
	NATS         struct {
		External string `json:"external" mapstructure:"external" default:""`
		Bind     string `json:"bind" mapstructure:"bind" default:"0.0.0.0"`
		Port     int    `json:"port" mapstructure:"port" default:"4444"`
		User     string `json:"user" mapstructure:"user" default:""`
		Password string `json:"password" mapstructure:"password" default:""`
	} `json:"nats" mapstructure:"nats"`
	PROMETHEUS_SCRAPE        []PrometheusScrape `json:"prometheus_scrape" mapstructure:"prometheus_scrape"`
	PrometheusScrapeInterval string             `json:"prometheus_scrape_interval" mapstructure:"prometheus_scrape_interval" default:"15s"`
}
