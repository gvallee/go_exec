//
// Copyright (c) 2023, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

package remote

import "testing"

func TestRemoteCmd(t *testing.T) {
	cmd := "/usr/bin/date"
	host := "localhost"
	err := ExecCmd("localhost", cmd, nil, nil)
	if err.Err != nil {
		t.Fatalf("unabel to run %s on %s", cmd, host)
	}
}
