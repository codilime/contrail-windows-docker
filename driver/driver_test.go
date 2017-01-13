package driver

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Juniper/contrail-go-api/types"
	"github.com/Microsoft/hcsshim"
	"github.com/codilime/contrail-windows-docker/common"
	"github.com/codilime/contrail-windows-docker/controller"
	"github.com/codilime/contrail-windows-docker/hns"
	dockerTypes "github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-plugins-helpers/network"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDriver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Contrail Network Driver test suite")
}

var contrailController *controller.Controller
var contrailDriver *ContrailDriver
var project *types.Project

var _ = BeforeSuite(func() {
	err := common.HardResetHNS()
	Expect(err).ToNot(HaveOccurred())
	contrailDriver, contrailController, project = startDriver("agatka")
})

var _ = AfterSuite(func() {
	stopDriver(contrailDriver)
	err := common.HardResetHNS()
	Expect(err).ToNot(HaveOccurred())
})

var _ = Describe("Contrail Network Driver", func() {

	Context("upon starting", func() {

		PIt("listens on a pipe", func() {
			_, err := sockets.DialPipe("//./pipe/"+common.DriverName, time.Second*3)
			Expect(err).ToNot(HaveOccurred())
		})

		PIt("tries to connect to existing HNS switch", func() {
		})

		It("if HNS switch doesn't exist, creates a new one", func() {
			// net, err := hns.GetHNSNetworkByName(common.NetworkHNSname)
			// Expect(net).To(BeNil())
			// Expect(err).To(HaveOccurred())

			// d, _, _ := startDriver("")
			// defer stopDriver(d)

			// net, err = hns.GetHNSNetworkByName(common.NetworkHNSname)
			// Expect(net).ToNot(BeNil())
			// Expect(net.Id).To(Equal(d.HnsNetworkID))
			// Expect(err).ToNot(HaveOccurred())
		})

	})

	Describe("allocating resources in Contrail Controller", func() {

		PContext("given correct tenant and subnet id", func() {
			It("works", func() {})
		})

		PContext("given incorrect tenant and subnet id", func() {
			It("returns proper error message", func() {})
		})

	})

	Context("upon shutting down", func() {
		PIt("HNS switch isn't removed", func() {})
	})

	Describe("allocating resources in Contrail Controller", func() {

		PContext("given correct tenant and subnet id", func() {
			It("works", func() {})
		})

		PContext("given incorrect tenant and subnet id", func() {
			It("returns proper error message", func() {})
		})

	})

	Context("on request from docker daemon", func() {

		const (
			tenantName  = "agatka"
			networkName = "test_net"
		)

		// BeforeEach(func() {
		// 	contrailDriver, contrailController, project = startDriver(tenantName)
		// })

		// AfterEach(func() {
		// 	stopDriver(contrailDriver)
		// })

		Context("on GetCapabilities", func() {
			It("returns local scope CapabilitiesResponse, nil", func() {
				resp, err := contrailDriver.GetCapabilities()
				Expect(resp).To(Equal(&network.CapabilitiesResponse{Scope: "local"}))
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("on CreateNetwork", func() {

			Context("subnet doesn't exist in Contrail", func() {
				It("responds with error", func() {
					req := &network.CreateNetworkRequest{
						NetworkID: "MyAwesomeNet",
						Options:   make(map[string]interface{}),
					}
					req.Options["tenant"] = tenantName
					req.Options["network"] = "nonexistingNetwork"
					err := contrailDriver.CreateNetwork(req)
					Expect(err).To(HaveOccurred())
				})
			})

			Context("tenant doesn't exist in Contrail", func() {
				It("responds with error", func() {
					req := &network.CreateNetworkRequest{
						NetworkID: "MyAwesomeNet",
						Options:   make(map[string]interface{}),
					}
					req.Options["tenant"] = "nonexistingTenant"
					req.Options["network"] = networkName
					err := contrailDriver.CreateNetwork(req)
					Expect(err).To(HaveOccurred())
				})
			})

			Context("tenant is not specified in request Options", func() {
				It("responds with error", func() {
					req := &network.CreateNetworkRequest{
						NetworkID: "MyAwesomeNet",
						Options:   make(map[string]interface{}),
					}
					req.Options["network"] = networkName
					err := contrailDriver.CreateNetwork(req)
					Expect(err).To(HaveOccurred())
				})
			})

			Context("network is not specified in request Options", func() {
				It("responds with error", func() {
					req := &network.CreateNetworkRequest{
						NetworkID: "MyAwesomeNet",
						Options:   make(map[string]interface{}),
					}
					req.Options["tenant"] = tenantName
					err := contrailDriver.CreateNetwork(req)
					Expect(err).To(HaveOccurred())
				})
			})

			Context("tenant and subnet exist in Contrail", func() {
				It("responds with nil", func() {
					controller.CreateMockedNetworkWithSubnet(contrailController.ApiClient,
						networkName, "10.10.0.0/24", project)

					req := &network.CreateNetworkRequest{
						NetworkID: "MyAwesomeNet",
						Options:   make(map[string]interface{}),
					}
					req.Options["network"] = networkName
					req.Options["tenant"] = tenantName
					err := contrailDriver.CreateNetwork(req)
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

		PContext("on AllocateNetwork", func() {
			It("responds with empty AllocateNetworkResponse, nil", func() {
				req := network.AllocateNetworkRequest{}
				resp, err := contrailDriver.AllocateNetwork(&req)
				Expect(resp).To(Equal(&network.AllocateNetworkResponse{}))
				Expect(err).ToNot(HaveOccurred())
			})
		})

		PContext("on DeleteNetwork", func() {
			It("responds with nil", func() {
				req := network.DeleteNetworkRequest{}
				err := contrailDriver.DeleteNetwork(&req)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		PContext("on FreeNetwork", func() {
			It("responds with nil", func() {
				req := network.FreeNetworkRequest{}
				err := contrailDriver.FreeNetwork(&req)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("on CreateEndpoint", func() {
			containerID := "2b7ddd5abfad4015b8984bf348a4d51e46a2dd06981f7f5040f9da034d8b631b"

			Context("docker network is set up correctly", func() {

				dockerNetID := ""
				var endpointsBefore []hcsshim.HNSEndpoint

				BeforeEach(func() {
					docker, err := dockerClient.NewEnvClient()
					Expect(err).ToNot(HaveOccurred())
					params := &dockerTypes.NetworkCreate{
						Driver: common.DriverName,
						Options: map[string]string{
							"tenant":  tenantName,
							"network": networkName,
						},
					}
					resp, err := docker.NetworkCreate(context.Background(), networkName, *params)
					Expect(err).ToNot(HaveOccurred())
					dockerNetID = resp.ID

					endpointsBefore, err = hns.ListHNSEndpoints()
					Expect(err).ToNot(HaveOccurred())
				})

				AfterEach(func() {
					endpoints, err := hns.ListHNSEndpoints()
					Expect(err).ToNot(HaveOccurred())
					for _, e := range endpoints {
						for _, originalEndpoint := range endpointsBefore {
							if e.Id == originalEndpoint.Id {
								// don't clean up endpoints that weren't created during tests.
								break
							}
						}
						err = hns.DeleteHNSEndpoint(e.Id)
						Expect(err).ToNot(HaveOccurred())
					}

					docker, err := dockerClient.NewEnvClient()
					Expect(err).ToNot(HaveOccurred())
					err = docker.NetworkRemove(context.Background(), dockerNetID)
					Expect(err).ToNot(HaveOccurred())
				})

				It("allocates Contrail resources", func() {
					req := &network.CreateEndpointRequest{
						EndpointID: containerID,
						NetworkID:  networkName,
					}
					_, err := contrailDriver.CreateEndpoint(req)
					Expect(err).ToNot(HaveOccurred())

					net, err := types.VirtualNetworkByName(contrailController.ApiClient,
						fmt.Sprintf("%s:%s:%s", common.DomainName, tenantName, networkName))
					Expect(err).ToNot(HaveOccurred())
					Expect(net).ToNot(BeNil())

					inst, err := types.VirtualMachineByName(contrailController.ApiClient,
						containerID)
					Expect(err).ToNot(HaveOccurred())
					Expect(inst).ToNot(BeNil())

					vif, err := types.VirtualMachineInterfaceByUuid(contrailController.ApiClient,
						inst.GetUuid())
					Expect(err).ToNot(HaveOccurred())
					Expect(vif).ToNot(BeNil())

					ip, err := types.InstanceIpByUuid(contrailController.ApiClient, inst.GetUuid())
					Expect(err).ToNot(HaveOccurred())
					Expect(ip).ToNot(BeNil())

					ipams, err := net.GetNetworkIpamRefs()
					Expect(err).ToNot(HaveOccurred())
					subnets := ipams[0].Attr.(types.VnSubnetsType).IpamSubnets
					gw := subnets[0].DefaultGateway
					Expect(gw).ToNot(Equal(""))

					macs := vif.GetVirtualMachineInterfaceMacAddresses()
					Expect(macs.MacAddress).To(HaveLen(1))
				})

				It("configures container network via HNS", func() {
					req := network.CreateEndpointRequest{
						EndpointID: containerID,
						NetworkID:  networkName,
					}
					resp, err := contrailDriver.CreateEndpoint(&req)
					Expect(err).To(HaveOccurred())

					endpoints, err := hns.ListHNSEndpoints()
					Expect(err).ToNot(HaveOccurred())
					ep := endpoints[0]
					Expect(ep.IPAddress).To(Equal(resp.Interface.Address))
					Expect(ep.MacAddress).To(Equal(resp.Interface.MacAddress))
				})

				PIt("configures vRouter agent", func() {})
			})

			Context("docker network doesn't exist", func() {
				It("responds with err", func() {
					req := &network.CreateEndpointRequest{
						EndpointID: containerID,
					}
					_, err := contrailDriver.CreateEndpoint(req)
					Expect(err).To(HaveOccurred())
				})
			})

			PContext("docker network is misconfigured", func() {
				It("responds with err", func() {
					docker, err := dockerClient.NewEnvClient()
					Expect(err).ToNot(HaveOccurred())
					params := &dockerTypes.NetworkCreate{
						Options: map[string]string{
							"tenant":  "agatka",
							"network": "lolel",
						},
					}
					_, err = docker.NetworkCreate(context.Background(), networkName, *params)
					Expect(err).ToNot(HaveOccurred())
					req := network.CreateEndpointRequest{
						EndpointID: containerID,
					}
					resp, err := contrailDriver.CreateEndpoint(&req)
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})
			})
		})

		Context("on DeleteEndpoint", func() {
			PIt("deallocates Contrail resources", func() {})
			PIt("configures vRouter Agent", func() {})
			It("responds with nil", func() {
				req := network.DeleteEndpointRequest{}
				err := contrailDriver.DeleteEndpoint(&req)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		PContext("on EndpointInfo", func() {
			It("responds with proper InfoResponse", func() {})
		})

		PContext("on Join", func() {
			Context("queried endpoint exists", func() {
				It("responds with proper JoinResponse", func() {}) // nil maybe?
			})

			Context("queried endpoint doesn't exist", func() {
				It("responds with err", func() {})
			})
		})

		PContext("on Leave", func() {

			Context("queried endpoint exists", func() {
				It("responds with proper JoinResponse, nil", func() {})
			})

			Context("queried endpoint doesn't exist", func() {
				It("responds with err", func() {})
			})
		})

		PContext("on DiscoverNew", func() {
			It("responds with nil", func() {
				req := network.DiscoveryNotification{}
				err := contrailDriver.DiscoverNew(&req)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		PContext("on DiscoverDelete", func() {
			It("responds with nil", func() {
				req := network.DiscoveryNotification{}
				err := contrailDriver.DiscoverDelete(&req)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		PContext("on ProgramExternalConnectivity", func() {
			It("responds with nil", func() {
				req := network.ProgramExternalConnectivityRequest{}
				err := contrailDriver.ProgramExternalConnectivity(&req)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		PContext("on RevokeExternalConnectivity", func() {
			It("responds with nil", func() {
				req := network.RevokeExternalConnectivityRequest{}
				err := contrailDriver.RevokeExternalConnectivity(&req)
				Expect(err).ToNot(HaveOccurred())
			})
		})

	})
})

func startDriver(tenant string) (*ContrailDriver, *controller.Controller, *types.Project) {
	mockedController, project := controller.NewMockedClientAndProject(tenant)

	d, err := NewDriver("172.100.0.0/16", "172.100.0.1", "Ethernet0", mockedController)
	Expect(err).ToNot(HaveOccurred())
	Expect(d.HnsNetworkID).ToNot(Equal(""))

	return d, mockedController, project
}

func stopDriver(d *ContrailDriver) {
	err := d.Teardown()
	Expect(err).ToNot(HaveOccurred())
	net, err := hns.GetHNSNetworkByName(common.NetworkHNSname)
	Expect(net).To(BeNil())
	Expect(err).To(HaveOccurred())
}
