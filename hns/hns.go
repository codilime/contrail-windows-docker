package hns

import (
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/Microsoft/hcsshim"
	log "github.com/Sirupsen/logrus"
	"github.com/codilime/contrail-windows-docker/common"
)

func CreateHNSNetwork(configuration *hcsshim.HNSNetwork) (string, error) {
	log.Infoln("Creating HNS network")
	configBytes, err := json.Marshal(configuration)
	if err != nil {
		log.Errorln(err)
		return "", err
	}
	log.Debugln("Config:", string(configBytes))
	response, err := hcsshim.HNSNetworkRequest("POST", "", string(configBytes))
	if err != nil {
		log.Errorln(err)
		return "", err
	}

	// When the first HNS network is created, a vswitch is also created and attached to
	// specified network adapter. This adapter will temporarily lose internet connectivity
	// while it reacquires IPv4. We need to wait for it.
	// https://github.com/Microsoft/hcsshim/issues/108
	if err := waitForInterface(); err != nil {
		log.Errorln(err)
		return "", err
	}

	return response.Id, nil
}

func waitForInterface() error {
	pollingStart := time.Now()
	for {
		queryStart := time.Now()
		iface, err := net.InterfaceByName(common.HNSTransparentInterfaceName)
		if err != nil {
			return err
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		// We print query time because it turns out that above operations actually take quite a
		// while (from 5-400ms), and the time depends (I think) if underlying interface configs
		// are being changed. For example, query usually takes ~10ms, but if it's about to change,
		// it can take up to 400ms. In other words, there must be some kind of mutex there.
		// This information could be useful for debugging.
		log.Debugf("Current %s addresses: %s. Query took %s", common.HNSTransparentInterfaceName,
			addrs, time.Since(queryStart))

		// We're essentialy waiting for adapter to reacquire IPv4 (that's how they do it
		// in Microsoft: https://github.com/Microsoft/hcsshim/issues/108)
		for _, addr := range addrs {
			ip, err, _ := net.ParseCIDR(addr.String())
			if err != nil {
				if ip.To4() != nil {
					log.Debugf("Waited %s for IP reacquisition", time.Since(pollingStart))
					return nil
				}
			}
		}

		if time.Since(pollingStart) > time.Millisecond*common.AdapterReconnectTimeout {
			return errors.New("Waited for net adapter to reconnect for too long.")
		}
		time.Sleep(time.Millisecond * common.AdapterPollingRate)
	}
}

func DeleteHNSNetwork(hnsID string) error {
	log.Infoln("Deleting HNS network", hnsID)
	_, err := hcsshim.HNSNetworkRequest("DELETE", hnsID, "")
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func ListHNSNetworks() ([]hcsshim.HNSNetwork, error) {
	log.Infoln("Listing HNS networks")
	nets, err := hcsshim.HNSListNetworkRequest("GET", "", "")
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return nets, nil
}

func GetHNSNetwork(hnsID string) (*hcsshim.HNSNetwork, error) {
	log.Infoln("Getting HNS network", hnsID)
	net, err := hcsshim.HNSNetworkRequest("GET", hnsID, "")
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return net, nil
}

func GetHNSNetworkByName(name string) (*hcsshim.HNSNetwork, error) {
	log.Infoln("Getting HNS network by name:", name)
	nets, err := hcsshim.HNSListNetworkRequest("GET", "", "")
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	for _, n := range nets {
		if n.Name == name {
			return &n, nil
		}
	}
	return nil, nil
}

func CreateHNSEndpoint(configuration *hcsshim.HNSEndpoint) (string, error) {
	log.Infoln("Creating HNS endpoint")
	configBytes, err := json.Marshal(configuration)
	if err != nil {
		log.Errorln(err)
		return "", err
	}
	log.Debugln("Config: ", string(configBytes))
	response, err := hcsshim.HNSEndpointRequest("POST", "", string(configBytes))
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

func DeleteHNSEndpoint(endpointID string) error {
	log.Infoln("Deleting HNS endpoint", endpointID)
	_, err := hcsshim.HNSEndpointRequest("DELETE", endpointID, "")
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func GetHNSEndpoint(endpointID string) (*hcsshim.HNSEndpoint, error) {
	log.Infoln("Getting HNS endpoint", endpointID)
	endpoint, err := hcsshim.HNSEndpointRequest("GET", endpointID, "")
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	return endpoint, nil
}

func GetHNSEndpointByName(name string) (*hcsshim.HNSEndpoint, error) {
	log.Infoln("Getting HNS endpoint by name:", name)
	eps, err := hcsshim.HNSListEndpointRequest("GET", "", "")
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	for _, ep := range eps {
		if ep.Name == name {
			return &ep, nil
		}
	}
	return nil, nil
}

func ListHNSEndpoints() ([]hcsshim.HNSEndpoint, error) {
	endpoints, err := hcsshim.HNSListEndpointRequest("GET", "", "")
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}

func ListHNSEndpointsOfNetwork(netID string) ([]hcsshim.HNSEndpoint, error) {
	eps, err := ListHNSEndpoints()
	if err != nil {
		return nil, err
	}
	var epsInNetwork []hcsshim.HNSEndpoint
	for _, ep := range eps {
		if ep.VirtualNetwork == netID {
			epsInNetwork = append(epsInNetwork, ep)
		}
	}
	return epsInNetwork, nil
}
