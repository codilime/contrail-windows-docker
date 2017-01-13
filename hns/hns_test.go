package hns

import (
	"strings"
	"testing"

	"github.com/Microsoft/hcsshim"
	"github.com/codilime/contrail-windows-docker/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHNS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HNS wrapper test suite")
}

var _ = BeforeSuite(func() {
	err := common.HardResetHNS()
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := common.HardResetHNS()
	Expect(err).ToNot(HaveOccurred())
})

var _ = Describe("HNS wrapper", func() {

	var originalNumNetworks int

	BeforeEach(func() {
		nets, err := ListHNSNetworks()
		Expect(err).ToNot(HaveOccurred())
		originalNumNetworks = len(nets)
	})

	Context("HNS network exists", func() {

		testNetName := "TestNetwork"
		testHnsNetID := ""

		subnets := []hcsshim.Subnet{
			{
				AddressPrefix:  "1.1.1.0/24",
				GatewayAddress: "1.1.1.1",
			},
		}
		netConfiguration := &hcsshim.HNSNetwork{
			Name:               testNetName,
			Type:               "transparent",
			Subnets:            subnets,
			NetworkAdapterName: "Ethernet0",
		}

		BeforeEach(func() {
			Expect(testHnsNetID).To(Equal(""))
			var err error
			testHnsNetID, err = CreateHNSNetwork(netConfiguration)
			Expect(err).ToNot(HaveOccurred())
			Expect(testHnsNetID).ToNot(Equal(""))
		})

		AfterEach(func() {
			Expect(testHnsNetID).ToNot(Equal(""))
			err := DeleteHNSNetwork(testHnsNetID)
			Expect(err).ToNot(HaveOccurred())
			_, err = GetHNSNetwork(testHnsNetID)
			Expect(err).To(HaveOccurred())
			testHnsNetID = ""
			nets, err := ListHNSNetworks()
			Expect(err).ToNot(HaveOccurred())
			Expect(nets).ToNot(BeNil())
			Expect(len(nets)).To(Equal(originalNumNetworks))
		})

		Specify("listing all HNS networks works", func() {
			nets, err := ListHNSNetworks()
			Expect(err).ToNot(HaveOccurred())
			Expect(nets).ToNot(BeNil())
			Expect(len(nets)).To(Equal(originalNumNetworks + 1))
			found := false
			for _, n := range nets {
				if n.Id == testHnsNetID {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		Specify("getting a single HNS network works", func() {
			net, err := GetHNSNetwork(testHnsNetID)
			Expect(err).ToNot(HaveOccurred())
			Expect(net).ToNot(BeNil())
			Expect(net.Id).To(Equal(testHnsNetID))
		})

		Specify("getting a single HNS network by name works", func() {
			net, err := GetHNSNetworkByName(testNetName)
			Expect(err).ToNot(HaveOccurred())
			Expect(net).ToNot(BeNil())
			Expect(net.Id).To(Equal(testHnsNetID))
		})

		Specify("HNS endpoint operations work", func() {
			hnsEndpointConfig := &hcsshim.HNSEndpoint{
				VirtualNetwork: testHnsNetID,
			}

			endpointID, err := CreateHNSEndpoint(hnsEndpointConfig)
			Expect(err).ToNot(HaveOccurred())
			Expect(endpointID).ToNot(Equal(""))

			endpoint, err := GetHNSEndpoint(endpointID)
			Expect(err).ToNot(HaveOccurred())
			Expect(endpoint).ToNot(BeNil())

			endpoints, err := ListHNSEndpoints()
			Expect(err).ToNot(HaveOccurred())
			Expect(endpoints).To(HaveLen(1))

			err = DeleteHNSEndpoint(endpointID)
			Expect(err).ToNot(HaveOccurred())

			endpoint, err = GetHNSEndpoint(endpointID)
			Expect(err).To(HaveOccurred())
			Expect(endpoint).To(BeNil())
		})

		Specify("Listing HNS endpoints works", func() {
			hnsEndpointConfig := &hcsshim.HNSEndpoint{
				VirtualNetwork: testHnsNetID,
			}

			endpointsList, err := ListHNSEndpoints()
			Expect(err).ToNot(HaveOccurred())
			numEndpointsOriginal := len(endpointsList)

			var endpoints [2]string
			for i := 0; i < 2; i++ {
				endpoints[i], err = CreateHNSEndpoint(hnsEndpointConfig)
				Expect(err).ToNot(HaveOccurred())
				Expect(endpoints[i]).ToNot(Equal(""))
			}

			endpointsList, err = ListHNSEndpoints()
			Expect(err).ToNot(HaveOccurred())
			Expect(endpointsList).To(HaveLen(numEndpointsOriginal + 2))

			for _, ep := range endpoints {
				err = DeleteHNSEndpoint(ep)
				Expect(err).ToNot(HaveOccurred())
			}

			endpointsList, err = ListHNSEndpoints()
			Expect(err).ToNot(HaveOccurred())
			Expect(endpointsList).To(HaveLen(numEndpointsOriginal))
		})
	})

	Context("HNS network doesn't exist", func() {

		BeforeEach(func() {
			nets, err := ListHNSNetworks()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(nets)).To(Equal(originalNumNetworks))
		})

		AfterEach(func() {
			nets, err := ListHNSNetworks()
			Expect(err).ToNot(HaveOccurred())
			for _, n := range nets {
				if strings.Contains(n.Name, "nat") {
					continue
				}
				err = DeleteHNSNetwork(n.Id)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		Specify("getting single HNS network returns error", func() {
			net, err := GetHNSNetwork("1234abcd")
			Expect(err).To(HaveOccurred())
			Expect(net).To(BeNil())
		})

		Specify("getting single HNS network by name returns error", func() {
			net, err := GetHNSNetworkByName("asdf")
			Expect(err).To(HaveOccurred())
			Expect(net).To(BeNil())
		})
	})
})
