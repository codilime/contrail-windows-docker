package common

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"unicode/utf16"

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

	// trim leading and trailing whitespace
	stdout = strings.TrimSpace(stdout)
	stderr = strings.TrimSpace(stderr)

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
	//return fmt.Sprintf("%s", outBuf), fmt.Sprintf("%s", errBuf), nil
	return decodeUtf8(outBuf), decodeUtf8(errBuf), err
}

func decodeUtf8(rawMsg []byte) string {
	// word := "Héllo!"
	// // rawMsg = []byte(word)
	// log.Infoln("raw:", rawMsg)
	// var runes []rune
	// for len(rawMsg) > 0 {
	// 	decodedRune, decodedRuneSize := utf8.DecodeRune(rawMsg)
	// 	decodedRune = utf16.DecodeRune()
	// 	runes = append(runes, decodedRune)
	// 	log.Infoln(decodedRune)
	// 	rawMsg = rawMsg[decodedRuneSize:]
	// }
	// decoded := string(runes)
	// log.Infoln("runes", runes)
	// log.Infoln("dec", decoded)
	// return decoded

	msg, err := utf16toString(rawMsg)
	if err != nil {
		log.Errorln(err)
	}
	return msg
}

func utf16toString(b []uint8) (string, error) {
	if len(b)&1 != 0 {
		return "", errors.New("len(b) must be even")
	}

	// Check BOM
	var bom int
	if len(b) >= 2 {
		switch n := int(b[0])<<8 | int(b[1]); n {
		case 0xfffe:
			bom = 1
			fallthrough
		case 0xfeff:
			b = b[2:]
		}
	}

	w := make([]uint16, len(b)/2)
	for i := range w {
		w[i] = uint16(b[2*i+bom&1])<<8 | uint16(b[2*i+(bom+1)&1])
	}
	return string(utf16.Decode(w)), nil
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
