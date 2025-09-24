package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				// Send any remaining data as the last line if there's content
				if len(str) > 0 {
					out <- str
				}
				break
			}

			data = data[:n]
			str += string(data)

			// Check for complete lines (ending with \n)
			if i := bytes.IndexByte([]byte(str), '\n'); i != -1 {
				// Found a complete line
				line := str[:i]
				out <- line
				str = str[i+1:] // Remove the processed line including \n
			}
		}
	}()

	return out
}

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
