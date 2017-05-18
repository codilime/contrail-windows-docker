package hyperv

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codilime/contrail-windows-docker/common"
)

func EnableExtension(vswitchName common.VSwitchName, netAdapterName common.AdapterName) error {
	log.Infoln("Enabling vRouter Hyper-V Extension")
	if out, err := callOnSwitch(vswitchName, netAdapterName, "Enable-VMSwitchExtension"); err != nil {
		log.Errorf("When enabling Hyper-V Extension: %s", err, out)
		return err
	}
	return nil
}

func DisableExtension(vswitchName common.VSwitchName, netAdapterName common.AdapterName) error {
	log.Infoln("Disabling vRouter Hyper-V Extension")
	if out, err := callOnSwitch(vswitchName, netAdapterName, "Disable-VMSwitchExtension"); err != nil {
		log.Errorf("When disabling Hyper-V Extension: %s", err, out)
		return err
	}
	return nil
}

func IsExtensionEnabled(vswitchName common.VSwitchName, netAdapterName common.AdapterName) (bool,
	error) {
	out, err := inspectExtensionProperty(vswitchName, netAdapterName, "Enabled")
	if err != nil {
		log.Errorf("When inspecting Hyper-V Extension: %s", err, out)
		return false, err
	}
	return out == "True", nil
}

func IsExtensionRunning(vswitchName common.VSwitchName, netAdapterName common.AdapterName) (bool,
	error) {
	out, err := inspectExtensionProperty(vswitchName, netAdapterName, "Running")
	if err != nil {
		log.Errorf("When inspecting Hyper-V Extension: %s", err, out)
		return false, err
	}
	return out == "True", nil
}

func inspectExtensionProperty(vswitchName common.VSwitchName, netAdapterName common.AdapterName,
	property string) (string, error) {
	log.Infoln("Inspecting vRouter Hyper-V Extension for property:", property)
	// we use -Expand, because otherwise, we get an object instead of single string value
	out, err := callOnSwitch(vswitchName, netAdapterName, "Get-VMSwitchExtension", "|", "Select",
		"-Expand", fmt.Sprintf("\"%s\"", property))
	log.Debugln("Inspect result:", out)
	return out, err
}

func callOnSwitch(vswitchName common.VSwitchName, netAdapterName common.AdapterName,
	command string, optionals ...string) (string, error) {
	c := []string{command,
		"-VMSwitchName", string(vswitchName),
		"-Name", fmt.Sprintf("\"%s\"", common.HyperVExtensionName)}
	for _, opt := range optionals {
		c = append(c, opt)
	}
	stdout, _, err := common.CallPowershell(c...)
	return stdout, err
}
