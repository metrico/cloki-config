package deleter

// ============================ DELETER ONLY ================================ //

// Provider names for ClokiDeleterSettingServer.Provider.
const (
	// PROVIDER_CONFIG takes the single ClickHouse server from database_data.
	PROVIDER_CONFIG = "config"
	// PROVIDER_HTTP asks an HTTP endpoint which servers to sweep.
	PROVIDER_HTTP = "http"
)

// ClokiDeleterSettingServer configures the deleter daemon, which sweeps the
// pending delete requests of every ClickHouse server it is given.
//
// Interval and Headers are strings rather than a time.Duration and a map: the
// env binding walks the struct and binds every non-struct, non-slice field to a
// single environment variable, so both have to survive a round trip through one
// string. They are parsed by the consumer.
type ClokiDeleterSettingServer struct {
	// Interval is the pause between two sweeps, as a Go duration. "0" runs a
	// single sweep and exits.
	Interval string `json:"interval" mapstructure:"interval" default:"1m"`
	// Provider is where the list of ClickHouse servers comes from, either
	// PROVIDER_CONFIG or PROVIDER_HTTP.
	Provider string `json:"provider" mapstructure:"provider" default:"config"`
	// ProviderURL is the endpoint listing the servers, for PROVIDER_HTTP. It may
	// carry basic auth credentials, which ProviderUser overrides.
	ProviderURL string `json:"provider_url" mapstructure:"provider_url" default:""`
	// ProviderUser and ProviderPassword are the basic auth of the endpoint.
	ProviderUser     string `json:"provider_user" mapstructure:"provider_user" default:""`
	ProviderPassword string `json:"provider_password" mapstructure:"provider_password" default:""`
	// ProviderHeaders are extra headers for the endpoint, as a JSON object:
	// {"X-Api-Key": "secret"}. For endpoints behind something basic auth cannot
	// express - an API token, a tenant selector, a proxy keyed off a header.
	ProviderHeaders string `json:"provider_headers" mapstructure:"provider_headers" default:""`
	// Concurrency is how many servers are swept in parallel.
	Concurrency int `json:"concurrency" mapstructure:"concurrency" default:"10"`
}
