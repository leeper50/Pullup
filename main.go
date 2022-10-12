package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
)

// Adds correct slash for path depending on OS
func osappend(dir *string) {
	os := runtime.GOOS
	switch os {
	case "windows":
		if !strings.HasSuffix(*dir, `\`) {
			*dir += `\`
		}
	case "macos":
		fallthrough
	case "linux":
		if !strings.HasSuffix(*dir, "/") {
			*dir += "/"
		}
	}
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
	// Directory path is checked before run
	files, _ := os.ReadDir(dir)
	extensions := []string{".exe", ".jpg", ".nfo", ".png", ".txt"}
	for _, file := range files {
		if file.IsDir() {
			println("Directory: ", file.Name())
			sub_dir := dir + file.Name()
			osappend(&sub_dir)
			PullUp(sub_dir)
			os.RemoveAll(sub_dir)
			println("Deleted directory: ", sub_dir)
		} else {
			source_file := dir + file.Name()
			ext := filepath.Ext(file.Name())
			if slices.Contains(extensions, ext) {
				os.Remove(source_file)
				println("Deleted file: ", source_file)
			} else {
				dest_file := root_dir + file.Name()
				if source_file != dest_file {
					err := os.Rename(source_file, dest_file)
					if err != nil {
						log.Fatal(err)
					}
					println(source_file, " --> ", dest_file)
				} else {
					println("File already at root", source_file)
				}
			}
		}
	}
}

var root_dir string

func main() {
	root_dir = ParseArgs(root_dir)
	PullUp(root_dir)
}
