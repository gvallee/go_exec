// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package advexec

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gvallee/go_exec/pkg/manifest"
	"github.com/gvallee/go_util/pkg/util"
)

const (
	// CmdTimeout is the maximum time we allow a command to run
	CmdTimeout = 30
)

// Result represents the result of the execution of a command
type Result struct {
	// Err is the Go error associated to the command execution
	Err error
	// Stdout is the messages that were displayed on stdout during the execution of the command
	Stdout string
	// Stderr is the messages that were displayed on stderr during the execution of the command
	Stderr string
}

// Advcmd represents an advanced command to be executed
type Advcmd struct {
	// Cmd represents the command to execute to submit the job
	Cmd *exec.Cmd

	// Timeout is the maximum time a command can run
	Timeout time.Duration

	// BinPath is the path to the binary to execute
	BinPath string

	// CmdArgs is a slice of string representing the command's arguments
	CmdArgs []string

	// ExecDir is the directory where to execute the command
	ExecDir string

	// Env is a slice of string representing the environment to be used with the command
	Env []string

	// Ctx is the context of the command to execute to submit a job
	Ctx context.Context

	// CancelFn is the function to cancel the command to submit a job
	CancelFn context.CancelFunc

	// ManifestName is the name of the manifest, it will default to "exec.MANIFEST" if not defined
	ManifestName string

	// ManifestDir is the directory where to create the manifest related to the command execution
	ManifestDir string

	// ManifestData is extra content to add to the manifest
	ManifestData []string

	// ManifestFileHash is a list of absolute path to files for which we want a hash in the manifest
	ManifestFileHash []string
}

// Run executes a syexec command and creates the appropriate manifest (when possible)
func (c *Advcmd) Run() Result {
	var res Result

	cmdTimeout := c.Timeout
	if cmdTimeout == 0 {
		cmdTimeout = CmdTimeout * time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()

	var stderr, stdout bytes.Buffer
	if c.Cmd == nil {
		c.Cmd = exec.CommandContext(ctx, c.BinPath, c.CmdArgs...)
		c.Cmd.Stdout = &stdout
		c.Cmd.Stderr = &stderr
		c.Cmd.Env = append(c.Cmd.Env, c.Env...)
	}

	if c.Cmd.Dir == "" {
		c.Cmd.Dir = c.ExecDir
	}

	log.Printf("-> Running %s %s from %s\n", c.BinPath, strings.Join(c.CmdArgs, " "), c.Cmd.Dir)
	err := c.Cmd.Run()
	res.Stderr = stderr.String()
	res.Stdout = stdout.String()
	if err != nil {
		res.Err = err
		return res
	}

	if c.ManifestDir != "" {
		if !util.PathExists(c.ManifestDir) {
			err := util.DirInit(c.ManifestDir)
			if err != nil {
				// This is not a fatal error, we log it and exit
				log.Printf("failed to create destination directory for the manifest: %s", err)
				return res
			}
		}
		path := filepath.Join(c.ManifestDir, "exec.MANIFEST")
		if c.ManifestName != "" {
			path = filepath.Join(c.ManifestDir, c.ManifestName+".MANIFEST")
		}
		if !util.FileExists(path) {
			currentTime := time.Now()
			data := []string{"Command: " + c.BinPath + " " + strings.Join(c.CmdArgs, " ") + "\n"}
			data = append(data, "Execution path: "+c.ExecDir)
			data = append(data, "Execution time: "+currentTime.Format("2006-01-02 15:04:05"))
			data = append(data, c.ManifestData...)

			// We transform relative paths into absolute path
			if c.BinPath[0] == '.' && c.BinPath[1] == '/' {
				c.BinPath = filepath.Join(c.ExecDir, c.BinPath[2:])
			}
			filesToHash := []string{c.BinPath} // we always get the fingerprint of the binary we execute
			filesToHash = append(filesToHash, c.ManifestFileHash...)
			hashData := manifest.HashFiles(filesToHash)
			data = append(data, hashData...)

			err := manifest.Create(path, data)
			if err != nil {
				// This is not a fatal error, we just log it
				log.Printf("failed to create manifest: %s", err)
			}
			log.Printf("-> Manifest successfully created (%s)", path)

		} else {
			log.Printf("Manifest %s already exists, skipping...", err)
		}
	}

	return res
}
