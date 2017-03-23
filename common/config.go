package common

import (
	"os"
	"path/filepath"
)

const (
	// DomainName specifies domain name in Contrail
	DomainName = "default-domain"

	// DriverName is name of the driver that is to be specified during docker network creation
	DriverName = "Contrail"

	// HNSNetworkPrefix is a prefix given too all HNS network names managed by the driver
	HNSNetworkPrefix = "Contrail"

	// RootNetworkName is a name of root HNS network created solely for the purpose of
	// having a virtual switch
	RootNetworkName = "ContrailRootNetwork"

	// AdapterReconnectTimeout is a time (in ms) to wait for adapter to reacquire IP after a new
	// HNS network is created. https://github.com/Microsoft/hcsshim/issues/108
	AdapterReconnectTimeout = 15000

	// AdapterPollingRate is rate (in ms) of polling of network adapter while waiting for it to
	// reacquire IP.
	AdapterPollingRate = 300

	// HNSTransparentInterfaceName is the name of transparent HNS vswitch interface name
	HNSTransparentInterfaceName = "vEthernet (HNSTransparent)"
)

// PluginSpecDir returns path to directory where docker daemon looks for plugin spec files.
func PluginSpecDir() string {
	return filepath.Join(os.Getenv("programdata"), "docker", "plugins")
}

// PluginSpecFilePath returns path to plugin spec file.
func PluginSpecFilePath() string {
	return filepath.Join(PluginSpecDir(), DriverName+".spec")
}
