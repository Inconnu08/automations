// Copyright © 2021 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var buildexecuteCmd = &cobra.Command{
	Use:   "buildexecute",
	Short: "A brief description of builddir",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		buildDir, _ := cmd.Flags().GetString("builddir") // copydir will use this
		copyDir, _ := cmd.Flags().GetString("copydir")
		exeName, _ := cmd.Flags().GetString("exe")

		if exeName == "" {
			log.Fatal("Exe name not provided")
		}

		buildexecuteArgs(buildDir, copyDir, exeName)
	},
}

func init() {
	rootCmd.AddCommand(buildexecuteCmd)
	buildexecuteCmd.Flags().StringP("builddir", "b", "./", "directory specified for build")
	buildexecuteCmd.Flags().StringP("copydir", "c", "./", "copy directory")
	buildexecuteCmd.Flags().StringP("exe", "e", "", "directory specified for build")
}

func buildexecuteArgs(buildDir, copyDir, exeName string) {
	//fmt.Printf("Addition of %s", args)
	//var sources []string
	//fmt.Println(args[1])
	//srcs, err := filepath.Glob(args[1])
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//sources = append(sources, srcs...)
	//println(sources)

	if err := CopyDir(copyDir, buildDir); err != nil {
		log.Fatalln(err)
	}

	runBuild(exeName, buildDir)
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {

	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	if src == dst {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	//if err == nil {
	//	return fmt.Errorf("destination already exists")
	//}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func runBuild(name, buildDir string) {
	buildPath := buildDir + string(filepath.Separator) + name + ".exe"

	cmd := exec.Command("go", "build", "-o", buildPath)
	//cmd.Dir = "./"+buildDir
	b, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
