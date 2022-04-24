// Sections of this code are Copyright 2009 The Go Authors. All rights reserved.

package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

func TestProcessFile(t *testing.T) {
	match, err := filepath.Glob("testdata/*.input")
	if err != nil {
		t.Fatal(err)
	}

	for _, in := range match {
		out := strings.Replace(in, ".input", ".golden", 1)
		t.Run(in, func(t *testing.T) {
			var buf bytes.Buffer
			processFile(in, nil, nil, &buf)

			expected, err := os.ReadFile(out)
			if err != nil && !*update {
				t.Error(err)
			}

			if got := buf.Bytes(); !bytes.Equal(got, expected) {
				if *update {
					if err := os.WriteFile(out, got, 0666); err != nil {
						t.Error(err)
					}
					return
				}

				t.Errorf("nginxfmt %s != %s (see %s.nginxfmt)\n", in, out, in)
				if err := os.WriteFile(in+".nginxfmt", got, 0666); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
