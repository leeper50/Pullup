package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"runtime"
	"strings"
)

func CheckExtensions(file fs.DirEntry, extensions []string) bool {
	if file.IsDir() {
		log.Fatal("can not check extension of a folder")
	}
	for _, ext := range extensions {
		if strings.Contains(file.Name(), ext) {
			return true
		}
	}
	return false
}

func osappend(dir *string) {
	os := runtime.GOOS
	switch os {
	case "windows":
		if !strings.HasSuffix(*dir, `\`) {
			*dir += `\`
		}
	case "linux":
		if !strings.HasSuffix(*dir, "/") {
			*dir += "/"
		}
	}
}

func DeleteDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatal(err)
	}
	println("DELETED DIR: ", dir)
}

func MoveFile(source_path, dest_path string) error {
	err := os.Rename(source_path, dest_path)
	if err != nil {
		return fmt.Errorf("failed moving file: %s", err)
	}
	println(source_path, " --> ", dest_path)
	return nil
}

func ParseArgs(path string) string {
	flag.StringVar(&path, "p", ``, "")
	flag.StringVar(&path, "path", ``, "")
	flag.Usage = func() {
		fmt.Printf("You must provide a path with the syntax:\n")
		fmt.Printf("./main.go -p path\n")
		fmt.Printf("./main.go --path path\n")
	}
	flag.Parse()
	if path == "" {
		var user_input string
		for invalid_path := true; invalid_path; {
			print("Enter a valid path\n")
			fmt.Scanln(&user_input)
			_, err := os.ReadDir(user_input)
			if err == nil {
				break
			}
		}
	}
	osappend(&path)
	println("Output path: " + path)
	return path
}

func PullUp(files []fs.DirEntry, dir string) {
	var (
		bad_exts  = [4]string{".exe", ".jpg", ".nfo", ".txt"}
		err       error
		filepath  string
		sub_files []fs.DirEntry
	)
	// Only for sub-directories
	for _, file := range files {
		if file.IsDir() {
			filepath = dir + file.Name()
			osappend(&filepath)
			sub_files, err = os.ReadDir(filepath)
			if err != nil {
				log.Fatal(err)
			}
			// Individual Files in sub-directories
			for _, sub_file := range sub_files {
				source_file := filepath + sub_file.Name()
				file_dest := dir + sub_file.Name()
				if !sub_file.IsDir() && !CheckExtensions(sub_file, bad_exts[:]) {
					err = MoveFile(source_file, file_dest)
					if err != nil {
						log.Fatal(err)
					}
				} else if sub_file.IsDir() {
					DeleteDir(source_file)
				} else {
					os.Remove(source_file)
					println("DELETED FILE: ", source_file)
					continue
				}
			}
			DeleteDir(filepath)
		}
	}
}

func main() {
	var (
		err        error
		root_dir   string
		root_files []fs.DirEntry
	)

	root_dir = ParseArgs(root_dir)

	// All files in original directory
	root_files, err = os.ReadDir(root_dir)
	if err != nil {
		log.Fatal(err)
	}

	PullUp(root_files, root_dir)
}
