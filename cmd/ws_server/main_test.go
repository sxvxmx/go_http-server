package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPutAndGetSnapshotFlow(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/replace/", putHandler)
	mux.HandleFunc("/get/", getHandler)
	server := httptest.NewServer(mux)
	defer server.Close()

	//payload
	payload := map[string]string{"Payload": "hello from test3"}
	bodyBytes, _ := json.Marshal(payload)
	putURL := server.URL + "/replace/test3"

	// send PUT request
	req, err := http.NewRequest(http.MethodPut, putURL, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("create put request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do put request: %v", err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from put, got %d", resp.StatusCode)
	}

	time.Sleep(10 * time.Second)

	snapName := "snapshots/test3.snap"
	getURL := server.URL + "/get/" + snapName

	// send GET request to fetch the snap file
	resp, err = client.Get(getURL)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from get, got %d", resp.StatusCode)
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read get body: %v", err)
	}

	// check contents
	if !bytes.Contains(got, []byte("hello from test3")) {
		t.Fatalf("unexpected snapshot content: %s", string(got))
	}
}
