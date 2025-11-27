package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

// only one at a time((
func main() {
	test1()
}

// file access test
func test1() {
	host := flag.String("server_host", "127.0.0.1", "des")
	port := flag.String("server_port", "9000", "des")
	file := flag.String("file", "get/test", "des")

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

// nonexistent file access test
func test2() {
	host := flag.String("server_host", "127.0.0.1", "des")
	port := flag.String("server_port", "9000", "des")
	file := flag.String("file", "get/test99", "des")

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

// correctness test PUT + payload write
func test3() {
	host := flag.String("server_host", "127.0.0.1", "des")
	port := flag.String("server_port", "9000", "des")
	file := flag.String("file", "replace/test2", "des")

	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", *host, *port))
	if err != nil {
		log.Printf("can not connect to host %v with port %v", *host, *port)
	}
	defer conn.Close()
	log.Printf("successfully connected")
	payload := `{"source":"example", "id":123, "payload":"Hello, GO v2"}`
	request := fmt.Sprintf("PUT /%s HTTP/1.1\r\nContent-Type: application/json\r\n\r\n%s", *file, payload)
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

// test with POST
func test4() {
	host := flag.String("server_host", "127.0.0.1", "des")
	port := flag.String("server_port", "9000", "des")
	file := flag.String("file", "replace/test2", "des")

	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", *host, *port))
	if err != nil {
		log.Printf("can not connect to host %v with port %v", *host, *port)
	}
	defer conn.Close()
	log.Printf("successfully connected")
	payload := `{"source":"example", "id":123, "payload":"Hello, GO v2"}`
	request := fmt.Sprintf("POST /%s HTTP/1.1\r\nContent-Type: application/json\r\n\r\n%s", *file, payload)
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
