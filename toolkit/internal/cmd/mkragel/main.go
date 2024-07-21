// Mkragel is a small helper to pretty-up ragel output.
//
// This helper is needed to add the magic "generated file" comment, as the
// ragel generator does not add it and strips comments before a `package`
// declaration.
//
// Each file provided as an argument is generated separately. Use the "native"
// ragel include/import mechanisms to handle machines split across multiple
// files.
//
// This is tested/generated with ragel shipped in Fedora 39 (ragel 7.0.4). If
// using ragel 6.x, the `-Z` flag to the `ragel` binary should be used instead
// of a `ragel-go` command.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	debugFlag := flag.Bool("D", false, "debug output")
	flag.Parse()

	if *debugFlag {
		debug = debugOutput
	}
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "no input files specified")
		os.Exit(1)
	}

	for _, name := range flag.Args() {
		if err := runOne(name); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func runOne(name string) error {
	dir, err := os.MkdirTemp("", "mkragel.")
	if err != nil {
		return err
	}
	defer func() { os.RemoveAll(dir) }()
	debug("work directory: %q", dir)
	gensrc := filepath.Join(dir, "out.go")
	debug("generated output: %q", gensrc)

	cmd := exec.Command("ragel-go", "--string-tables", "-F1", "-o", gensrc, name)
	debug("exec: %q", cmd.Args)
	if err := cmd.Run(); err != nil {
		return err
	}
	gen, err := os.Open(gensrc)
	if err != nil {
		return err
	}
	defer gen.Close()

	var buf bytes.Buffer
	_, err = fmt.Fprintf(&buf, "// Code generated by ragel-go from %s. DO NOT EDIT.\n\n", name)
	if err != nil {
		return err
	}
	if _, err := io.Copy(&buf, gen); err != nil {
		return err
	}

	fmtd, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	outname := strings.Replace(name, ".rl", ".go", 1)
	debug("final output: %q", outname)
	out, err := os.Create(outname)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, bytes.NewReader(fmtd)); err != nil {
		return err
	}

	return nil
}

var debug func(string, ...any) = debugNoOp

func debugNoOp(_ string, _ ...any) {}

func debugOutput(f string, v ...any) {
	if f[len(f)-1] != '\n' {
		f = f + "\n"
	}
	fmt.Fprintf(os.Stderr, f, v...)
}