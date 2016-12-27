package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codilime/contrail-windows-docker/driver"
)

func main() {
	var d driver.ContrailDriver
	var err error

	if d, err = driver.NewDriver(); err != nil {
		logrus.Error(err)
		return
	}

	if err = d.Serve(); err != nil {
		logrus.Error(err)
		return
	}

	if err = d.Teardown(); err != nil {
		logrus.Error(err)
	}
}
