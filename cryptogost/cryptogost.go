package cryptogost

import (
	"bytes"
	"fmt"
	"os"

	"go.cypherpunks.ru/gogost/v5/gost34112012512"
	"go.cypherpunks.ru/gogost/v5/gost3412128"
	"go.cypherpunks.ru/gogost/v5/gost341264"
	"go.cypherpunks.ru/gogost/v5/gost3413"
)

/*
	func Test() {
		data := []byte("data to be signed")

		password := []byte("Супертест")
		hasher := gost34112012512.New()

		key, err := bcrypt.GenerateFromPassword(password, 10)
		fmt.Print(len(key))
		if err != nil {
			panic("signature is invalid")
		}

		cryptoKuz := gost3412128.NewCipher(key[:32])
		f := gost3413.Pad1(data, cryptoKuz.BlockSize())
		dst := make([]byte, cryptoKuz.BlockSize())

		bufdst := make([]byte, gost3412128.BlockSize)
		cryptoKuz.Encrypt(dst, f[0:16])
		for n := 16; n < len(f); n += 16 {
			cryptoKuz.Encrypt(bufdst, f[n:n+16])
			dst = append(dst, bufdst...)
		}
		decodet := make([]byte, cryptoKuz.BlockSize())
		cryptoKuz.Decrypt(decodet, dst[0:16])
		for n := 16; n < len(dst); n += 16 {
			cryptoKuz.Decrypt(bufdst, dst[n:n+16])
			decodet = append(decodet, bufdst...)
		}
		_, err = hasher.Write(data)
		if err != nil {
			panic("signature is invalid")
		}
		//dgst := hasher.Sum(nil)
		//fmt.Print("\n", string(dgst))
		//	fmt.Print("\n", string(data))
		//	fmt.Print("\n", string(dst))
		index := bytes.IndexByte(decodet, 0)
		if index != -1 {
			trimmedData := decodet[:index]
			fmt.Println(trimmedData)
		} else {
			trimmedData := decodet
			fmt.Println(trimmedData)
		}
		fmt.Print("\n", (decodet))
	}
*/
func genKey(password []byte) []byte {
	hasher := gost34112012512.New()
	_, err := hasher.Write(password)
	if err != nil {
		panic("signature is invalid")
	}
	return hasher.Sum(nil)
}
func EncryptK(path string, password string, path_out string) {
	key := genKey([]byte(password))
	cryptoKuz := gost3412128.NewCipher(key[:32])
	blockSize := cryptoKuz.BlockSize()
	file_input, err := os.Open(path)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file_input.Close()
	fileInfo, err := file_input.Stat()
	if err != nil {
		fmt.Println("Ошибка при получение размера:", err)
		return
	}
	bufdata := make([]byte, fileInfo.Size())

	nData, err := file_input.Read(bufdata)
	if err != nil {
		fmt.Println("Ошибка при чтение:", err)
		return
	}

	f := gost3413.Pad1(bufdata[:nData], blockSize)
	dst := make([]byte, blockSize)
	bufdst := make([]byte, blockSize)
	cryptoKuz.Encrypt(dst, f[0:blockSize])
	for n := blockSize; n < len(f); n += blockSize {
		cryptoKuz.Encrypt(bufdst, f[n:n+blockSize])
		dst = append(dst, bufdst...)
	}
	file, err := os.Create(path_out)
	if err != nil {
		fmt.Println("Ошибка при создание файла:", err)
		return
	}
	_, err = file.Write(dst)
	if err != nil {
		fmt.Println("Ошибка при записи шифрованных данных:", err)
		return
	}
	file.Close()
}
func DecryptK(path string, password string, path_out string) {
	key := genKey([]byte(password))
	cryptoKuz := gost3412128.NewCipher(key[:32])
	blockSize := cryptoKuz.BlockSize()
	file_input, err := os.Open(path)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file_input.Close()
	fileInfo, err := file_input.Stat()
	if err != nil {
		fmt.Println("Ошибка при получение размера:", err)
		return
	}
	bufdata := make([]byte, fileInfo.Size())

	nData, err := file_input.Read(bufdata)
	if err != nil {
		fmt.Println("Ошибка при чтение:", err)
		return
	}
	dst := bufdata[:nData]
	bufdst := make([]byte, blockSize)
	decodet := make([]byte, blockSize)
	cryptoKuz.Decrypt(decodet, dst[0:blockSize])
	for n := blockSize; n < len(dst); n += blockSize {
		cryptoKuz.Decrypt(bufdst, dst[n:n+blockSize])
		decodet = append(decodet, bufdst...)
	}
	var trimmedData []byte
	index := bytes.IndexByte(decodet, 0)
	if index != -1 {
		trimmedData = decodet[:index]

	} else {
		trimmedData = decodet
	}
	file, err := os.Create(path_out)
	if err != nil {
		fmt.Println("Ошибка при создание файла:", err)
		return
	}
	defer file.Close()
	_, err = file.Write(trimmedData)
	if err != nil {
		fmt.Println("Ошибка при записи шифрованных данных:", err)
		return
	}
}
func EncryptM(path string, password string, path_out string) {
	key := genKey([]byte(password))
	cryptoM := gost341264.NewCipher(key[:32])
	blockSize := cryptoM.BlockSize()
	file_input, err := os.Open(path)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file_input.Close()
	fileInfo, err := file_input.Stat()
	if err != nil {
		fmt.Println("Ошибка при получение размера:", err)
		return
	}
	bufdata := make([]byte, fileInfo.Size())

	nData, err := file_input.Read(bufdata)
	if err != nil {
		fmt.Println("Ошибка при чтение:", err)
		return
	}

	f := gost3413.Pad1(bufdata[:nData], blockSize)
	dst := make([]byte, blockSize)
	bufdst := make([]byte, blockSize)
	cryptoM.Encrypt(dst, f[0:blockSize])
	for n := blockSize; n < len(f); n += blockSize {
		cryptoM.Encrypt(bufdst, f[n:n+blockSize])
		dst = append(dst, bufdst...)
	}
	file, err := os.Create(path_out)
	if err != nil {
		fmt.Println("Ошибка при создание файла:", err)
		return
	}
	_, err = file.Write(dst)
	if err != nil {
		fmt.Println("Ошибка при записи шифрованных данных:", err)
		return
	}
	file.Close()
}
func DecryptM(path string, password string, path_out string) {
	key := genKey([]byte(password))
	cryptoM := gost341264.NewCipher(key[:32])
	blockSize := cryptoM.BlockSize()
	file_input, err := os.Open(path)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file_input.Close()
	fileInfo, err := file_input.Stat()
	if err != nil {
		fmt.Println("Ошибка при получение размера:", err)
		return
	}
	bufdata := make([]byte, fileInfo.Size())

	nData, err := file_input.Read(bufdata)
	if err != nil {
		fmt.Println("Ошибка при чтение:", err)
		return
	}
	dst := bufdata[:nData]
	bufdst := make([]byte, blockSize)
	decodet := make([]byte, blockSize)
	cryptoM.Decrypt(decodet, dst[0:blockSize])
	for n := blockSize; n < len(dst); n += blockSize {
		cryptoM.Decrypt(bufdst, dst[n:n+blockSize])
		decodet = append(decodet, bufdst...)
	}
	var trimmedData []byte
	index := bytes.IndexByte(decodet, 0)
	if index != -1 {
		trimmedData = decodet[:index]

	} else {
		trimmedData = decodet
	}
	file, err := os.Create(path_out)
	if err != nil {
		fmt.Println("Ошибка при создание файла:", err)
		return
	}
	defer file.Close()
	_, err = file.Write(trimmedData)
	if err != nil {
		fmt.Println("Ошибка при записи шифрованных данных:", err)
		return
	}
}
