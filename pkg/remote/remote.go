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

	goerrs "github.com/gvallee/go_errs/pkg/goerrs"
	"github.com/gvallee/go_exec/v2/pkg/advexec"
)

const (
	errCodeInvalidInput = goerrs.Code("invalid_input")
	errCodeUnavailable  = goerrs.Code("unavailable")
	errCodeInternal     = goerrs.Code("internal")
)

func ExecCmd(host, binPath string, args []string, env []string) advexec.Result {
	host = strings.TrimSpace(host)
	if host == "" {
		var newErr advexec.Result
		newErr.Err = goerrs.Wrap("remote.ExecCmd", errCodeInvalidInput, fmt.Errorf("host cannot be empty"))
		return newErr
	}

	if strings.ContainsAny(host, " \t\n\r") {
		var newErr advexec.Result
		newErr.Err = goerrs.Wrap("remote.ExecCmd", errCodeInvalidInput, fmt.Errorf("host cannot contain whitespace"))
		return newErr
	}

	binPath = strings.TrimSpace(binPath)
	if binPath == "" {
		var newErr advexec.Result
		newErr.Err = goerrs.Wrap("remote.ExecCmd", errCodeInvalidInput, fmt.Errorf("binPath cannot be empty"))
		return newErr
	}

	sshBinPath, err := exec.LookPath("ssh")
	if err != nil {
		var newErr advexec.Result
		newErr.Err = goerrs.Wrap("remote.ExecCmd", errCodeUnavailable, fmt.Errorf("unable to find ssh: %w", err))
		return newErr
	}

	var cmd advexec.Advcmd
	cmd.BinPath = sshBinPath
	cmd.CmdArgs = []string{host, binPath}
	cmd.CmdArgs = append(cmd.CmdArgs, args...)
	cmd.Env = env

	res := cmd.Run()
	if res.Err != nil {
		res.Err = goerrs.Wrap("remote.ExecCmd", errCodeInternal, fmt.Errorf("unable to run %s %s on %s: %w - stderr: %s - stdout: %s", binPath, strings.Join(args, " "), host, res.Err, res.Stderr, res.Stdout))
		return res
	}

	return res
}
