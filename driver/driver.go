// Implemented according to
// https://github.com/docker/libnetwork/blob/master/docs/remote.md

package driver

import (
	"errors"
	"fmt"

	"context"

	"github.com/Microsoft/hcsshim"
	log "github.com/Sirupsen/logrus"
	"github.com/codilime/contrail-windows-docker/common"
	"github.com/codilime/contrail-windows-docker/controller"
	"github.com/codilime/contrail-windows-docker/hns"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/go-plugins-helpers/network"
	"github.com/docker/go-plugins-helpers/sdk"
)

type ContrailDriver struct {
	controller *controller.Controller
	HnsID      string
}

func NewDriver(subnet, gateway, adapter string, controller *controller.Controller) (*ContrailDriver,
	error) {

	subnets := []hcsshim.Subnet{
		{
			AddressPrefix:  subnet,
			GatewayAddress: gateway,
		},
	}

	configuration := &hcsshim.HNSNetwork{
		Name:               common.NetworkHNSname,
		Type:               "transparent",
		Subnets:            subnets,
		NetworkAdapterName: adapter,
	}

	hnsID, err := hns.CreateHNSNetwork(configuration)
	if err != nil {
		return nil, err
	}

	d := &ContrailDriver{
		controller: controller,
		HnsID:      hnsID,
	}
	return d, nil
}

func (d *ContrailDriver) Serve() error {
	h := network.NewHandler(d)

	config := sdk.WindowsPipeConfig{
		// This will set permissions for Service, System, Adminstrator group and account to have full access
		SecurityDescriptor: "D:(A;ID;FA;;;SY)(A;ID;FA;;;BA)(A;ID;FA;;;LA)(A;ID;FA;;;LS)",

		InBufferSize:  4096,
		OutBufferSize: 4096,
	}

	h.ServeWindows("//./pipe/"+common.DriverName, common.DriverName, &config)
	return nil
}

func (d *ContrailDriver) Teardown() error {
	err := hns.DeleteHNSNetwork(d.HnsID)
	return err
}

func (d *ContrailDriver) GetCapabilities() (*network.CapabilitiesResponse, error) {
	log.Debugln("=== GetCapabilities")
	r := &network.CapabilitiesResponse{}
	r.Scope = network.LocalScope
	return r, nil
}

func (d *ContrailDriver) CreateNetwork(req *network.CreateNetworkRequest) error {
	log.Debugln("=== CreateNetwork")
	log.Debugln("network.NetworkID =", req.NetworkID)
	log.Debugln(req)
	log.Debugln("IPv4:")
	for _, n := range req.IPv4Data {
		log.Debugln(n)
	}
	log.Debugln("IPv6:")
	for _, n := range req.IPv6Data {
		log.Debugln(n)
	}
	log.Debugln("options:")
	for k, v := range req.Options {
		fmt.Printf("%v: %v\n", k, v)
	}

	tenant, exists := req.Options["tenant"]
	if !exists {
		return errors.New("Tenant not specified")
	}

	netName, exists := req.Options["network"]
	if !exists {
		return errors.New("Network name not specified")
	}

	_, err := d.controller.GetNetwork(tenant.(string), netName.(string))
	if err != nil {
		return err
	}
	return nil

}

func (d *ContrailDriver) AllocateNetwork(req *network.AllocateNetworkRequest) (*network.AllocateNetworkResponse, error) {
	log.Debugln("=== AllocateNetwork")
	log.Debugln(req)
	r := &network.AllocateNetworkResponse{}
	return r, nil
}

func (d *ContrailDriver) DeleteNetwork(req *network.DeleteNetworkRequest) error {
	log.Debugln("=== DeleteNetwork")
	log.Debugln(req)
	return nil
}

func (d *ContrailDriver) FreeNetwork(req *network.FreeNetworkRequest) error {
	log.Debugln("=== FreeNetwork")
	log.Debugln(req)
	return nil
}

func (d *ContrailDriver) CreateEndpoint(req *network.CreateEndpointRequest) (*network.CreateEndpointResponse, error) {
	log.Debugln("=== CreateEndpoint")
	log.Debugln(req)
	log.Debugln(req.Interface)
	log.Debugln("options:")
	for k, v := range req.Options {
		fmt.Printf("%v: %v\n", k, v)
	}

	docker, err := dockerClient.NewEnvClient()
	if err != nil {
		return nil, err
	}

	dockerNetwork, err := docker.NetworkInspect(context.Background(), req.NetworkID)
	if err != nil {
		return nil, err
	}

	tenant, exists := dockerNetwork.Options["tenant"]
	if !exists {
		return nil, errors.New("Retreived network has no Contrail tenant specified")
	}

	netName, exists := dockerNetwork.Options["network"]
	if !exists {
		return nil, errors.New("Retreived network has no Contrail network name specfied")
	}

	log.Infoln(tenant, netName)

	if contrailNetwork, err := d.controller.GetNetwork(tenant, netName)
	if err != nil {
		return nil, err
	}

	contrailInstance, err := d.controller.GetOrCreateInstance(tenant, req.EndpointID)
	if err != nil {
		return nil, err
	}

	contrailVif, err := d.controller.GetOrCreateInterface(contrailNetwork, contrailInstance)
	if err != nil {
		return nil, err
	}

	contrailIP, err := d.controller.GetOrCreateInstanceIp(contrailNetwork, contrailVif)
	if err != nil {
		return nil, err
	}

	contrailGateway, err := d.controller.GetDefaultGatewayIp(contrailNetwork)
	if err != nil {
		return nil, err
	}

	contrailMac, err := d.controller.GetInterfaceMac(contrailVif)
	if err != nil {
		return nil, err
	}

	r := &network.CreateEndpointResponse{}
	return r, nil
}

func (d *ContrailDriver) DeleteEndpoint(req *network.DeleteEndpointRequest) error {
	log.Debugln("=== DeleteEndpoint")
	log.Debugln(req)
	return nil
}

func (d *ContrailDriver) EndpointInfo(req *network.InfoRequest) (*network.InfoResponse, error) {
	log.Debugln("=== EndpointInfo")
	log.Debugln(req)
	r := &network.InfoResponse{}
	return r, nil
}

func (d *ContrailDriver) Join(req *network.JoinRequest) (*network.JoinResponse, error) {
	log.Debugln("=== Join")
	log.Debugln(req)
	log.Debugln("options:")
	for k, v := range req.Options {
		fmt.Printf("%v: %v\n", k, v)
	}
	r := &network.JoinResponse{}
	return r, nil
}

func (d *ContrailDriver) Leave(req *network.LeaveRequest) error {
	log.Debugln("=== Leave")
	log.Debugln(req)
	return nil
}

func (d *ContrailDriver) DiscoverNew(req *network.DiscoveryNotification) error {
	log.Debugln("=== DiscoverNew")
	log.Debugln(req)
	return nil
}

func (d *ContrailDriver) DiscoverDelete(req *network.DiscoveryNotification) error {
	log.Debugln("=== DiscoverDelete")
	log.Debugln(req)
	return nil
}

func (d *ContrailDriver) ProgramExternalConnectivity(req *network.ProgramExternalConnectivityRequest) error {
	log.Debugln("=== ProgramExternalConnectivity")
	log.Debugln(req)
	return nil
}

func (d *ContrailDriver) RevokeExternalConnectivity(req *network.RevokeExternalConnectivityRequest) error {
	log.Debugln("=== RevokeExternalConnectivity")
	log.Debugln(req)
	return nil
}
