package federizer

import (
	"federizer/internal/shared/config"
	"io"
	"log"
	"net/http"
	"os"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		io.WriteString(w, "Hello, World!")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Start() error {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Error reading config file: %v\n", err)
		return err
	}

	var cfg config.Config
	if err := cfg.Load(configFile); err != nil {
		log.Printf("Error loading config: %v\n", err)
		return err
	}

	http.HandleFunc("/", helloWorldHandler)

	log.Printf("Starting server on %s:%s\n", cfg.ServerHost, cfg.ServerPort)
	return http.ListenAndServe(cfg.ServerHost+":"+cfg.ServerPort, nil)
}
