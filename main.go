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

func PullUp(dir string) {
	// All files in original directory
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	bad_exts := [4]string{".exe", ".jpg", ".nfo", ".txt"}
	// Only for sub-directories
	for _, file := range files {
		if file.IsDir() {
			println("dir: ", file.Name())
			sub_dir := dir + file.Name()
			osappend(&sub_dir)
			PullUp(sub_dir)
			// DeleteDir(sub_dir)
		} else {
			// Individual Files in sub-directories
			for _, temp := range files {
				println(temp.Name())
			}
			for _, sub_file := range files {
				source_file := dir + sub_file.Name()
				dest_file := root_dir + sub_file.Name()
				if sub_file.IsDir() {
					println("Sub file: ", sub_file.Name(), " is a dir")
					osappend(&source_file)
					PullUp(source_file)
					DeleteDir(source_file)
					continue
				} else if CheckExtensions(sub_file, bad_exts[:]) {
					os.Remove(source_file)
					println("DELETED FILE: ", source_file)
					continue
				} else {
					if source_file != dest_file {
						err := MoveFile(source_file, dest_file)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						print("Source and dest are the same")
					}
					continue
				}
			}
		}
	}
}

var root_dir string

func main() {

	// root_dir = ParseArgs(root_dir)
	root_dir = `C:\Users\walter\Downloads\backup - Copy\`
	PullUp(root_dir)
}
