package main

import (
	"flag"
	"os"
	"path/filepath"
	"os/exec"
	"log"
	"strings"
	"fmt"
)

const (
	MAX_DIRITEMS_READ=100000
)

var (
	GoCommand = flag.String("go-command", "go","command for run go compiler")
	GoRoot = flag.String("go-path", os.Getenv("GOROOT"), "path for goroot")
	Environments = flag.String("envs", "windows/amd64,windows/386,linux/amd64,linux/386", "Comma separated OS/ARCH for build go env.")
)

func main(){
	flag.Parse()

	srcPrefix := filepath.Join(*GoRoot, "src")
	srcCmdPrefix := filepath.Join(srcPrefix, "cmd")
	builded := []string{}
	for _, env := range strings.Split(*Environments, ","){
		env = strings.TrimSpace(env)
		parts := strings.Split(env,"/")
		if len(parts) != 2 {
			log.Println("Error while split OS/ARCH:", env)
			continue
		}

		pkgPrefix := filepath.Join(*GoRoot, "pkg", parts[0] + "_" + parts[1])

		log.Println("Build for: ", env)
		os.Setenv("GOOS", parts[0])
		os.Setenv("GOARCH", parts[1])

		filepath.Walk(srcPrefix, func(path string, info os.FileInfo, err error) error{
			if path == srcPrefix {
				return nil
			}

			if path == srcCmdPrefix {
				return filepath.SkipDir
			}

			if !info.IsDir() {
				return nil
			}


			if dir, err := os.Open(path); err == nil {
				files, err := dir.Readdir(MAX_DIRITEMS_READ)
				if err != nil {
					log.Println("Can't")
				}
			}else{
				log.Println("Can't open dir:", path, err)
			}

			base := filepath.Base(path)
			if base == "internal" || base == "vendor" {
				log.Println("Skip personal build: ", path)
				return filepath.SkipDir
			}

			packageName := path[len(srcPrefix)+1:]
			pkgFileName := filepath.Join(pkgPrefix, packageName) + ".a"
			if _, err := os.Stat(pkgFileName); err == nil {
				log.Println("File already exists:", pkgFileName)
				return nil
			}
			cmd := exec.Command(*GoCommand, "install", packageName)
			log.Println(cmd.Args)
			cmd.Run()
			return nil
		})

		builded = append(builded, env)
	}
	fmt.Println("Builded for:", builded)
}
