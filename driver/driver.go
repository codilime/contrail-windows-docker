// Implemented according to
// https://github.com/docker/libnetwork/blob/master/docs/remote.md

package driver

import (
	"encoding/json"
	"fmt"

	"github.com/Microsoft/hcsshim"
	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/network"
)

const (
	// Drivername is name of the driver that is to be specified during docker network creation
	Drivername = "Contrail"
)

type ContrailDriver struct{}

func (d *ContrailDriver) GetCapabilities() (*network.CapabilitiesResponse, error) {
	logrus.Println("=== GetCapabilities")
	r := &network.CapabilitiesResponse{}
	r.Scope = network.LocalScope
	return r, nil
}

func (d *ContrailDriver) CreateNetwork(req *network.CreateNetworkRequest) error {
	logrus.Println("=== CreateNetwork")
	logrus.Println("network.NetworkID =", req.NetworkID)
	logrus.Println(req)
	logrus.Println("IPv4:")
	for _, n := range req.IPv4Data {
		logrus.Println(n)
	}
	logrus.Println("IPv6:")
	for _, n := range req.IPv6Data {
		logrus.Println(n)
	}
	logrus.Println("options:")
	for k, v := range req.Options {
		fmt.Printf("%v: %v\n", k, v)
	}

	subnets := []hcsshim.Subnet{}
	s := hcsshim.Subnet{
		AddressPrefix:  req.IPv4Data[0].Pool,
		GatewayAddress: req.IPv4Data[0].Gateway,
	}
	subnets = append(subnets, s)

	logrus.Println("subnets", subnets)

	configuration := &hcsshim.HNSNetwork{
		Name:    req.NetworkID,
		Type:    "transparent",
		Subnets: subnets,
	}

	request, err := json.Marshal(configuration)
	if err != nil {
		return err
	}
	logrus.Println("[HNS] Request ", string(request))

	response, err := hcsshim.HNSNetworkRequest("POST", "", string(request))
	if err != nil {
		logrus.Println("[HNS] Error ", err)
		return err
	}
	logrus.Println("[HNS] Response ", response)

	return nil
}

func (d *ContrailDriver) AllocateNetwork(req *network.AllocateNetworkRequest) (*network.AllocateNetworkResponse, error) {
	logrus.Println("=== AllocateNetwork")
	logrus.Println(req)
	r := &network.AllocateNetworkResponse{}
	return r, nil
}

func (d *ContrailDriver) DeleteNetwork(req *network.DeleteNetworkRequest) error {
	logrus.Println("=== DeleteNetwork")
	logrus.Println(req)
	return nil
}

func (d *ContrailDriver) FreeNetwork(req *network.FreeNetworkRequest) error {
	logrus.Println("=== FreeNetwork")
	logrus.Println(req)
	return nil
}

func (d *ContrailDriver) CreateEndpoint(req *network.CreateEndpointRequest) (*network.CreateEndpointResponse, error) {
	logrus.Println("=== CreateEndpoint")
	logrus.Println(req)
	logrus.Println(req.Interface)
	logrus.Println("options:")
	for k, v := range req.Options {
		fmt.Printf("%v: %v\n", k, v)
	}
	r := &network.CreateEndpointResponse{}
	return r, nil
}

func (d *ContrailDriver) DeleteEndpoint(req *network.DeleteEndpointRequest) error {
	logrus.Println("=== DeleteEndpoint")
	logrus.Println(req)
	return nil
}

func (d *ContrailDriver) EndpointInfo(req *network.InfoRequest) (*network.InfoResponse, error) {
	logrus.Println("=== EndpointInfo")
	logrus.Println(req)
	r := &network.InfoResponse{}
	return r, nil
}

func (d *ContrailDriver) Join(req *network.JoinRequest) (*network.JoinResponse, error) {
	logrus.Println("=== Join")
	logrus.Println(req)
	logrus.Println("options:")
	for k, v := range req.Options {
		fmt.Printf("%v: %v\n", k, v)
	}
	r := &network.JoinResponse{}
	r.DisableGatewayService = true
	return r, nil
}

func (d *ContrailDriver) Leave(req *network.LeaveRequest) error {
	logrus.Println("=== Leave")
	logrus.Println(req)
	return nil
}

func (d *ContrailDriver) DiscoverNew(req *network.DiscoveryNotification) error {
	logrus.Println("=== DiscoverNew")
	logrus.Println(req)
	return nil
}

func (d *ContrailDriver) DiscoverDelete(req *network.DiscoveryNotification) error {
	logrus.Println("=== DiscoverDelete")
	logrus.Println(req)
	return nil
}

func (d *ContrailDriver) ProgramExternalConnectivity(req *network.ProgramExternalConnectivityRequest) error {
	logrus.Println("=== ProgramExternalConnectivity")
	logrus.Println(req)
	return nil
}

func (d *ContrailDriver) RevokeExternalConnectivity(req *network.RevokeExternalConnectivityRequest) error {
	logrus.Println("=== RevokeExternalConnectivity")
	logrus.Println(req)
	return nil
}
