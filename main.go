package main

import (
	"github.com/codilime/contrail-windows-docker/driver"

	"github.com/docker/go-plugins-helpers/network"
	"github.com/docker/go-plugins-helpers/sdk"
)

func main() {
	d := &driver.ContrailDriver{}
	h := network.NewHandler(d)

	config := sdk.WindowsPipeConfig{
		// This will set permissions for Everyone user allowing him to open, write, read the pipe
		SecurityDescriptor: "S:(ML;;NW;;;LW)D:(A;;0x12019f;;;WD)",
		InBufferSize:       4096,
		OutBufferSize:      4096,
	}

	h.ServeWindows("//./pipe/"+driver.Drivername, driver.Drivername, &config)
}
