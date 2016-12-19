package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codilime/contrail-windows-docker/driver"
	"github.com/docker/go-plugins-helpers/network"
	"github.com/docker/go-plugins-helpers/sdk"
)

func main() {
	d, err := driver.NewDriver("127.0.0.1", 8082)
	if err != nil {
		logrus.Error(err)
	}
	h := network.NewHandler(d)

	config := sdk.WindowsPipeConfig{
		// This will set permissions for Everyone user allowing him to open, write, read the pipe
		SecurityDescriptor: "S:(ML;;NW;;;LW)D:(A;;0x12019f;;;WD)",
		InBufferSize:       4096,
		OutBufferSize:      4096,
	}

	h.ServeWindows("//./pipe/"+driver.DriverName, driver.DriverName, &config)
	err = d.Teardown()
	if err != nil {
		logrus.Error(err)
	}
}
