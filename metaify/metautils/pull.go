package metautils

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func PullUpdate() (string, error) {
	cmd := exec.Command("git", "pull")
	so, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	eo, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	// io.Copy(os.Stdout, so)
	// io.Copy(os.Stderr, eo)

	norout, err := ioutil.ReadAll(so)
	if err != nil {
		return "", err
	}

	errout, err := ioutil.ReadAll(eo)
	if err != nil {
		return "", err
	}

	if len(errout) != 0 {
		return "", fmt.Errorf(string(errout))
	}

	return string(norout), nil
}
