package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 2 {
		log.Fatal("Argument missing")
	}

	dir := os.Args[1]
	customEnv, err := ReadDir(dir)
	if err != nil {
		log.Fatalf("Error processing dir '%s': %s", dir, err)
	}

	code, err := RunCmd(os.Args[2:], customEnv)
	if err != nil {
		log.Fatalf("Error running command: %s", err)
	}

	os.Exit(code)
}
