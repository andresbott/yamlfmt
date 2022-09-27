// ADOBE CONFIDENTIAL
// ___________________
//
// Copyright 2022 Adobe
// All Rights Reserved.
//
// NOTICE: All information contained herein is, and remains
// the property of Adobe and its suppliers, if any. The intellectual
// and technical concepts contained herein are proprietary to Adobe
// and its suppliers and are protected by all applicable intellectual
// property laws, including trade secret and copyright laws.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Adobe.

package cli

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/andresbott/yamlfmt/internal/filematch"
	"github.com/andresbott/yamlfmt/internal/yamlfmt"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

// Execute is the entry point for the command line
func Execute() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	ShaVer    string // sha1 revision used to build the program
	BuildTime string // when the executable was built
	Version   = "development"
)

const name = "yamlfmt"

func rootCmd() *cobra.Command {
	var dry bool
	var verbose bool
	var quiet bool
	var version bool

	cmd := &cobra.Command{
		Use: name + " <glob>",
		Long: name + " format yaml files to a opinionated defaults, it is inspired ing go fmt. \n" +
			"The input accepts any regular glob patterns plus some convenience aliases:" +
			"\n- \"./\" is an alias for: \"./*.yaml\" and \"./*.yml\"" +
			"\n- \"./..\" is an alias for: \"./**/*.yaml\" and \"./**/*.yml\" to search recursively" +
			"\n- any existing directory will be searched for \"./*.yaml\" and \"./*.yml\"",

		Short: "format yaml files",
		RunE: func(cmd *cobra.Command, args []string) error {

			if version {
				fmt.Printf("Version: %s\n", Version)
				fmt.Printf("Build date: %s\n", BuildTime)
				fmt.Printf("Commit sha: %s\n", ShaVer)
				fmt.Printf("Compiler: %s\n", runtime.Version())

				return nil
			}

			pattern := ""
			if len(args) > 0 {
				pattern = args[0]

			}

			files, err := filematch.FindFiles(pattern)
			if err != nil {
				return fmt.Errorf("unable to find files: %v", err)
			}

			p := printer{
				quiet:   quiet,
				verbose: verbose,
			}

			if len(files) == 0 {
				p.print("no files found")
				return nil
			}

			if dry {
				p.print("(running in dry run mode, no changes will be written)")
			}

			for _, f := range files {

				abs, err := filepath.Abs(f)
				if err != nil {
					return fmt.Errorf("unable to get absolute path for: %s,  %s", f, err.Error())
				}
				fInfo, err := os.Stat(abs)
				if err != nil {
					return fmt.Errorf("unable to stat file: %s,  %s", f, err.Error())
				}
				if fInfo.IsDir() {
					return fmt.Errorf("%s is a directoru:  %s", f, err.Error())
				}

				maxFileSize := int64(50 * 1024 * 1024)
				if fInfo.Size() > maxFileSize {
					return fmt.Errorf("file %s is bigger than 50MB", f)
				}

				inBytes, err := os.ReadFile(abs)
				if err != nil {
					return fmt.Errorf("unable to read file: %s", f)
				}

				inHash := fileHash(inBytes)
				got, err := yamlfmt.Format(inBytes)
				if err != nil {
					return fmt.Errorf("unable to format file %s:%s", f, err.Error())
				}
				outHash := fileHash(got)

				if inHash != outHash {
					p.print(f)
					diff := cmp.Diff(string(inBytes), string(got))
					p.printVerbose(diff)

					if !dry {
						err = os.WriteFile(abs, got, fInfo.Mode())
						if err != nil {
							return fmt.Errorf("unable to write file: %s", f)
						}
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&dry, "dry-run", "d", false, "do not persists changes")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print the diff for every file to change")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "dont print any output, takes precedence over verbose")
	cmd.Flags().BoolVar(&version, "version", false, "print version and build information")

	return cmd
}

func fileHash(in []byte) string {
	r := bytes.NewReader(in)

	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
