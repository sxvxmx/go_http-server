package internal

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func Handler(conn net.Conn) {
	defer conn.Close()
	log.Printf("connection established")

	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		log.Printf("error while reading data from user: %v", err)

		return
	}

	request := string(buf)
	log.Printf("get data from tcp connection: %v", request)
	lines := strings.Split(request, "/r/n")
	if len(lines) == 0 {
		log.Println("invalid request")

		return
	}

	parts := strings.Fields(lines[0])
	content, err := os.ReadFile(strings.TrimPrefix(parts[1], "/"))
	if err != nil {
		response := "HTTP/1.1 404 Not Found\r\n\r\nFile Not Found"
		conn.Write([]byte(response))
		log.Printf("can not read file: %v", err)

		return
	}

	response := "HTTP/1.1 200 OK\r\nContent-Length: " + fmt.Sprint(len(content)) + "\r\n\r\n" + string(content)
	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("error while writing request: %v", err)
	}

	log.Printf("request proceed successfully")
}
