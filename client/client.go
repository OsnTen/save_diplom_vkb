package client

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bartmeuris/progressio"

	"github.com/quic-go/quic-go"
)

func copyProgress(w io.Writer, r io.Reader, size int64) (written int64, err error) {

	// Wrap your io.Writer:
	pw, ch := progressio.NewProgressWriter(w, size)
	defer pw.Close()

	// Launch a Go-Routine reading from the progress channel
	go func() {
		for p := range ch {
			fmt.Printf("\rProgress: %s", p.String())
		}
		fmt.Printf("\nDone\n")
	}()

	// Copy the data from the reader to the new writer
	return io.Copy(pw, r)
}
func Client(path, file_name string) {
	// Создайте собственный сертификат и ключ TLS (или используйте существующие)
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("Ошибка при загрузке сертификата и ключа:", err)
		return
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ServerName:         "example.com", // измените на домен сервера
		InsecureSkipVerify: true,          // ВНИМАНИЕ: не используйте в продакшн!
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5s крайний срок выполнения
	defer cancel()

	addr := "127.0.0.1:4242" // адрес сервера
	session, err := quic.DialAddr(ctx, addr, tlsConfig, nil)
	// Устанавливаем соединение с сервером на порту 4242
	if err != nil {
		fmt.Println("Ошибка при установке соединения:", err)
		return
	}

	//file_name := "test.zst"

	// Открываем файл на чтение
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Получаем первый поток сессии QUIC
	stream, err := session.OpenStreamSync(ctx)
	if err != nil {
		fmt.Println("Ошибка при открытии потока:", err)
		return
	}

	byteFilename := []byte(file_name + "\x00")
	FileNameSize := make([]byte, 2)

	binary.BigEndian.PutUint16(FileNameSize, uint16(len(byteFilename)))

	_, err = stream.Write(FileNameSize)
	if err != nil {
		fmt.Println("Ошибка при отправке имени файла:", err)
		return
	}
	if SyncReadCode(stream) {
		fmt.Println("Ошибка при синхронизации размера имени файла:", err)
	}

	// Отправляем имя файла и размер
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Ошибка при получении информации о файле:", err)
		return
	}
	bs := make([]byte, 8)
	binary.PutUvarint(bs, uint64(fileInfo.Size()))
	bs = append([]byte(file_name+"\x00"), bs...)
	_, err = stream.Write(bs)
	if err != nil {
		fmt.Println("Ошибка при отправке размера файла:", err)
		return
	}
	if SyncReadCode(stream) {
		fmt.Println("Ошибка при синхронизации размера имени файла:", err)
	}

	_, err = copyProgress(stream, file, fileInfo.Size())
	if err != nil {
		fmt.Println("Ошибка при передаче файла:", err)
		return
	}
	stream.Close()
	fmt.Println("Файл успешно отправлен")
}

func SyncReadCode(stream quic.Stream) bool {
	buf := make([]byte, 1)
	_, err := stream.Read(buf)
	if err != nil || uint8(buf[0]) != uint8(200) {
		fmt.Println("Ошибка при отправке размера файла:", err)
		return true
	}
	return false
}
