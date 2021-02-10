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

package main

import (
	"fmt"
	"gobuildercli/cmd"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	cmd.Execute()
	runCopy()
	//var out, stderr bytes.Buffer
	//	//cmd := exec.Command(fmt.Sprintf("%s/go build", goPath))
	//	//cmd.Stdout = &out
	//	//cmd.Stderr = &stderr
	//	//err := cmd.Run()
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//}

	//var sources []string
	//srcs, err := filepath.Glob("H:/file-handler/*")
	//if err != nil {
	//	fmt.Println("oh")
	//	fmt.Println(err.Error())
	//}
	//sources = append(sources, srcs...)

	//for _,v := range sources {
	//	println(v)
	//}
	//dst := "F:/0"
	//src := "H:/file-handler"
	//fmt.Println(src)
	//if err := CopyDir("/Users/taufiq/Documents/projects/asha/lib/models/",
	//	"/Users/taufiq/Documents/projects/test"); err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(FilterDirs("/Users/taufiq/Documents/projects/asha/lib/models", "*.dart"))
	//fmt.Println()
	//fmt.Println(FilterDirsGlob("/Users/taufiq/Documents/projects/asha/lib/models", "*.dart"))
	//fmt.Println()
	//fmt.Println(glob("/Users/taufiq/Documents/projects/asha/lib/models", "*.dart"))
	//fmt.Println()\

	// ====================================================================================
	//matches, _ := filepath.Glob("H:\\file-handler")
	//fmt.Println(matches)
	//for _, file := range matches {
	//fmt.Println(file)
	//	CopyDir("/Users/taufiq/Documents/projects/asha/lib", dst)
	//	//fmt.Println(file)
	//}
}

func runBuild() {
	cmd := exec.Command("go", "build", "./main.go")
	b, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	//cmd = exec.Command("go", "help")
	//b, err = cmd.Output()
	//
	//fmt.Println("string(b):", string(b))
	//cmd.Stdout = &out
	//cmd.Stderr = &stderr
	//if err != nil {
	//	fmt.Println("err")
	//	log.Fatal(err)
	//}
}

func runCopy() {
	cmd := exec.Command("cp", "./cmd", "./hello")
	b, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	//cmd = exec.Command("go", "help")
	//b, err = cmd.Output()
	//
	//fmt.Println("string(b):", string(b))
	//cmd.Stdout = &out
	//cmd.Stderr = &stderr
	//if err != nil {
	//	fmt.Println("err")
	//	log.Fatal(err)
	//}
}

func FilterDirs(dir, suffix string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), suffix) {
			res = append(res, filepath.Join(dir, f.Name()))
		}
	}
	return res, nil
}

func glob(dir string, ext string) ([]string, error) {

	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})

	return files, err
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

func FilterDirsGlob(dir, suffix string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, suffix))
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}

	//if !si.IsDir() {
	//	fmt.Errorf("source is not a directory")
	//	fmt.Println(si.Name())
	//	fmt.Println(src)
	//	fmt.Println(dst)
	//	return CopyFile(src, dst)
	//}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("156", err)
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		fmt.Println("165", err)
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		fmt.Println("171", err)

	}

	for _, entry := range entries {
		fmt.Println("entry.Name()", entry.Name())
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		fmt.Println(srcPath)
		fmt.Println(dstPath)
		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				fmt.Println("184", err)
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			fmt.Println("copy file")
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				fmt.Println("195", err)
				return
			}
		}
	}

	return
}
