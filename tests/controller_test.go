package driver_test

import (
	"fmt"
	"testing"

	"github.com/codilime/contrail-windows-docker/driver"

	"github.com/Juniper/contrail-go-api"
	"github.com/Juniper/contrail-go-api/config"
	"github.com/Juniper/contrail-go-api/mocks"
	"github.com/Juniper/contrail-go-api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller client test suite")
}

const (
	tenantName   = "agatka"
	networkName  = "test_net"
	subnetPrefix = "10.10.10.0"
	subnetMask   = 24
	defaultGW    = "10.10.10.1"
	ifaceMac     = "contrail_pls_check_macs"
	containerId  = "12345678901"
)

var _ = Describe("Controller", func() {

	var client *driver.ContrailClient
	var project *types.Project

	BeforeEach(func() {
		client = newMockedClient()
		project = new(types.Project)
		project.SetFQName("domain", []string{driver.DomainName, tenantName})
		err := client.ApiClient.Create(project)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("getting Contrail network", func() {
		Context("when network already exists in Contrail", func() {
			var testNetwork *types.VirtualNetwork
			BeforeEach(func() {
				testNetwork = createMockedNetworkWithSubnet(client.ApiClient, project)
			})
			It("returns it", func() {
				net, err := client.GetNetwork(tenantName, networkName)
				Expect(err).ToNot(HaveOccurred())
				Expect(net).To(BeEquivalentTo(testNetwork))
			})
		})
		Context("when network doesn't exist in Contrail", func() {
			It("returns an error", func() {
				net, err := client.GetNetwork(tenantName, networkName)
				Expect(err).To(HaveOccurred())
				Expect(net).To(BeNil())
			})
		})
	})

	Describe("getting Contrail default gateway IP", func() {
		Context("network has subnet with default gateway", func() {
			var testNetwork *types.VirtualNetwork
			BeforeEach(func() {
				testNetwork = createMockedNetworkWithSubnet(client.ApiClient, project)
				addDefaultGatewayToSubnet(client.ApiClient, testNetwork)
			})
			It("returns gateway IP", func() {
				gwAddr, err := client.GetDefaultGatewayIp(testNetwork)
				Expect(err).ToNot(HaveOccurred())
				Expect(gwAddr).ToNot(Equal(""))
			})
		})
		Context("network has subnet without default gateway", func() {
			var testNetwork *types.VirtualNetwork
			BeforeEach(func() {
				testNetwork = createMockedNetworkWithSubnet(client.ApiClient, project)
			})
			It("returns error", func() {
				gwAddr, err := client.GetDefaultGatewayIp(testNetwork)
				Expect(err).To(HaveOccurred())
				Expect(gwAddr).To(Equal(""))
			})
		})
		Context("network doesn't have subnets", func() {
			var testNetwork *types.VirtualNetwork
			BeforeEach(func() {
				testNetwork = createMockedNetwork(client.ApiClient, project)
			})
			It("returns error", func() {
				gwAddr, err := client.GetDefaultGatewayIp(testNetwork)
				Expect(err).To(HaveOccurred())
				Expect(gwAddr).To(Equal(""))
			})
		})
	})

	Describe("getting Contrail instance", func() {
		Context("when instance already exists in Cotnrail", func() {
			var testInstance *types.VirtualMachine
			BeforeEach(func() {
				testInstance = createMockedInstance(client.ApiClient)
			})
			It("returns existing instance", func() {
				instance, err := client.GetOrCreateInstance(tenantName, containerId)
				Expect(err).ToNot(HaveOccurred())
				Expect(instance).ToNot(BeNil())
				Expect(instance).To(BeEquivalentTo(testInstance))
			})
		})
		Context("when instance doesn't exist in Contrail", func() {
			It("creates a new instance", func() {
				instance, err := client.GetOrCreateInstance(tenantName, containerId)
				Expect(err).ToNot(HaveOccurred())
				Expect(instance).ToNot(BeNil())
			})
		})
	})

	Describe("getting Contrail virtual interface", func() {
		var testNetwork *types.VirtualNetwork
		var testInstance *types.VirtualMachine
		BeforeEach(func() {
			testNetwork = createMockedNetworkWithSubnet(client.ApiClient, project)
			testInstance = createMockedInstance(client.ApiClient)
		})
		Context("when vif already exists in Contrail", func() {
			var testInterface *types.VirtualMachineInterface
			BeforeEach(func() {
				testInterface = createMockedInterface(client.ApiClient, testInstance,
					testNetwork)
			})
			It("returns existing vif", func() {
				iface, err := client.GetOrCreateInterface(testNetwork, testInstance)
				Expect(err).ToNot(HaveOccurred())
				Expect(iface).ToNot(BeNil())
				Expect(iface).To(BeEquivalentTo(testInterface))
			})
		})
		Context("when vif doesn't exist in Contrail", func() {
			It("creates a new vif", func() {
				iface, err := client.GetOrCreateInterface(testNetwork, testInstance)
				Expect(err).ToNot(HaveOccurred())
				Expect(iface).ToNot(BeNil())
			})
		})
	})

	Describe("getting virtual interface MAC", func() {
		var testNetwork *types.VirtualNetwork
		var testInstance *types.VirtualMachine
		var testInterface *types.VirtualMachineInterface
		BeforeEach(func() {
			testNetwork = createMockedNetworkWithSubnet(client.ApiClient, project)
			testInstance = createMockedInstance(client.ApiClient)
			testInterface = createMockedInterface(client.ApiClient, testInstance,
				testNetwork)
		})
		Context("when vif has MAC", func() {
			BeforeEach(func() {
				addMacToInterface(client.ApiClient, testInterface)
			})
			It("returns MAC address", func() {
				mac, err := client.GetInterfaceMac(testInterface)
				Expect(err).ToNot(HaveOccurred())
				Expect(mac).To(Equal(ifaceMac))
			})
		})
		Context("when vif doesn't have a MAC", func() {
			It("returns error", func() {
				mac, err := client.GetInterfaceMac(testInterface)
				Expect(err).To(HaveOccurred())
				Expect(mac).To(Equal(""))
			})
		})
	})

	Describe("getting Contrail instance IP", func() {
		var testNetwork *types.VirtualNetwork
		var testInstance *types.VirtualMachine
		var testInterface *types.VirtualMachineInterface
		BeforeEach(func() {
			testNetwork = createMockedNetworkWithSubnet(client.ApiClient, project)
			testInstance = createMockedInstance(client.ApiClient)
			testInterface = createMockedInterface(client.ApiClient, testInstance,
				testNetwork)
		})
		Context("when instance IP already exists in Contrail", func() {
			var testInstanceIP *types.InstanceIp
			BeforeEach(func() {
				testInstanceIP = createMockedInstanceIp(client.ApiClient, testInterface,
					testNetwork)
			})
			It("returns existing instance IP", func() {
				instanceIP, err := client.GetOrCreateInstanceIp(testNetwork, testInterface)
				Expect(err).ToNot(HaveOccurred())
				Expect(instanceIP).ToNot(BeNil())
				Expect(instanceIP).To(BeEquivalentTo(testInstanceIP))
			})
		})
		Context("when instance IP doesn't exist in Contrail", func() {
			It("creates new instance IP", func() {
				instanceIP, err := client.GetOrCreateInstanceIp(testNetwork, testInterface)
				Expect(err).ToNot(HaveOccurred())
				Expect(instanceIP).ToNot(BeNil())
			})
		})
	})
})

func newMockedClient() *driver.ContrailClient {
	c := &driver.ContrailClient{}
	mockedApiClient := new(mocks.ApiClient)
	mockedApiClient.Init()
	c.ApiClient = mockedApiClient
	return c
}

func createMockedNetworkWithSubnet(c contrail.ApiClient, project *types.Project) *types.VirtualNetwork {
	netUuid, err := config.CreateNetwork(c, project.GetUuid(), networkName)
	Expect(err).ToNot(HaveOccurred())
	Expect(netUuid).ToNot(BeNil())
	testNetwork, err := types.VirtualNetworkByUuid(c, netUuid)
	Expect(err).ToNot(HaveOccurred())
	Expect(testNetwork).ToNot(BeNil())
	return testNetwork
}

func addDefaultGatewayToSubnet(c contrail.ApiClient, testNetwork *types.VirtualNetwork) {
	subnet := &types.IpamSubnetType{
		Subnet:         &types.SubnetType{IpPrefix: subnetPrefix, IpPrefixLen: subnetMask},
		DefaultGateway: defaultGW,
	}

	var ipamSubnets types.VnSubnetsType
	ipamSubnets.AddIpamSubnets(subnet)

	ipam, err := c.FindByName("network-ipam",
		"default-domain:default-project:default-network-ipam")
	Expect(err).ToNot(HaveOccurred())
	err = testNetwork.AddNetworkIpam(ipam.(*types.NetworkIpam), ipamSubnets)
	Expect(err).ToNot(HaveOccurred())

	err = c.Update(testNetwork)
	Expect(err).ToNot(HaveOccurred())
}

func createMockedNetwork(c contrail.ApiClient, project *types.Project) *types.VirtualNetwork {
	netUuid, err := config.CreateNetwork(c, project.GetUuid(), networkName)
	Expect(err).ToNot(HaveOccurred())
	Expect(netUuid).ToNot(BeNil())
	testNetwork, err := types.VirtualNetworkByUuid(c, netUuid)
	Expect(err).ToNot(HaveOccurred())
	Expect(testNetwork).ToNot(BeNil())
	return testNetwork
}

func createMockedInstance(c contrail.ApiClient) *types.VirtualMachine {
	testInstance := new(types.VirtualMachine)
	testInstance.SetFQName("project", []string{driver.DomainName, tenantName, containerId})
	err := c.Create(testInstance)
	Expect(err).ToNot(HaveOccurred())
	return testInstance
}

func createMockedInterface(c contrail.ApiClient, instance *types.VirtualMachine,
	net *types.VirtualNetwork) *types.VirtualMachineInterface {
	iface := new(types.VirtualMachineInterface)
	instanceFQName := instance.GetFQName()
	namespace := instanceFQName[len(instanceFQName)-2]
	iface.SetFQName("project", []string{driver.DomainName, namespace, instance.GetName()})
	err := iface.AddVirtualMachine(instance)
	Expect(err).ToNot(HaveOccurred())
	err = iface.AddVirtualNetwork(net)
	Expect(err).ToNot(HaveOccurred())
	err = c.Create(iface)
	Expect(err).ToNot(HaveOccurred())
	return iface
}

func addMacToInterface(c contrail.ApiClient, iface *types.VirtualMachineInterface) {
	macs := new(types.MacAddressesType)
	macs.AddMacAddress(ifaceMac)
	iface.SetVirtualMachineInterfaceMacAddresses(macs)
	err := c.Update(iface)
	Expect(err).ToNot(HaveOccurred())
}

func createMockedInstanceIp(c contrail.ApiClient, iface *types.VirtualMachineInterface,
	net *types.VirtualNetwork) *types.InstanceIp {
	name := fmt.Sprintf("%s_%s", tenantName, iface.GetName())

	instIp := &types.InstanceIp{}
	instIp.SetName(name)
	err := instIp.AddVirtualNetwork(net)
	Expect(err).ToNot(HaveOccurred())
	err = instIp.AddVirtualMachineInterface(iface)
	Expect(err).ToNot(HaveOccurred())
	err = c.Create(instIp)
	Expect(err).ToNot(HaveOccurred())

	allocatedIp, err := c.FindByUuid(instIp.GetType(), instIp.GetUuid())
	Expect(err).ToNot(HaveOccurred())
	Expect(allocatedIp.(*types.InstanceIp).GetInstanceIpAddress()).ToNot(BeNil())

	return allocatedIp.(*types.InstanceIp)
}
