//
// Copyright (c) 2023-2026, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE file for license information
//

package remote

import (
	"os/exec"
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
