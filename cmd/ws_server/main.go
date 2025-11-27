package main

import (
	"log"
	"net"

	handler_request "client/internal"
)

const baseport = ":9000"

func main() {
	for {
		listener, err := net.Listen("tcp", baseport)
		if err != nil {
			log.Printf("error while sw server starting: %v", err)

			return
		}

		defer listener.Close()

		log.Printf("Server is up on port%v", baseport)

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("error while adding new connection: %v", err)

				continue
			}

			go handler_request.Handler(conn)
		}
	}
}
