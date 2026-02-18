package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:42069")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Println("Connection Accepted")
		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
		fmt.Println("Connection Closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	fileLines := make(chan string)
	messageContent := ""
	readLength := 8
	go func() {
		for {
			fileContent := make([]byte, readLength)
			numberOfBytesRead, err := f.Read(fileContent)

			contentStr := string(fileContent)
			parts := strings.Split(contentStr, "\n")
			if len(parts) > 1 {
				messageContent += parts[0]
				fileLines <- messageContent
				messageContent = ""
				if len(parts) > 2 {
					for _, part := range parts[1:] {
						if part != parts[(len(parts)-1)] {
							fileLines <- messageContent
						} else if part != parts[(len(parts)-1)] {
							messageContent += part
						}
					}
				} else {
					messageContent += parts[1]
				}
			} else {
				messageContent += parts[0]
			}

			if err != nil {
				if (errors.Is(err, io.EOF)) || numberOfBytesRead < readLength {
					f.Close()
					close(fileLines)
					return
				} else {
					log.Fatalf("%v", err)
				}
			}
		}
	}()

	return fileLines
}
