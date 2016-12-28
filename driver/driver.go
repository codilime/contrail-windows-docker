// Implemented according to
// https://github.com/docker/libnetwork/blob/master/docs/remote.md

package driver

import (
	"fmt"

	"github.com/Microsoft/hcsshim"
	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/network"
	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	// DriverName is name of the driver that is to be specified during docker network creation
	DriverName = "Contrail"

	// HNSNetworkName is a constant name for HNS network
	NetworkHNSname = "ContrailHNSNet"
)

type ContrailDriver struct {
	HnsID string
}

func NewDriver(subnet, gateway, adapter string) (*ContrailDriver, error) {

	subnets := []hcsshim.Subnet{
		{
			AddressPrefix:  subnet,
			GatewayAddress: gateway,
		},
	}

	configuration := &hcsshim.HNSNetwork{
		Name:               NetworkHNSname,
		Type:               "transparent",
		Subnets:            subnets,
		NetworkAdapterName: adapter,
	}

	hnsID, err := CreateHNSNetwork(configuration)
	if err != nil {
		return nil, err
	}

	d := &ContrailDriver{
		HnsID: hnsID,
	}
	return d, nil
}

func (d *ContrailDriver) Serve() error {
	h := network.NewHandler(d)

	config := sdk.WindowsPipeConfig{
		// This will set permissions for Everyone user allowing him to open, write, read the pipe
		SecurityDescriptor: "S:(ML;;NW;;;LW)D:(A;;0x12019f;;;WD)",
		InBufferSize:       4096,
		OutBufferSize:      4096,
	}

	h.ServeWindows("//./pipe/"+DriverName, DriverName, &config)
	return nil
}

func (d *ContrailDriver) Teardown() error {
	err := DeleteHNSNetwork(d.HnsID)
	return err
}

func (d *ContrailDriver) GetCapabilities() (*network.CapabilitiesResponse, error) {
	logrus.Debugln("=== GetCapabilities")
	r := &network.CapabilitiesResponse{}
	r.Scope = network.LocalScope
	return r, nil
}

func (d *ContrailDriver) CreateNetwork(req *network.CreateNetworkRequest) error {
	logrus.Debugln("=== CreateNetwork")
	logrus.Debugln("network.NetworkID =", req.NetworkID)
	logrus.Debugln(req)
	logrus.Debugln("IPv4:")
	for _, n := range req.IPv4Data {
		logrus.Debugln(n)
	}
	logrus.Debugln("IPv6:")
	for _, n := range req.IPv6Data {
		logrus.Debugln(n)
	}
	logrus.Debugln("options:")
	for k, v := range req.Options {
		fmt.Printf("%v: %v\n", k, v)
	}

	// subnets := []hcsshim.Subnet{}
	// s := hcsshim.Subnet{
	// 	AddressPrefix:  req.IPv4Data[0].Pool,
	// 	GatewayAddress: req.IPv4Data[0].Gateway,
	// }
	// subnets = append(subnets, s)

	// logrus.Debugln("subnets", subnets)

	// configuration := &hcsshim.HNSNetwork{
	// 	Name:    req.NetworkID,

	return nil
}

func (d *ContrailDriver) AllocateNetwork(req *network.AllocateNetworkRequest) (*network.AllocateNetworkResponse, error) {
	logrus.Debugln("=== AllocateNetwork")
	logrus.Debugln(req)
	r := &network.AllocateNetworkResponse{}
	return r, nil
}

func (d *ContrailDriver) DeleteNetwork(req *network.DeleteNetworkRequest) error {
	logrus.Debugln("=== DeleteNetwork")
	logrus.Debugln(req)
	return nil
}

func (d *ContrailDriver) FreeNetwork(req *network.FreeNetworkRequest) error {
	logrus.Debugln("=== FreeNetwork")
	logrus.Debugln(req)
	return nil
}

func (d *ContrailDriver) CreateEndpoint(req *network.CreateEndpointRequest) (*network.CreateEndpointResponse, error) {
	logrus.Debugln("=== CreateEndpoint")
	logrus.Debugln(req)
	logrus.Debugln(req.Interface)
	logrus.Debugln("options:")
	for k, v := range req.Options {
		fmt.Printf("%v: %v\n", k, v)
	}
	r := &network.CreateEndpointResponse{}
	return r, nil
}

func (d *ContrailDriver) DeleteEndpoint(req *network.DeleteEndpointRequest) error {
	logrus.Debugln("=== DeleteEndpoint")
	logrus.Debugln(req)
	return nil
}

func (d *ContrailDriver) EndpointInfo(req *network.InfoRequest) (*network.InfoResponse, error) {
	logrus.Debugln("=== EndpointInfo")
	logrus.Debugln(req)
	r := &network.InfoResponse{}
	return r, nil
}

func (d *ContrailDriver) Join(req *network.JoinRequest) (*network.JoinResponse, error) {
	logrus.Debugln("=== Join")
	logrus.Debugln(req)
	logrus.Debugln("options:")
	for k, v := range req.Options {
		fmt.Printf("%v: %v\n", k, v)
	}
	r := &network.JoinResponse{}
	r.DisableGatewayService = true
	return r, nil
}

func (d *ContrailDriver) Leave(req *network.LeaveRequest) error {
	logrus.Debugln("=== Leave")
	logrus.Debugln(req)
	return nil
}

func (d *ContrailDriver) DiscoverNew(req *network.DiscoveryNotification) error {
	logrus.Debugln("=== DiscoverNew")
	logrus.Debugln(req)
	return nil
}

func (d *ContrailDriver) DiscoverDelete(req *network.DiscoveryNotification) error {
	logrus.Debugln("=== DiscoverDelete")
	logrus.Debugln(req)
	return nil
}

func (d *ContrailDriver) ProgramExternalConnectivity(req *network.ProgramExternalConnectivityRequest) error {
	logrus.Debugln("=== ProgramExternalConnectivity")
	logrus.Debugln(req)
	return nil
}

func (d *ContrailDriver) RevokeExternalConnectivity(req *network.RevokeExternalConnectivityRequest) error {
	logrus.Debugln("=== RevokeExternalConnectivity")
	logrus.Debugln(req)
	return nil
}
