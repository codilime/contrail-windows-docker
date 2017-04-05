package hyperv

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codilime/contrail-windows-docker/common"
)

func EnableExtension(netAdapterName string) error {
	log.Infoln("Enabling vRouter Hyper-V Extension")
	if out, err := callOnSwitch("Enable-VMSwitchExtension", netAdapterName); err != nil {
		log.Errorf("When enabling Hyper-V Extension: %s", err, out)
		return err
	}
	return nil
}

func DisableExtension(netAdapterName string) error {
	log.Infoln("Disabling vRouter Hyper-V Extension")
	if out, err := callOnSwitch("Disable-VMSwitchExtension", netAdapterName); err != nil {
		log.Errorf("When disabling Hyper-V Extension: %s", err, out)
		return err
	}
	return nil
}

func IsExtensionEnabled(netAdapterName string) (bool, error) {
	out, err := inspectExtensionProperty(netAdapterName, "Enabled")
	if err != nil {
		log.Errorf("When inspecting Hyper-V Extension: %s", err, out)
		return false, err
	}
	return out == "True", nil
}

func IsExtensionRunning(netAdapterName string) (bool, error) {
	out, err := inspectExtensionProperty(netAdapterName, "Running")
	if err != nil {
		log.Errorf("When inspecting Hyper-V Extension: %s", err, out)
		return false, err
	}
	return out == "True", nil
}

func inspectExtensionProperty(netAdapterName, property string) (string, error) {
	log.Infoln("Inspecting vRouter Hyper-V Extension for property:", property)
	// we use -Expand, because otherwise, we get an object instead of single string value
	out, err := callOnSwitch("Get-VMSwitchExtension", netAdapterName, "|", "Select",
		"-Expand", fmt.Sprintf("\"%s\"", property))
	log.Debugln("Inspect result:", out)
	return out, err
}

func callOnSwitch(command, netAdapterName string, optionals ...string) (string, error) {
	c := []string{command,
		"-VMSwitchName", fmt.Sprintf("\"%s\"", switchName(netAdapterName)),
		"-Name", fmt.Sprintf("\"%s\"", common.HyperVExtensionName)}
	for _, opt := range optionals {
		c = append(c, opt)
	}
	stdout, _, err := common.CallPowershell(c...)
	return stdout, err
}

func switchName(netAdapterName string) string {
	return fmt.Sprintf("Layered %s", netAdapterName)
}
