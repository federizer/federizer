package federizer

import (
	"federizer/config"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		io.WriteString(w, "Hello, World!")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Start() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("Error reading config file %s: %v\n", configPath, err)
		return err
	}

	var cfg config.Config
	if err := cfg.Load(configFile); err != nil {
		log.Printf("Error loading config: %v\n", err)
		return err
	}

	http.HandleFunc("/", helloWorldHandler)

	log.Printf("Starting server on %s:%d\n", cfg.ServerHost, cfg.ServerPort)
	return http.ListenAndServe(cfg.ServerHost+":"+strconv.Itoa(cfg.ServerPort), nil)
}
