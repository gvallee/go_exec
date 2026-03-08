// Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsSHA256Hash(t *testing.T) {
	valid := strings.Repeat("a", 64)
	if !isSHA256Hash(valid) {
		t.Fatalf("expected valid sha256 hash")
	}

	if isSHA256Hash("ABC") {
		t.Fatalf("expected invalid hash")
	}

	if isSHA256Hash(strings.Repeat("g", 64)) {
		t.Fatalf("expected invalid hash")
	}
}

func TestHashFiles(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "sample.txt")
	content := []byte("hello world")
	err := os.WriteFile(file, content, 0644)
	if err != nil {
		t.Fatalf("failed to create file: %s", err)
	}

	h := sha256.Sum256(content)
	expected := hex.EncodeToString(h[:])

	data := HashFiles([]string{file})
	if len(data) != 1 {
		t.Fatalf("unexpected hash count: %d", len(data))
	}

	if data[0] != file+": "+expected {
		t.Fatalf("unexpected hash data: %q", data[0])
	}
}

func TestCreateLoadAndCheck(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	err := os.WriteFile(target, []byte("original"), 0644)
	if err != nil {
		t.Fatalf("failed to create target file: %s", err)
	}

	manifestPath := filepath.Join(dir, "exec.MANIFEST")
	entries := HashFiles([]string{target})
	err = Create(manifestPath, entries)
	if err != nil {
		t.Fatalf("failed to create manifest: %s", err)
	}

	dataByFile, dataByHash, err := Load(manifestPath)
	if err != nil {
		t.Fatalf("failed to load manifest: %s", err)
	}

	if len(dataByFile) != 1 || len(dataByHash) != 1 {
		t.Fatalf("unexpected manifest content size: %d/%d", len(dataByFile), len(dataByHash))
	}

	if err = Check(manifestPath); err != nil {
		t.Fatalf("check failed for unchanged file: %s", err)
	}

	err = os.WriteFile(target, []byte("changed"), 0644)
	if err != nil {
		t.Fatalf("failed to modify target file: %s", err)
	}

	if err = Check(manifestPath); err == nil {
		t.Fatalf("expected hash mismatch after file modification")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, _, err := Load(filepath.Join(t.TempDir(), "missing.MANIFEST"))
	if err == nil {
		t.Fatalf("expected error for missing manifest")
	}
}
