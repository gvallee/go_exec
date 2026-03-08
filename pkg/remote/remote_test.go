//
// Copyright (c) 2023-2026, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE file for license information
//

package remote

import "testing"

func TestRemoteCmd(t *testing.T) {
	cmd := "/usr/bin/date"
	host := "localhost"
	err := ExecCmd(host, cmd, nil, nil)
	if err.Err != nil {
		t.Fatalf("unable to run %s on %s", cmd, host)
	}
}
