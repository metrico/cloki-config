package reader

// ============================ READER ONLY ================================= //
type ClokiReaderSettingServer struct {
	ViewPath    string `json:"view_path" mapstructure:"view_path" default:""`
	WriterProxy string `json:"writer_proxy" mapstructure:"writer_proxy" default:""`
}
