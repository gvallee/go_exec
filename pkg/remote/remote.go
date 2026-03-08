//
// Copyright (c) 2023-2026, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE file for license information
//

package remote

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gvallee/go_exec/pkg/advexec"
)

func ExecCmd(host, binPath string, args []string, env []string) advexec.Result {
	if host == "" {
		var newErr advexec.Result
		newErr.Err = fmt.Errorf("host cannot be empty")
		return newErr
	}

	if binPath == "" {
		var newErr advexec.Result
		newErr.Err = fmt.Errorf("binPath cannot be empty")
		return newErr
	}

	sshBinPath, err := exec.LookPath("ssh")
	if err != nil {
		var newErr advexec.Result
		newErr.Err = fmt.Errorf("unable to find ssh: %w", err)
		return newErr
	}

	var cmd advexec.Advcmd
	cmd.BinPath = sshBinPath
	cmd.CmdArgs = []string{host, binPath}
	cmd.CmdArgs = append(cmd.CmdArgs, args...)
	cmd.Env = env

	res := cmd.Run()
	if res.Err != nil {
		res.Err = fmt.Errorf("unable to run %s %s on %s: %w - stderr: %s - stdout: %s", binPath, strings.Join(args, " "), host, res.Err, res.Stderr, res.Stdout)
		return res
	}

	return res
}
