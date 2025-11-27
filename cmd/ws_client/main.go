package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	host := flag.String("server_host", "127.0.0.1", "des")
	port := flag.String("server_port", "9000", "des")
	file := flag.String("file", "files/example_file.txt", "des")

	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", *host, *port))
	if err != nil {
		log.Printf("can not connect to host %v with port %v", *host, *port)
	}
	defer conn.Close()
	log.Printf("successfully connected")
	request := fmt.Sprintf("GET /%s HTTP/1.1\r\nHost: %s\r\n\r\n", *file, fmt.Sprintf("%v:%v", *host, *port))
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("successfully send mes")
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		log.Printf("error while reading data from user: %v", err)

		return
	}
	response := string(buf)
	fmt.Println(response)
}
