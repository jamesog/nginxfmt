// Sections of this code are Copyright 2009 The Go Authors. All rights reserved.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/aluttik/go-crossplane"
)

var (
	list  = flag.Bool("l", false, "list files whose formatting differs from nginxfmt's")
	write = flag.Bool("w", false, "write result to (source) file instead of stdout")
)

const (
	indent = 2
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: nginxfmt [flags] [path ...]\n")
	flag.PrintDefaults()
}

func format(config crossplane.Config) ([]byte, error) {
	var buf bytes.Buffer
	err := crossplane.Build(&buf, config, &crossplane.BuildOptions{Indent: indent})
	return buf.Bytes(), err
}

// If in == nil, the source is the contents of the file with the given filename.
func processFile(filename string, info fs.FileInfo, in io.Reader, out io.Writer) error {
	src, err := readFile(filename, in)
	if err != nil {
		return err
	}

	p, err := crossplane.Parse(filename, &crossplane.ParseOptions{
		Open:               func(string) (io.Reader, error) { return bytes.NewReader(src), nil },
		ParseComments:      true,
		StopParsingOnError: true,
	})
	if err != nil {
		return err
	}

	res, err := format(p.Config[0])
	if err != nil {
		return err
	}

	if !bytes.Equal(src, res) {
		if *list {
			fmt.Println(filename)
		}
		if *write {
			perm := info.Mode().Perm()
			bakname := filename + ".bak"
			err = os.Rename(filename, bakname)
			if err != nil {
				return err
			}
			err = os.WriteFile(filename, res, perm)
			if err != nil {
				os.Rename(bakname, filename)
				return err
			}
			err = os.Remove(bakname)
			if err != nil {
				return err
			}
		}
	}

	if !*list && !*write {
		res = append(res, '\n')
		_, err = out.Write(res)
	}

	return err
}

func readFile(filename string, in io.Reader) ([]byte, error) {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		in = f
		defer f.Close()
	}
	src, err := io.ReadAll(in)
	return src, err
}

func main() {
	err := nginxfmtMain(os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func nginxfmtMain(out io.Writer) error {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		if *write {
			return fmt.Errorf("error: cannot use -w with standard input")
		}
		return processFile("<stdin>", nil, os.Stdin, out)
	}

	for _, arg := range args {
		switch info, err := os.Stat(arg); {
		case err != nil:
			return err
		case !info.IsDir():
			err := processFile(arg, info, nil, out)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error processing %s: %v\n", arg, err)
			}
		default:
			err := filepath.WalkDir(arg, func(path string, f fs.DirEntry, err error) error {
				if err != nil || f.IsDir() {
					return err
				}
				info, err := f.Info()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: error: %v", path, err)
					return nil
				}
				return processFile(path, info, nil, out)
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: error walking directory: %v", arg, err)
			}
		}
	}

	return nil
}
