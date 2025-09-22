package reader

// ============================ READER ONLY ================================= //
type ClokiReaderSettingServer struct {
	ViewPath        string `json:"view_path" mapstructure:"view_path" default:"/etc/qryn-view"`
	WriterProxy     string `json:"writer_proxy" mapstructure:"writer_proxy" default:""`
	OmitEmptyValues bool   `json:"omit_empty_values" mapstructure:"omit_empty_values" default:"false"`
}
