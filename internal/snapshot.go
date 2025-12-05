package internal

import (
	"io"
	"log"
	"os"
)

const storage = "../../files/"

func Snap(file string) {
	log.Printf("started with %q", file)
	copyFile(storage+file, storage+"snapshots/"+file+".snap")
}

func copyFile(src, dst string) {
	sourceFile, err := os.Open(src)
	if err != nil {
		log.Printf("failed to open source file %q: %v", src, err)
		return
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		log.Printf("failed to create destination file %q: %v", dst, err)
		return
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		destinationFile.Close()
		log.Printf("copy content failed from %q to %q: %v", src, dst, err)
		return
	}
}
