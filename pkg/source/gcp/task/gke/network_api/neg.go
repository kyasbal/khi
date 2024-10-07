package network_api

type NegAttachOrDetachRequestEndpoint struct {
	Instance  string `yaml:"instance"`
	IpAddress string `yaml:"ipAddress"`
	Port      string `yaml:"port"`
}

type NegAttachOrDetachRequest struct {
	NetworkEndpoints []*NegAttachOrDetachRequestEndpoint `yaml:"networkEndpoints"`
}
