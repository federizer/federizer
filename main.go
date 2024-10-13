package main

import (
	federizer "federizer/cmd"
	"log"
)

func main() {
	err := federizer.Start()
	if err != nil {
		log.Fatalf("federizer error: %v", err)
	}
	log.Print("federizer shutdown gracefully")
}
