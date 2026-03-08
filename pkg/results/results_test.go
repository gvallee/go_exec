// Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package results

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingFile(t *testing.T) {
	data, err := Load(filepath.Join(t.TempDir(), "missing.out"))
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %s", err)
	}

	if len(data) != 0 {
		t.Fatalf("expected no result, got %d", len(data))
	}
}

func TestLoadParseResults(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "results.out")
	content := "test1 PASS\n\n test2 FAIL\ntest3    PASS   note\n"
	err := os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create results file: %s", err)
	}

	data, err := Load(outputFile)
	if err != nil {
		t.Fatalf("failed to load output file: %s", err)
	}

	if len(data) != 3 {
		t.Fatalf("expected 3 results, got %d", len(data))
	}

	if !data[0].Pass {
		t.Fatalf("expected first result to pass")
	}

	if data[1].Pass {
		t.Fatalf("expected second result to fail")
	}

	if !data[2].Pass {
		t.Fatalf("expected third result to pass")
	}
}
