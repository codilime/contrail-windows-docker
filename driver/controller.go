package driver

import (
	"errors"
	"fmt"

	"github.com/Juniper/contrail-go-api"
	"github.com/Juniper/contrail-go-api/types"
)

const (
	DomainName = "default-domain"
)

type Info struct {
}

type ContrailClient struct {
	ApiClient contrail.ApiClient
}

func NewContrailClient(ip string, port int) *ContrailClient {
	client := &ContrailClient{}
	client.ApiClient = contrail.NewClient(ip, port)
	return client
}

func (c *ContrailClient) GetNetwork(tenantName, networkName string) (*types.VirtualNetwork,
	error) {
	name := fmt.Sprintf("%s:%s:%s", DomainName, tenantName, networkName)
	net, err := types.VirtualNetworkByName(c.ApiClient, name)
	if err != nil {
		return nil, err
	}
	return net, nil
}

func (c *ContrailClient) GetDefaultGatewayIp(net *types.VirtualNetwork) (string, error) {
	ipamReferences, err := net.GetNetworkIpamRefs()
	if err != nil {
		return "", err
	}
	if len(ipamReferences) == 0 {
		return "", errors.New("Ipam references list is empty")
	}
	attribute := ipamReferences[0].Attr
	ipamSubnets := attribute.(types.VnSubnetsType).IpamSubnets
	if len(ipamSubnets) == 0 {
		return "", errors.New("Ipam subnets list is empty")
	}
	gw := ipamSubnets[0].DefaultGateway
	if gw == "" {
		return "", errors.New("Default GW is empty")
	}
	return gw, nil
}

func (c *ContrailClient) GetOrCreateInstance(tenantName, containerId string) (*types.VirtualMachine, error) {
	name := fmt.Sprintf("%s:%s:%s", DomainName, tenantName, containerId)
	instance, err := types.VirtualMachineByName(c.ApiClient, name)
	if err == nil && instance != nil {
		return instance, nil
	}

	instance = new(types.VirtualMachine)
	instance.SetFQName("project", []string{DomainName, tenantName, containerId})
	err = c.ApiClient.Create(instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (c *ContrailClient) GetOrCreateInterface(net *types.VirtualNetwork,
	instance *types.VirtualMachine) (*types.VirtualMachineInterface, error) {
	instanceFQName := instance.GetFQName()
	namespace := instanceFQName[len(instanceFQName)-2]
	name := fmt.Sprintf("%s:%s:%s", DomainName, namespace, instance.GetName())
	iface, err := types.VirtualMachineInterfaceByName(c.ApiClient, name)
	if err == nil && iface != nil {
		return iface, nil
	}

	iface = new(types.VirtualMachineInterface)
	iface.SetFQName("project", []string{DomainName, namespace, instance.GetName()})
	err = iface.AddVirtualMachine(instance)
	if err != nil {
		return nil, err
	}
	err = iface.AddVirtualNetwork(net)
	if err != nil {
		return nil, err
	}
	err = c.ApiClient.Create(iface)
	if err != nil {
		return nil, err
	}
	return iface, nil
}

func (c *ContrailClient) GetInterfaceMac(iface *types.VirtualMachineInterface) (string, error) {
	macs := iface.GetVirtualMachineInterfaceMacAddresses()
	if len(macs.MacAddress) == 0 {
		return "", errors.New("Empty MAC list")
	}
	return macs.MacAddress[0], nil
}

func (c *ContrailClient) GetOrCreateInstanceIp(net *types.VirtualNetwork,
	iface *types.VirtualMachineInterface) (*types.InstanceIp, error) {
	ifaceFQName := iface.GetFQName()
	tenantName := ifaceFQName[len(ifaceFQName)-2]
	name := fmt.Sprintf("%s_%s", tenantName, iface.GetName())
	instIp, err := types.InstanceIpByName(c.ApiClient, name)
	if err == nil && instIp != nil {
		return instIp, nil
	}

	instIp = &types.InstanceIp{}
	err = instIp.AddVirtualNetwork(net)
	if err != nil {
		return nil, err
	}
	err = instIp.AddVirtualMachineInterface(iface)
	if err != nil {
		return nil, err
	}
	err = c.ApiClient.Create(instIp)
	if err != nil {
		return nil, err
	}
	return instIp, nil
}
