package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	FirstPort  = "8080"
	SecondPort = "8081"
)

type Request struct {
	XMLName  xml.Name `xml:"file"`
	FilePath string   `xml:",chardata"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("no files")
		os.Exit(1)
	}

	file1, file2 := os.Args[1], os.Args[2]
	if file1 == "" || file2 == "" {
		fmt.Println("no files")
		os.Exit(1)
	}

	firstCmd := exec.Command("./bin/t1", FirstPort)
	secondCmd := exec.Command("./bin/t1", SecondPort)

	// starting xml-rpc servers
	firstReader := startServer(firstCmd, FirstPort)
	secondReader := startServer(secondCmd, SecondPort)

	b1, err := getBytesFromServer(FirstPort, file1, firstReader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	b2, err := getBytesFromServer(SecondPort, file2, secondReader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ioutil.WriteFile("res.txt", xorBytes(b1, b2), fs.FileMode(0777))

	firstCmd.Process.Kill()
	secondCmd.Process.Kill()
}

func startServer(command *exec.Cmd, port string) io.Reader {
	reader, writer := io.Pipe()
	command.Stdout = writer
	if err := command.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	waitServer(port)
	return reader
}

func waitServer(port string) {
	for !isPortOPen(port) {
		time.Sleep(time.Second / 2)
	}
}

func isPortOPen(port string) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func getBytesFromServer(port string, file string, reader io.Reader) ([]byte, error) {
	req := Request{FilePath: file}
	body, _ := xml.Marshal(req)
	_, err := http.Post(fmt.Sprintf("http://localhost:%s/rpc/print", port), "text/xml", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	data := make([]byte, 0, 1024)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		data = append(data, buf[:n]...)
		if n != 1024 {
			break
		}

	}

	return data, nil
}

func xorBytes(text, key []byte) []byte {
	textLen := len(text)
	dif := textLen - len(key)
	if dif > 0 {
		for i := 0; i < dif; i++ {
			key = append(key, 0)
		}
	}

	result := make([]byte, textLen)
	for i := 0; i < textLen; i++ {
		result[i] = key[i] ^ text[i]
	}

	return result
}
