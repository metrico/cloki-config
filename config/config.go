package config

import (
	"github.com/metrico/cloki-config/config/reader"
	"github.com/metrico/cloki-config/config/writer"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

// ============================ BASE ONLY ================================= //
type ClokiBaseDataBase struct {
	User                   string `json:"user" mapstructure:"user" default:"cloki_user"`
	Node                   string `json:"node" mapstructure:"node" default:"clokinode"`
	Password               string `json:"pass" mapstructure:"pass" default:"cloki_pass"`
	Name                   string `json:"name" mapstructure:"name" default:"cloki_data"`
	Host                   string `json:"host" mapstructure:"host" default:"127.0.0.1"`
	TableSamples           string `json:"table_samples" mapstructure:"table_samples" default:""`
	TableSeries            string `json:"table_series" mapstructure:"table_series" default:""`
	TableMetrics           string `json:"table_metrics" mapstructure:"table_metrics" default:""`
	Debug                  bool   `json:"debug" mapstructure:"debug" default:"false"`
	Port                   uint32 `json:"port" mapstructure:"port" default:"9000"`
	HttpPort               uint32 `json:"http_port" mapstructure:"http_port" default:"8123"`
	ReadTimeout            uint32 `json:"read_timeout" mapstructure:"read_timeout" default:"30"`
	WriteTimeout           uint32 `json:"write_timeout" mapstructure:"write_timeout" default:"30"`
	MaxIdleConn            int    `json:"max_idle_connection" mapstructure:"max_idle_connection" default:"5"`
	MaxOpenConn            int    `json:"max_open_connection" mapstructure:"max_open_connection" default:"50"`
	Primary                bool   `json:"primary" mapstructure:"primary" default:"false"`
	Strategy               string `json:"strategy" mapstructure:"strategy" default:"failover"`
	TTLDays                int    `json:"ttl_days" mapstructure:"ttl_days" default:"0"`
	StoragePolicy          string `json:"storage_policy" mapstructure:"storage_policy" default:""`
	Secure                 bool   `json:"secure" mapstructure:"secure" default:"false"`
	Cloud                  bool   `json:"cloud" mapstructure:"cloud" default:"false"`
	ClusterName            string `json:"cluster_name" mapstructure:"cluster_name" default:""`
	AsyncInsert            bool   `json:"async_insert" mapstructure:"async_insert" default:"false"`
	EmergencySweepLimit    string `json:"emergency_sweep_limit" mapstructure:"emergency_sweep_limit" default:""`
	EmergencySweepAdvanced string `json:"emergency_sweep_advanced" mapstructure:"emergency_sweep_advanced" default:""`
	TextCodec              string `json:"text_codec" mapstructure:"text_codec" default:""`
	LogsIndex              string `json:"logs_index" mapstructure:"logs_index" default:""`
	LogsIndexGranularity   uint32 `json:"logs_index_granularity" mapstructure:"logs_index_granularity" default:""`
	ReplicatedClusterName  string `json:"replicated_cluster_name" mapstructure:"replicated_cluster_name" default:""`
	TestDistributed        bool   `json:"test_distributed" mapstructure:"test_distributed" default:"false"`
	Distributed            bool   `json:"distributed" mapstructure:"distributed" default:"true"`
	SamplesOrdering        string `json:"samples_ordering" mapstructure:"samples_ordering" default:"false"`
	SkipUnavailableShards  bool   `json:"skip_unavailable_shards" mapstructure:"skip_unavailable_shards" default:"false"`
	InsecureSkipVerify     bool   `json:"insecure_skip_verify" mapstructure:"insecure_skip_verify" default:"false"`

	TTLPolicy []struct {
		Timeout string `json:"ttl_policy" mapstructure:"ttl_policy" default:""`
		MoveTo  string `json:"move_to" mapstructure:"move_to" default:""`
	} `json:"ttl_policy" mapstructure:"ttl_policy" default:""`
}

type ClokiBaseSettingServer struct {
	ClokiWriter writer.ClokiWriterSettingServer `json:"writer" mapstructure:"writer"`
	ClokiReader reader.ClokiReaderSettingServer `json:"reader" mapstructure:"reader"`

	//Base
	FingerPrintType          uint `default:"1"`
	DataBaseStrategy         uint `default:"0"`
	CurrentDataBaseIndex     uint `default:"0"`
	DataDatabaseGroupNodeMap map[string][]string
	Validate                 *validator.Validate
	EnvPrefix                string `default:"QRYN"`
	PluginsPath              string `default:""`
	AnalyticsDatabase        string `json:"analytics_database" mapstructure:"analytics_database" default:""`
	//clickhouse://default:pass4default@127.0.0.1:8123/analytics

	DATABASE_DATA []ClokiBaseDataBase `json:"database_data" mapstructure:"database_data"`

	CLUSTER_SETTINGS struct {
		DistributionHeader bool `json:"distribution_header" mapstructure:"distribution_header" default:"false"`
	} `json:"cluster_settings" mapstructure:"cluster_settings"`

	DRILLDOWN_SETTINGS struct {
		LogDrilldown            bool    `json:"log_drilldown" mapstructure:"log_drilldown" default:"false"`
		LogPatternsSimilarity   float64 `json:"log_patterns_similarity" mapstructure:"pattern_similarity" default:"0.7"`
		LogPatternsReadLimit    int     `json:"log_patterns_read_limit" mapstructure:"patterns_read_limit" default:"300"`
		LogPatternsDownsampling float64 `json:"log_patterns_downsampling" mapstructure:"patterns_downsampling" default:"1"`
	} `json:"drilldown_settings" mapstructure:"drilldown_settings"`

	SYSTEM_SETTINGS struct {
		HostName                    string  `json:"hostname" mapstructure:"hostname" default:"hostname"`
		EnableUserAuditLogin        bool    `json:"user_audit_login" mapstructure:"user_audit_login" default:"true"`
		UUID                        string  `json:"uuid" mapstructure:"uuid" default:""`
		DBBulk                      int64   `json:"db_bulk" mapstructure:"db_bulk" default:"0"`
		DBTimer                     float64 `json:"db_timer" mapstructure:"db_timer" default:"0.2"`
		RetryAttempts               int     `json:"retry_attempts" mapstructure:"retry_attempts" default:"10"`
		RetryTimeoutS               int     `json:"retry_timeout_s" mapstructure:"retry_timeout_s" default:"1"`
		BufferSizeSample            uint32  `json:"buffer_size_sample" mapstructure:"buffer_size_sample" default:"200000"`
		BufferSizeTimeSeries        uint32  `json:"buffer_size_timeseries" mapstructure:"buffer_size_timeseries" default:"200000"`
		ChannelsSample              int     `json:"channels_sample" mapstructure:"channels_sample" default:"2"`
		ChannelsTimeSeries          int     `json:"channels_timeseries" mapstructure:"channels_timeseries" default:"2"`
		HashType                    string  `json:"hash_type" mapstructure:"hash_type" default:""`
		CPUMaxProcs                 int     `json:"cpu_max_procs" mapstructure:"cpu_max_procs" default:"0"`
		NoForceRotation             bool    `json:"no_force_rotation" mapstructure:"no_force_rotation" default:"false"`
		QueryStats                  bool    `json:"query_stats" mapstructure:"query_stats" default:"false"`
		DynamicDatabases            bool    `json:"dynamic_databases" mapstructure:"dynamic_databases" default:"false"`
		DynamicDatabasesReadTimeout float64 `json:"dynamic_databases_read_timeout" mapstructure:"dynamic_databases_read_timeout" default:"0"`
		AWSLambda                   bool    `json:"aws_lambda" mapstructure:"aws_lambda" default:"false"`
		LicenseAutoShutdown         bool    `json:"license_auto_shutdown" mapstructure:"license_auto_shutdown" default:"false"`
		DynamicFolder               string  `json:"dynamic_folder" mapstructure:"dynamic_folder" default:""`
		MetricsMaxSamples           int     `json:"metrics_max_samples" mapstructure:"metrics_max_samples" default:"5000000"`
		MaxSeries                   int     `json:"max_series" mapstructure:"max_series" default:"0"`
		Mode                        string  `json:"mode" mapstructure:"mode" default:""`
		TotalRateLimitMB            int     `json:"total_rate_limit_mb" mapstructure:"total_rate_limit_mb" default:"80"`
		RateLimitPerDBMB            int     `json:"rate_limit_per_db_mb" mapstructure:"rate_limit_per_db_mb" default:"50"`
		MaxParallelQueries          int     `json:"max_parallel_queries" mapstructure:"max_parallel_queries" default:"0"`
		PyroscopeServerURL          string  `json:"pyroscope_server_url" mapstructure:"pyroscope_server_url" default:""`
		PyroscopeExtraTags          string  `json:"pyroscope_extra_tags" mapstructure:"pyroscope_extra_tags"`
	} `json:"system_settings" mapstructure:"system_settings"`

	WORKER struct {
		Type                    string `json:"type" mapstructure:"type" default:""`
		SyncUrl                 string `json:"sync_url" mapstructure:"sync_url" default:""`
		MQUrl                   string `json:"mq_url" mapstructure:"mq_url" default:""`
		AwsLambdaARN            string `json:"aws_lambda_arn" mapstructure:"aws_lambda_arn" default:""`
		AlertManagerURL         string `json:"alert_manager_url" mapstructure:"alert_manager_url" default:""`
		RecordingRulesWriterURL string `json:"recording_rules_writer" mapstructure:"recording_rules_writer" default:""`
	} `json:"worker" mapstructure:"worker"`

	MULTITENANCE_SETTINGS struct {
		Enabled bool `json:"enabled" mapstructure:"enabled" default:"false"`
	} `json:"multitenance_settings" mapstructure:"multitenance_settings"`

	AUTH_SETTINGS struct {
		AuthTokenHeader string `json:"token_header" mapstructure:"token_header" default:"Auth-Token"`
		AuthTokenExpire uint32 `json:"token_expire" mapstructure:"token_expire" default:"1200"`
		BASIC           struct {
			Username string `json:"username" mapstructure:"username" default:""`
			Password string `json:"password" mapstructure:"password" default:""`
		}
	} `json:"auth_settings" mapstructure:"auth_settings"`

	API_SETTINGS struct {
		EnableForceSync   bool `json:"sync_map_force" mapstructure:"sync_map_force" default:"false"`
		EnableTokenAccess bool `json:"enable_token_access" mapstructure:"enable_token_access" default:"true"`
	} `json:"api_settings" mapstructure:"api_settings"`

	SCRIPT_SETTINGS struct {
		Enable bool   `json:"enable" mapstructure:"enable" default:"false"`
		Engine string `json:"engine" mapstructure:"engine" default:"lua"`
		Folder string `json:"folder" mapstructure:"folder" default:"/usr/local/qryn/scripts/"`
	} `json:"script_settings" mapstructure:"script_settings"`

	FORWARD_SETTINGS struct {
		ForwardUrl     string `json:"forward_url" mapstructure:"forward_url" default:""`
		ForwardLabels  string `json:"forward_labels" mapstructure:"forward_labels" default:""`
		ForwardHeaders string `json:"forward_headers" mapstructure:"forward_headers" default:""`
	} `json:"forward_settings" mapstructure:"forward_settings"`

	HTTP_SETTINGS struct {
		Host          string `json:"host" mapstructure:"host" default:"0.0.0.0"`
		Port          int    `json:"port" mapstructure:"port" default:"0"`
		ApiPrefix     string `json:"api_prefix" mapstructure:"api_prefix" default:""`
		ApiPromPrefix string `json:"api_prom_prefix" mapstructure:"api_prom_prefix" default:""`
		Prefork       bool   `json:"prefork" mapstructure:"prefork" default:"false"`
		Gzip          bool   `json:"gzip" mapstructure:"gzip" default:"true"`
		GzipStatic    bool   `json:"gzip_static" mapstructure:"gzip_static" default:"true"`
		Debug         bool   `json:"debug" mapstructure:"debug" default:"false"`
		Concurrency   int    `json:"concurrency" mapstructure:"concurrency" default:"350"`
		Cors          struct {
			Origin string `json:"origin" mapstructure:"origin" default:"*"`
			Enable bool   `json:"enable" mapstructure:"enable" default:"false"`
		} `json:"cors" mapstructure:"cors"`
		WebSocket struct {
			Enable bool `json:"enable" mapstructure:"enable" default:"false"`
		} `json:"websocket" mapstructure:"websocket"`
		Enable        bool `json:"enable" mapstructure:"enable" default:"true"`
		InputBufferMB int  `json:"input_buffer_mb" mapstructure:"input_buffer_mb" default:"200"`
	} `json:"http_settings" mapstructure:"http_settings"`

	HTTPS_SETTINGS struct {
		Host                string `json:"host" mapstructure:"host" default:"0.0.0.0"`
		Port                int    `json:"port" mapstructure:"port" default:"3201"`
		Cert                string `json:"cert" mapstructure:"cert" default:""`
		Key                 string `json:"key" mapstructure:"key" default:""`
		HttpRedirect        bool   `json:"http_redirect" mapstructure:"http_redirect" default:"false"`
		Enable              bool   `json:"enable" mapstructure:"enable" default:"false"`
		MinTLSVersionString string `json:"min_tls_version" mapstructure:"min_tls_version" default:"0"`
		MaxTLSVersionString string `json:"max_tls_version" mapstructure:"max_tls_version" default:"0"`
		MinTLSVersion       uint16 `default:"0"`
		MaxTLSVersion       uint16 `default:"0"`
	} `json:"https_settings" mapstructure:"https_settings"`

	LOG_SETTINGS struct {
		Enable        bool   `json:"enable" mapstructure:"enable" default:"true"`
		MaxAgeDays    uint32 `json:"max_age_days" mapstructure:"max_age_days" default:"7"`
		RotationHours uint32 `json:"rotation_hours" mapstructure:"rotation_hours" default:"24"`
		Path          string `json:"path" mapstructure:"path" default:"./"`
		Level         string `json:"level" mapstructure:"level" default:"info"`
		Name          string `json:"name" mapstructure:"name" default:"ClokiBase.log"`
		Stdout        bool   `json:"stdout" mapstructure:"stdout" default:"false"`
		Json          bool   `json:"json" mapstructure:"json" default:"true"`
		SysLogLevel   string `json:"syslog_level" mapstructure:"syslog_level" default:"LOG_INFO"`
		SysLog        bool   `json:"syslog" mapstructure:"syslog" default:"false"`
		SyslogUri     string `json:"syslog_uri" mapstructure:"syslog_uri" default:""`
		TracesUrl     string `json:"traces_url" mapstructure:"traces_url" default:""`
		InstanceName  string `json:"instance_name" mapstructure:"instance_name" default:""`
		Qryn          struct {
			Url         string `json:"url" mapstructure:"url" default:""`
			App         string `json:"app" mapstructure:"app" default:""`
			AddHostname bool   `json:"add_hostname" mapstructure:"add_hostname" default:"false"`
			Headers     string `json:"headers" mapstructure:"headers" default:""`
		} `json:"qryn" mapstructure:"qryn"`
	} `json:"log_settings" mapstructure:"log_settings"`

	PROMETHEUS_CLIENT struct {
		PushURL      string   `json:"push_url" mapstructure:"push_url" default:""`
		TargetIP     string   `json:"target_ip" mapstructure:"target_ip" default:""`
		PushInterval uint32   `json:"push_interval" mapstructure:"push_interval" default:"60"`
		PushName     string   `json:"push_name" mapstructure:"push_name" default:""`
		ServiceName  string   `json:"service_name" mapstructure:"service_name" default:"prometheus"`
		MetricsPath  string   `json:"metrics_path" mapstructure:"metrics_path" default:"/metrics"`
		Enable       bool     `json:"enable" mapstructure:"enable" default:"false"`
		AllowIP      []string `json:"allow_ip" mapstructure:"allow_ip" default:"[127.0.0.1]"`
	} `json:"prometheus_client" mapstructure:"prometheus_client"`

	LICENSE_SETTINGS struct {
		ProxyServer              string `json:"proxy_server" mapstructure:"proxy_server" default:""`
		RemoteCheckInterval      string `json:"remote_check_interval" mapstructure:"remote_check_interval" default:"90d"`
		GracePeriodCheckInterval string `json:"grace_check_interval" mapstructure:"grace_check_interval" default:"14d"`
	} `json:"license_settings" mapstructure:"license_settings"`

	WEBHOOKS struct {
		LimitYellowZone string `json:"limit_yellow_zone" mapstructure:"limit_yellow_zone" default:""`
		LimitRedZone    string `json:"limit_red_zone" mapstructure:"limit_red_zone" default:""`
	} `json:"webhooks" mapstructure:"webhooks"`

	EXPORTER_SETTINGS struct {
		ServerEnable bool   `json:"server_enable" mapstructure:"server_enable" default:"false"`
		ExportEnable bool   `json:"export_enable" mapstructure:"export_enable" default:"false"`
		Cron         string `json:"cron" mapstructure:"cron" default:""`
		From         string `json:"from" mapstructure:"from" default:""`
		To           string `json:"to" mapstructure:"to" default:""`
	} `json:"exporter_settings" mapstructure:"exporter_settings"`
}

func (c ClokiBaseSettingServer) PyroscopeExtraTags() map[string]string {
	res := make(map[string]string)
	kvs := strings.Split(c.SYSTEM_SETTINGS.PyroscopeExtraTags, ";")
	for _, kv := range kvs {
		_kv := strings.SplitN(kv, "=", 2)
		if len(_kv) != 2 {
			continue
		}
		_kv[0] = strings.TrimSpace(_kv[0])
		_kv[1] = strings.TrimSpace(_kv[1])
		if _kv[0] != "" && _kv[1] != "" {
			res[_kv[0]] = _kv[1]
		}
	}
	return res
}
