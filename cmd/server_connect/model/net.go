package model

// NetConfig info about listen and connect restriction
type NetConfig struct {
	TCPPort              int    `ini:"TCP_PORT"`
	UDPPort              int    `ini:"UDP_PORT"`
	MaxConnectionsPerIP  int    `ini:"MaxConnectionsPerIP"`
	MaxPacketsPerSecond  int    `ini:"MaxPacketsPerSecond"`
	LauncherProxyWhiteIP string `ini:"LauncherProxyWhiteListIP"`
}
