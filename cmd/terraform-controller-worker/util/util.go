package util

import (
	"bytes"
	"fmt"
	"os/exec"

	"k8s.io/klog/v2"
)

func RunTerraformCommand(command string, stdin *string) error {
	cmd := exec.Command("terraform", command)
	var out bytes.Buffer
	var stdErr bytes.Buffer

	if stdin != nil {
		cmd.Stdin = bytes.NewReader([]byte(*stdin))
	}

	cmd.Stdout = &out
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		klog.Infof("tf %s output:\n%v\n", command, stdErr.String())
		return fmt.Errorf("Error running terraform %s: %v", command, err)
	}

	klog.Infof("tf %s output:\n%v\n", command, out.String())
	return nil
}
