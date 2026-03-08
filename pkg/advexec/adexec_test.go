// Copyright (c) 2021-2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package advexec

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestExecWithEnv(t *testing.T) {
	var c Advcmd
	var err error
	c.BinPath, err = exec.LookPath("env")
	if err != nil {
		t.Skip("'env' command not available, skipping...")
	}
	dummyEnv := "TOTO=titi"
	c.Env = append(c.Env, dummyEnv)
	res := c.Run()
	if res.Err != nil {
		t.Fatalf("execution failed: %s, stdout:%s, stderr:%s", res.Err, res.Stdout, res.Stderr)
	}

	if !strings.Contains(res.Stdout, dummyEnv) {
		t.Fatalf("%s does not contain %s", res.Stdout, dummyEnv)
	}
}

func TestExecTimeout(t *testing.T) {
	var c Advcmd
	var err error

	c.BinPath, err = exec.LookPath("sleep")
	if err != nil {
		t.Skip("'sleep' command not available, skipping...")
	}
	c.CmdArgs = append(c.CmdArgs, "10")
	// Set a timeout of 5 seconds
	c.Timeout = 5 * time.Second
	res := c.Run()
	if res.Err != nil {
		t.Logf("Timeout detected: %s", res.Err)
	} else {
		t.Fatalf("Timeout not detected")
	}
}
