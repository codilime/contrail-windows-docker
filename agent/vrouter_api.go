package agent

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codilime/contrail-windows-docker/common"
)

func AddPort(vmUuid, vifUuid, ifName, mac, dockerID string) error {
	stdout, stderr, err := common.Call("python", common.AgentAPIWrapperScriptPath,
		"add", vmUuid, vifUuid, ifName, mac, dockerID)
	log.Debugln("Called Agent API wrapper: ", stdout)
	if err != nil {
		log.Errorf("When calling Agent API wrapper script: %s, %s", stdout, stderr)
		return err
	}
	return nil
}

func DeletePort(vifUuid string) error {
	stdout, stderr, err := common.Call("python", common.AgentAPIWrapperScriptPath,
		"delete", vifUuid)
	log.Debugln("Called Agent API wrapper: ", stdout)
	if err != nil {
		log.Errorf("When calling Agent API wrapper script: %s, %s", stdout, stderr)
		return err
	}
	return nil
}
