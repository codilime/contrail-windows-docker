package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func CallPowershell(args ...string) (string, string, error) {
	c := []string{"-NonInteractive"}
	for _, arg := range args {
		c = append(c, arg)
	}
	cmd := exec.Command("powershell", c...)

	log.Debugf("Running Powershell command: %s. ", strings.Join(c[1:], " "))

	stdoutPipe, stderrPipe, err := setupOutputCollection(cmd)
	if err != nil {
		return "", "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", "", err
	}

	stdout, stderr, err := collectOutput(stdoutPipe, stderrPipe)
	if err != nil {
		return "", "", err
	}

	// we drop the last two characters, because they are just newline
	if len(stdout) > 2 {
		stdout = stdout[:len(stdout)-2]
	}

	err = cmd.Wait()
	if err != nil {
		return "", "", err
	}

	printDebugInfo(c, stdout, stderr)

	return stdout, stderr, nil
}

func setupOutputCollection(cmd *exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	return stdoutPipe, stderrPipe, err
}

func collectOutput(stdoutPipe, stderrPipe io.ReadCloser) (string, string, error) {
	outBuf, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		return "", "", err
	}
	errBuf, err := ioutil.ReadAll(stderrPipe)
	if err != nil {
		return "", "", err
	}

	// ReadAll returns in []byte, so convert to string
	return fmt.Sprintf("%s", outBuf), fmt.Sprintf("%s", errBuf), nil
}

func printDebugInfo(cmds []string, stdout, stderr string) {
	logMsg := ""
	if stdout != "" {
		logMsg += fmt.Sprintf("stdout: %s;", stdout)
	}
	if stderr != "" {
		logMsg += fmt.Sprintf("stderr: %s;", stderr)
	}
	if logMsg != "" {
		log.Debug(logMsg)
	}
}
