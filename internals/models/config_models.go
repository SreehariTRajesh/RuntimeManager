package models

type Config struct {
	HostConf           HostConfig           `yaml:"host-config"`
	VXLanConf          VXLanConfig          `yaml:"vxlan-config"`
	MacVLANNetworkConf MacVLANNetworkConfig `yaml:"macvlan-network-config"`
}

type HostConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type VXLanConfig struct {
	Bridge     string      `yaml:"bridge"`
	VXLanPeers []VXLanPeer `yaml:"vxlan-peers"`
}

type VXLanPeer struct {
	Name    string `yaml:"name"`
	VXLanId int    `yaml:"vxlan-id"`
	Remote  string `yaml:"remote"`
	DstPort int    `yaml:"dst-port"`
	Device  string `yaml:"device"`
}

type MacVLANNetworkConfig struct {
	Name    string `yaml:"name"`
	Subnet  string `yaml:"subnet"`
	Gateway string `yaml:"gateway"`
	Parent  string `yaml:"parent"`
}
