// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package results

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// Result represents the result of a given experiment
type Result struct {
	Pass bool
	Note string
}

// Load reads a output file and load the list of experiments that are in the file
func Load(outputFile string) ([]Result, error) {
	var existingResults []Result

	log.Println("Reading results from", outputFile)

	f, err := os.Open(outputFile)
	if err != nil {
		// No result file, it is okay
		return existingResults, nil
	}
	defer f.Close()

	lineReader := bufio.NewScanner(f)
	if lineReader == nil {
		return existingResults, fmt.Errorf("failed to create file reader")
	}

	for lineReader.Scan() {
		line := lineReader.Text()
		var newResult Result
		if strings.Contains(line, "PASS") {
			newResult.Pass = true
		} else {
			newResult.Pass = false
		}
		existingResults = append(existingResults, newResult)
	}

	return existingResults, nil
}
