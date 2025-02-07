package main

import (
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"log"
	"os"
	"path"
	"sync"
)

var (
	wg sync.WaitGroup
)

func walk(dest, dir string) {
	defer wg.Done()
	if files, err := os.ReadDir(dir); err == nil {
		for _, file := range files {
			fullpath := path.Join(dir, file.Name())
			if file.IsDir() {
				wg.Add(1)
				go walk(dest, fullpath)
			} else if file.Type().IsRegular() {
				// hash, move, link
				f, err := os.Open(fullpath)
				if err != nil {
					log.Println("open", err)
					break
				}
				s := sha512.New()
				if _, err := f.WriteTo(s); err != nil {
					log.Println("hash", err)
					break
				}
				f.Close()
				filesum := hex.EncodeToString(s.Sum(nil))
				log.Println(fullpath, filesum)
				os.Rename(fullpath, path.Join(dest, filesum))
				os.Symlink(path.Join(dest, filesum), fullpath)
			}
		}
	}
}

func main() {
	clipath := flag.String("path", "/tmp/links", "Path to turn into symlinks.")
	clidest := flag.String("dest", "/tmp/hashes", "Path to store hash content.")
	flag.Parse()

	wg.Add(1)
	go walk(*clidest, *clipath)
	wg.Wait()
}
