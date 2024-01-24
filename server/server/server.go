package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/quic-go/quic-go"
)

func ServerReceiving() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("Ошибка при загрузке ключей:", err)
		return
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	listener, err := quic.ListenAddr("0.0.0.0:4242", tlsConfig, &quic.Config{})
	if err != nil {
		fmt.Println("Ошибка при создании слушателя:", err)
		return
	}
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			fmt.Println("Ошибка при приеме сессии:", err)
			continue
		}
		go handleQUICSession(sess)
	}
}
func handleQUICSession(sess quic.Connection) {
	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		fmt.Println("Ошибка при открытии потока:", err)
		return
	}
	fmt.Println("Получили новое соединение:", stream.StreamID())
	nameSize := make([]byte, 2)
	_, err = stream.Read(nameSize)
	if err != nil {
		fmt.Println("Ошибка при чтении размера имени файла:", err)
		return
	}
	if err = SyncWriteCode(stream); err != nil {
		fmt.Println("Ошибка при синхронизации размера имени файла:", err)
	}
	nameFileSize := binary.BigEndian.Uint16(nameSize)
	headerBuf := make([]byte, nameFileSize+uint16(8))
	_, err = stream.Read(headerBuf)
	if err != nil && err != io.EOF {
		fmt.Println("Ошибка при чтении имени файла:", err)
		return
	}
	headerDelimiter := bytes.IndexByte(headerBuf, 0)
	if headerDelimiter == -1 {
		fmt.Println("Ошибка: неверный формат заголовка")
		return
	}
	if err = SyncWriteCode(stream); err != nil {
		fmt.Println("Ошибка при синхронизации размера имени файла:", err)
	}
	filename := string(headerBuf[:headerDelimiter])
	fmt.Println(filename)
	fileSize, _ := binary.Uvarint(headerBuf[headerDelimiter+1 : headerDelimiter+9])
	fmt.Println(fileSize)

	file, err := os.Create("./tmp/server/" + filename)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	bufSize := min(fileSize, min(268435456, max(8388608, m.Sys-m.Alloc/4*8)))
	if err != nil {
		panic(err)
	}
	fileBuf := make([]byte, bufSize)
	bytesReceived := uint64(0)
	for bytesReceived < fileSize {
		n, err := stream.Read(fileBuf)
		if err != nil && err != io.EOF {
			fmt.Println("Ошибка при получении файла:", err)
			return
		}
		_, err = file.Write(fileBuf[:n])
		if err != nil {
			fmt.Println("Ошибка при записи в файл:", err)
			return
		}
		//file.Sync()
		bytesReceived += uint64(n)
	}

	fmt.Println("Файл успешно получен и сохранен:", filename)
	file.Close()

}
func SyncWriteCode(stream quic.Stream) error {
	ok := []byte{200}
	_, err := stream.Write(ok)
	return err
}
