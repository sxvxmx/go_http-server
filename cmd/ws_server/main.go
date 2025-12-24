package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	snapshot "client/internal"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

const baseport = ":9000"
const storage = "./files/"

type Data struct {
	Source  string
	ID      uint64
	Payload string
}

var snaps = make(map[string]bool)
var table = make(map[string]int)

// ADD USERS HERE!
var users = []string{}

func main() {
	//dirs creation
	if err := os.MkdirAll(storage, 0o755); err != nil {
		log.Fatalf("can't create storage dir: %v", err)
	}
	if err := os.MkdirAll(storage+"snapshots", 0o755); err != nil {
		log.Fatalf("can't create snapshots dir: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/get/", getHandler)
	mux.HandleFunc("/test", testHandler)
	mux.HandleFunc("/bcast", bcastHandler)
	mux.HandleFunc("/vclock", clockHandler)
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
	if name == "" {
		name = "test.json"
	}
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
	if name == "" {
		name = "test.json"
	}
	path := storage + name
	log.Printf("%s", fmt.Sprint(path))

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
	val, ok := table[data.Source]
	if !ok {
		table[data.Source] = 0
		val = 0
	}
	if val < int(data.ID) {
		table[data.Source] = int(data.ID)
		//payload logic
		patch, err := jsonpatch.DecodePatch([]byte(data.Payload))
		if err != nil {
			http.Error(w, "invalid json patch", http.StatusBadRequest)
			return
		}
		log.Printf("%s", fmt.Sprint(patch))

		orig, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				orig = []byte(`{}`)
			} else {
				http.Error(w, "cannot read file", http.StatusInternalServerError)
				return
			}
		}

		modified, err := patch.Apply(orig)
		if err != nil {
			http.Error(w, "patch apply failed", http.StatusUnprocessableEntity)
			return
		}

		if err := os.WriteFile(path, modified, os.FileMode(0644)); err != nil {
			http.Error(w, "cannot write file", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "try bigger id", http.StatusBadRequest)
		return
	}

	_, ok = snaps[path]
	if !ok {
		// start snapshot
		go snapshot.Snap(name)
		snaps[path] = true
		log.Printf("started snaping for: %v", path)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("write successful"))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	file, err := os.ReadFile("index.html")
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "file index.html not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Printf("read error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", strconv.Itoa(len(file)))
	if _, err := w.Write(file); err != nil {
		log.Printf("write error: %v", err)
	}
}

func clockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if err := json.NewEncoder(w).Encode(table); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
	}
}

// send /replace
func bcastHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/replace/")
	if name == "" {
		name = "test.json"
	}
	path := storage + name
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "cannot read file", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}

	for _, user := range users {
		req, err := http.NewRequest(
			http.MethodPut,
			user,
			bytes.NewReader(data),
		)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
	}
}
