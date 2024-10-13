package federizer

import (
	"federizer/internal/shared/config"
	"fmt"
	"io"
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
		fmt.Printf("Error reading config file: %v\n", err)
		return err
	}

	var cfg config.Config
	if err := cfg.Load(configFile); err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		return err
	}

	http.HandleFunc("/hello", helloWorldHandler)

	port := fmt.Sprintf(":%s", cfg.Port)

	fmt.Printf("Starting server at port %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}

	return err
}
