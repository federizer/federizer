package federizer

import (
	"federizer/config"
	"io"
	"log"
	"net/http"
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
	cfg := config.Config{}
	if err := cfg.LoadConfig(); err != nil {
		return err
	}

	http.HandleFunc("/", helloWorldHandler)

	log.Printf("Starting server on %s:%d\n", cfg.ServerHost, cfg.ServerPort)
	return http.ListenAndServe(cfg.ServerHost+":"+strconv.Itoa(cfg.ServerPort), nil)
}
