package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	snapshot "client/internal"
)

const baseport = ":9000"
const storage = "../../files/"

type Data struct {
	Source  string
	ID      uint64
	Payload string
}

func main() {
	//dirs creation
	if err := os.MkdirAll(storage, 0o755); err != nil {
		log.Fatalf("can't create storage dir: %v", err)
	}
	if err := os.MkdirAll(storage+"/snapshot", 0o755); err != nil {
		log.Fatalf("can't create snapshot dir: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/get/", getHandler)
	mux.HandleFunc("/replace/", putHandler)

	srv := &http.Server{
		Addr:    baseport,
		Handler: mux,
	}
	log.Printf("listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/get/")
	file, err := os.ReadFile(storage + name)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Printf("read error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(file)))
	if _, err := w.Write(file); err != nil {
		log.Printf("write error: %v", err)
	}
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/replace/")
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var data Data
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		log.Printf("json parse error: %v", err)
		return
	}

	if err := os.WriteFile(storage+name, []byte(data.Payload), 0o644); err != nil {
		http.Error(w, "cannot write file", http.StatusInternalServerError)
		log.Printf("write error: %v", err)
		return
	}

	// start snapshot
	go snapshot.Snap(name)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("write successful"))
}
