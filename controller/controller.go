package controller

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/Juniper/contrail-go-api"
	"github.com/Juniper/contrail-go-api/types"
	log "github.com/Sirupsen/logrus"
	"github.com/codilime/contrail-windows-docker/common"
)

type Info struct {
}

type Controller struct {
	ApiClient contrail.ApiClient
}

type KeystoneEnvs struct {
	Os_auth_url    string
	Os_username    string
	Os_tenant_name string
	Os_password    string
	Os_token       string
}

func (k *KeystoneEnvs) LoadFromEnvironment() {

	k.Os_auth_url = k.GetenvIfNil(k.Os_auth_url, "OS_AUTH_URL")
	k.Os_username = k.GetenvIfNil(k.Os_username, "OS_USERNAME")
	k.Os_tenant_name = k.GetenvIfNil(k.Os_tenant_name, "OS_TENANT_NAME")
	k.Os_password = k.GetenvIfNil(k.Os_password, "OS_PASSWORD")
	k.Os_token = k.GetenvIfNil(k.Os_token, "OS_TOKEN")

	// print a warning for every empty variable
	keysReflection := reflect.ValueOf(*k)
	for i := 0; i < keysReflection.NumField(); i++ {
		if keysReflection.Field(i).String() == "" {
			log.Warn("Keystone variable empty: ", keysReflection.Type().Field(i).Name)
		}
	}
	log.Infoln(k)
}

func (k *KeystoneEnvs) GetenvIfNil(currentVal, envVar string) string {
	if currentVal == "" {
		return os.Getenv(envVar)
	}
	return currentVal
}

func NewController(ip string, port int, keys *KeystoneEnvs) (*Controller, error) {
	client := &Controller{}
	client.ApiClient = contrail.NewClient(ip, port)

	if keys.Os_auth_url == "" {
		// this corner case is not handled by keystone.Authenticate. Causes panic.
		return nil, errors.New("Empty Keystone auth URL")
	}

	keystone := contrail.NewKeystoneClient(keys.Os_auth_url, keys.Os_tenant_name,
		keys.Os_username, keys.Os_password, keys.Os_token)
	err := keystone.Authenticate()
	if err != nil {
		log.Errorln("Keystone error:", err)
		return nil, err
	}
	client.ApiClient.(*contrail.Client).SetAuthenticator(keystone)
	return client, nil
}

func (c *Controller) GetNetwork(tenantName, networkName string) (*types.VirtualNetwork,
	error) {
	name := fmt.Sprintf("%s:%s:%s", common.DomainName, tenantName, networkName)
	net, err := types.VirtualNetworkByName(c.ApiClient, name)
	if err != nil {
		log.Errorf("Failed to get virtual network %s by name: %v", name, err)
		return nil, err
	}
	return net, nil
}

func (c *Controller) GetIpamSubnet(net *types.VirtualNetwork) (*types.IpamSubnetType, error) {
	ipamReferences, err := net.GetNetworkIpamRefs()
	if err != nil {
		log.Errorf("Failed to get ipam references: %v", err)
		return nil, err
	}
	if len(ipamReferences) == 0 {
		err = errors.New("Ipam references list is empty")
		log.Error(err)
		return nil, err
	}
	attribute := ipamReferences[0].Attr
	ipamSubnets := attribute.(types.VnSubnetsType).IpamSubnets
	if len(ipamSubnets) == 0 {
		err = errors.New("Ipam subnets list is empty")
		log.Error(err)
		return nil, err
	}
	return &ipamSubnets[0], nil
}

func (c *Controller) GetDefaultGatewayIp(net *types.VirtualNetwork) (string, error) {
	subnet, err := c.GetIpamSubnet(net)
	if err != nil {
		return "", err
	}
	gw := subnet.DefaultGateway
	if gw == "" {
		err = errors.New("Default GW is empty")
		log.Error(err)
		return "", err
	}
	return gw, nil
}

func (c *Controller) GetOrCreateInstance(vif *types.VirtualMachineInterface, containerId string) (
	*types.VirtualMachine, error) {
	instance, err := types.VirtualMachineByName(c.ApiClient, containerId)
	if err == nil && instance != nil {
		return instance, nil
	}

	instance = new(types.VirtualMachine)
	instance.SetName(containerId)
	err = c.ApiClient.Create(instance)
	if err != nil {
		log.Errorf("Failed to create instance: %v", err)
		return nil, err
	}

	createdInstance, err := types.VirtualMachineByName(c.ApiClient, containerId)
	if err != nil {
		log.Errorf("Failed to retreive instance %s by name: %v", containerId, err)
		return nil, err
	}
	log.Infoln("Created instance: ", createdInstance.GetFQName())

	err = vif.AddVirtualMachine(createdInstance)
	if err != nil {
		log.Errorf("Failed to add instance to vif")
		return nil, err
	}
	err = c.ApiClient.Update(vif)
	if err != nil {
		log.Errorf("Failed to update vif")
		return nil, err
	}

	return createdInstance, nil
}

