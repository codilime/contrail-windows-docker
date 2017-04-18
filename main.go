package main

import (
	"flag"
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/codilime/contrail-windows-docker/controller"
	"github.com/codilime/contrail-windows-docker/driver"
)

func main() {
	var adapter = flag.String("adapter", "Ethernet0",
		"net adapter for HNS switch, must be physical")
	var controllerIP = flag.String("controllerIP", "127.0.0.1",
		"IP address of Contrail Controller API")
	var controllerPort = flag.Int("controllerPort", 8082,
		"port of Contrail Controller API")
	var logLevelString = flag.String("logLevel", "Info",
		"log verbosity (possible values: Debug|Info|Warn|Error|Fatal|Panic)")
	flag.Parse()

	logLevel, err := log.ParseLevel(*logLevelString)
	if err != nil {
		log.Error(err)
		return
	}
	log.SetLevel(logLevel)

	keys := &controller.KeystoneEnvs{}
	keys.LoadFromEnvironment()

	c, err := controller.NewController(*controllerIP, *controllerPort, keys)
	if err != nil {
		log.Error(err)
		return
	}

	d := driver.NewDriver(*adapter, c)
	if err = d.StartServing(); err != nil {
		log.Error(err)
	} else {
		defer d.StopServing()
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
	}
}
