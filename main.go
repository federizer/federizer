package main

import (
	federizer "federizer/cmd"
	"log"
)

func main() {
	err := federizer.Start()
	if err != nil {
		log.Fatalf("Server encountered an error: %v\n", err)
	}
	log.Printf("Server shutdown gracefully")
}