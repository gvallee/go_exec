// Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package advexec

import (
	"os/exec"
	"strings"
	"testing"
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
		t.Fatalf("%s does not containt %s", res.Stdout, dummyEnv)
	}
}
