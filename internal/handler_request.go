package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const storage = "../../files/"

type Data struct {
	Source  string
	ID      uint64
	Payload string
}

func Handler(conn net.Conn) {
	defer conn.Close()
	log.Printf("connection established")

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("error while reading data from user: %v", err)
		return
	}

	request := string(buf[:n])
	log.Printf("get data from tcp connection: %v", request)
	lines := strings.Split(request, "/r/n")
	if len(lines) == 0 {
		log.Println("invalid request")
		return
	}

	parts := strings.Fields(lines[0])
	switch method := parts[0]; method {
	case "GET":
		getHandler(conn, lines)
	case "PUT":
		putHandler(conn, lines)
	default:
		log.Printf("no method for %v", method)
	}

}

func resp(conn net.Conn, content []byte) {
	response := "HTTP/1.1 200 OK\r\nContent-Length: " + fmt.Sprint(len(content)) + "\r\n\r\n" + string(content)
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("error while writing request: %v", err)
	}
	log.Printf("request proceed successfully")
}

func getHandler(conn net.Conn, lines []string) {
	if strings.Contains(strings.Fields(lines[0])[1], "/get") {
		content, err := os.ReadFile(storage + strings.TrimPrefix(strings.Fields(lines[0])[1], "/get/"))
		if err != nil {
			response := "HTTP/1.1 404 Not Found\r\n\r\nFile Not Found"
			resp(conn, []byte(response))
			log.Printf("can not read file: %v", err)
			return
		}
		resp(conn, content)
	}
}

func putHandler(conn net.Conn, lines []string) {
	if strings.Contains(strings.Fields(lines[0])[1], "/replace") {
		var data Data
		jn := []byte(strings.SplitN(lines[0], "\r\n\r\n", 2)[1])

		err := json.Unmarshal(jn, &data)
		if err != nil {
			response := "HTTP/1.1 500"
			resp(conn, []byte(response))
			log.Printf("parse: %v", err)
			return
		}
		log.Printf("payload %v", data.Payload)
		err = os.WriteFile(storage+strings.TrimPrefix(strings.Fields(lines[0])[1], "/replace/"), []byte(data.Payload), os.FileMode(0644))
		if err != nil {
			response := "HTTP/1.1 500"
			resp(conn, []byte(response))
			log.Printf("can not write file: %v", err)
			return
		}
		resp(conn, []byte("write successful"))
	}
}
