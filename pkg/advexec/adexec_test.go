// Copyright (c) 2021-2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package advexec

import (
	"os"
	"os/exec"
	"path/filepath"
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

func TestEmptyBinPath(t *testing.T) {
	var c Advcmd
	res := c.Run()
	if res.Err == nil {
		t.Fatalf("expected error for empty bin path")
	}
}

func TestExecCreateManifest(t *testing.T) {
	echoBin, err := exec.LookPath("echo")
	if err != nil {
		t.Skip("'echo' command not available, skipping...")
	}

	execDir := t.TempDir()
	manifestDir := t.TempDir()

	var c Advcmd
	c.BinPath = echoBin
	c.CmdArgs = []string{"hello"}
	c.ExecDir = execDir
	c.ManifestDir = manifestDir
	c.ManifestName = "myexec"
	c.ManifestData = []string{"Meta: value"}

	res := c.Run()
	if res.Err != nil {
		t.Fatalf("execution failed: %s", res.Err)
	}

	manifestPath := filepath.Join(manifestDir, "myexec.MANIFEST")
	if _, statErr := os.Stat(manifestPath); statErr != nil {
		t.Fatalf("manifest was not created: %s", statErr)
	}
}
