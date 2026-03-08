//
// Copyright (c) 2023-2026, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE file for license information
//

package remote

import (
	"os/exec"
	"strings"
	"testing"
)

func TestRemoteCmd(t *testing.T) {
	cmd, err := exec.LookPath("date")
	if err != nil {
		t.Skip("'date' command not available, skipping...")
	}
	host := "localhost"
	res := ExecCmd(host, cmd, nil, nil)
	if res.Err != nil {
		t.Fatalf("unable to run %s on %s", cmd, host)
	}
}

func TestRemoteCmdWithArgs(t *testing.T) {
	cmd, err := exec.LookPath("echo")
	if err != nil {
		t.Skip("'echo' command not available, skipping...")
	}

	res := ExecCmd("localhost", cmd, []string{"hello"}, nil)
	if res.Err != nil {
		t.Fatalf("unable to run remote command with args: %s", res.Err)
	}

	if !strings.Contains(res.Stdout, "hello") {
		t.Fatalf("unexpected stdout: %q", res.Stdout)
	}
}

func TestRemoteCmdValidation(t *testing.T) {
	res := ExecCmd("", "/bin/true", nil, nil)
	if res.Err == nil {
		t.Fatalf("expected error for empty host")
	}

	res = ExecCmd("bad host", "/bin/true", nil, nil)
	if res.Err == nil {
		t.Fatalf("expected error for host with whitespace")
	}

	res = ExecCmd("localhost", "", nil, nil)
	if res.Err == nil {
		t.Fatalf("expected error for empty binPath")
	}
}
