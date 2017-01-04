package driver

import (
	"context"
	"testing"
	"time"

	"github.com/Juniper/contrail-go-api/types"
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

var _ = BeforeSuite(func() {
	err := common.HardResetHNS()
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := common.HardResetHNS()
	Expect(err).ToNot(HaveOccurred())
})

var _ = Describe("Contrail Network Driver", func() {

	Context("upon starting", func() {

		PIt("listens on a pipe", func() {
			d, _, _ := startDriver("")
			defer stopDriver(d)

			_, err := sockets.DialPipe("//./pipe/"+common.DriverName, time.Second*3)
			Expect(err).ToNot(HaveOccurred())
		})

		PIt("tries to connect to existing HNS switch", func() {
			d, _, _ := startDriver("")
			defer stopDriver(d)
		})

		It("if HNS switch doesn't exist, creates a new one", func() {
			net, err := hns.GetHNSNetworkByName(common.NetworkHNSname)
			Expect(net).To(BeNil())
			Expect(err).To(HaveOccurred())

			d, _, _ := startDriver("")
			defer stopDriver(d)

			net, err = hns.GetHNSNetworkByName(common.NetworkHNSname)
			Expect(net).ToNot(BeNil())
			Expect(net.Id).To(Equal(d.HnsID))
			Expect(err).ToNot(HaveOccurred())
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

		var contrailController *controller.Controller
		var contrailDriver *ContrailDriver
		var project *types.Project

		const (
			tenantName  = "agatka"
			networkName = "test_net"
		)

		BeforeEach(func() {
			contrailDriver, contrailController, project = startDriver(tenantName)
		})

		AfterEach(func() {
			stopDriver(contrailDriver)
		})

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

			Context("docker network doesn't exist", func() {
				It("responds with err", func() {
					req := network.CreateEndpointRequest{
						EndpointID: containerID,
					}
					_, err := contrailDriver.CreateEndpoint(&req)
					Expect(err).To(HaveOccurred())
				})
			})

			Context("docker network is misconfigured", func() {
				It("responds with err", func() {
					docker, err := dockerClient.NewEnvClient()
					Expect(err).ToNot(HaveOccurred())
					params := &dockerTypes.NetworkCreate{
						Driver: "Contrail",
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
					_, err = contrailDriver.CreateEndpoint(&req)
					Expect(err).To(HaveOccurred())
				})
			})

			It("allocates Contrail resources", func() {
				// req := &network.CreateEndpointRequest{
				// 	EndpointID: containerID,
				// 	NetworkID:  networkName,
				// }
				// err := contrailDriver.CreateEndpoint(req)
				// Expect(err).ToNot(HaveOccurred())
			})

			PIt("configures container network via HNS", func() {})
			PIt("configures vRouter agent", func() {})
			It("responds with proper CreateEndpointResponse, nil", func() {
				req := network.CreateEndpointRequest{}
				resp, err := contrailDriver.CreateEndpoint(&req)
				Expect(resp).To(Equal(&network.CreateEndpointResponse{}))
				Expect(err).ToNot(HaveOccurred())
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
	Expect(d.HnsID).ToNot(Equal(""))

	return d, mockedController, project
}

func stopDriver(d *ContrailDriver) {
	err := d.Teardown()
	Expect(err).ToNot(HaveOccurred())
	net, err := hns.GetHNSNetworkByName(common.NetworkHNSname)
	Expect(net).To(BeNil())
	Expect(err).To(HaveOccurred())
}
