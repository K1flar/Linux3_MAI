package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
)

func main() {
	// if len(os.Args) < 3 {
	// 	os.Exit(1)
	// }

	// file1, file2 := os.Args[1], os.Args[2]
	// if file1 == "" || file2 == "" {
	// 	os.Exit(1)
	// }

	// Удаляем существующий сокет (если такой существует)
	os.Remove("/tmp/mysocket")

	// Создаем сокет для прослушивания
	addr, err := net.ResolveUnixAddr("unixgram", "/tmp/mysocket")
	if err != nil {
		fmt.Println("Ошибка при разрешении адреса:", err)
		return
	}

	// Слушаем сокет
	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		fmt.Println("Ошибка при прослушивании сокета:", err)
		return
	}
	defer os.Remove("/tmp/mysocket") // Удаляем сокет после завершения
	firstCmd := exec.Command("./bin/t1", "aboba")
	secondCmd := exec.Command("./bin/t1", "hui")
	if err := firstCmd.Run(); err != nil {
		fmt.Println(err)
	}
	if err := secondCmd.Run(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Ожидание данных...")

	// Читаем данные
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUnix(buf)
	if err != nil {
		fmt.Println("Ошибка при чтении данных:", err)
		return
	}

	fmt.Printf("Получено %d байт данных: %s\n", n, string(buf[:n]))

	n, _, err = conn.ReadFromUnix(buf)
	if err != nil {
		fmt.Println("Ошибка при чтении данных:", err)
		return
	}

	fmt.Printf("Получено %d байт данных: %s\n", n, string(buf[:n]))

	// firstCmd := exec.Command("./bin/t1", file1)
	// secondCmd := exec.Command("./bin/t1", file2)
	// firstBytes := getBytes(firstCmd)
	// secondBytes := getBytes(secondCmd)

	// ioutil.WriteFile("res.txt", xorBytes(firstBytes, secondBytes), fs.FileMode(0777))
}

func getBytes(command *exec.Cmd) []byte {
	reader, writer := io.Pipe()
	command.Stdout = writer
	command.Start()
	data := make([]byte, 1024)
	n, _ := reader.Read(data)
	return data[:n]
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