func (c *Controller) GetOrCreateInterface(net *types.VirtualNetwork, tenantName,
	containerId string) (*types.VirtualMachineInterface, error) {

	fqName := fmt.Sprintf("%s:%s:%s", common.DomainName, tenantName, containerId)
	iface, err := types.VirtualMachineInterfaceByName(c.ApiClient, fqName)
	if err == nil && iface != nil {
		return iface, nil
	}

	iface = new(types.VirtualMachineInterface)
	iface.SetFQName("project", []string{common.DomainName, tenantName, containerId})
	err = iface.AddVirtualNetwork(net)
	if err != nil {
		log.Errorf("Failed to add network to interface: %v", err)
		return nil, err
	}
	err = c.ApiClient.Create(iface)
	if err != nil {
		log.Errorf("Failed to create interface: %v", err)
		return nil, err
	}

	createdIface, err := types.VirtualMachineInterfaceByName(c.ApiClient, fqName)
	if err != nil {
		log.Errorf("Failed to retreive vmi %s by name: %v", fqName, err)
		return nil, err
	}
	log.Infoln("Created instance: ", createdIface.GetFQName())
	return createdIface, nil
}

func (c *Controller) GetInterfaceMac(iface *types.VirtualMachineInterface) (string, error) {
	macs := iface.GetVirtualMachineInterfaceMacAddresses()
	if len(macs.MacAddress) == 0 {
		err := errors.New("Empty MAC list")
		log.Error(err)
		return "", err
	}
	return macs.MacAddress[0], nil
}

func (c *Controller) GetOrCreateInstanceIp(net *types.VirtualNetwork,
	iface *types.VirtualMachineInterface) (*types.InstanceIp, error) {
	instIp, err := types.InstanceIpByName(c.ApiClient, iface.GetName())
	if err == nil && instIp != nil {
		return instIp, nil
	}

	instIp = &types.InstanceIp{}
	instIp.SetName(iface.GetName())
	err = instIp.AddVirtualNetwork(net)
	if err != nil {
		log.Errorf("Failed to add network to instanceIP object: %v", err)
		return nil, err
	}
	err = instIp.AddVirtualMachineInterface(iface)
	if err != nil {
		log.Errorf("Failed to add vmi to instanceIP object: %v", err)
		return nil, err
	}
	err = c.ApiClient.Create(instIp)
	if err != nil {
		log.Errorf("Failed to instanceIP: %v", err)
		return nil, err
	}

	allocatedIP, err := types.InstanceIpByUuid(c.ApiClient, instIp.GetUuid())
	if err != nil {
		log.Errorf("Failed to retreive instanceIP object %s by name: %v", instIp.GetUuid(), err)
		return nil, err
	}
	return allocatedIP, nil
}

func (c *Controller) DeleteElementRecursive(parent contrail.IObject) error {
	log.Debugln("Deleting", parent.GetType(), parent.GetUuid())
	for err := c.ApiClient.Delete(parent); err != nil; err = c.ApiClient.Delete(parent) {
		if strings.Contains(err.Error(), "404 Resource") {
			log.Errorln("Failed to delete Contrail resource", err.Error())
			break
		} else if strings.Contains(err.Error(), "409 Conflict") {
			msg := err.Error()
			// example error message when object has children:
			// `409 Conflict: Delete when children still present:
			// ['http://10.7.0.54:8082/virtual-network/23e300f4-ab1a-4d97-a1d9-9ed69b601e17']`

			// This regex finds all strings like:
			// `virtual-network/23e300f4-ab1a-4d97-a1d9-9ed69b601e17`
			var re *regexp.Regexp
			re, err = regexp.Compile(
				"([a-z-]*\\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})")
			if err != nil {
				return err
			}
			found := re.FindAll([]byte(msg), -1)

			for _, f := range found {
				split := strings.Split(string(f), "/")
				typename := split[0]
				UUID := split[1]
				var child contrail.IObject
				child, err = c.ApiClient.FindByUuid(typename, UUID)
				if err != nil {
					return err
				}
				if child == nil {
					return errors.New("Child object is nil")
				}
				err = c.DeleteElementRecursive(child)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
