package transform

import (
	"core/compression"
	"core/cryptogost"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/klauspost/compress/zstd"
)

func RecoveryFile(path string, filename string, CompressionStat bool, ExpansionFormat string, CryptoName string, CryptoPassword string) {
	inputFile, err := os.Open(path)
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	path_encrypto := "./tmp/server/zt/" + strings.Split(filename, ".")[0] + ".zst"
	path_encoded := "./org/" + strings.Split(filename, ".")[0] + ExpansionFormat
	encryptoFile, err := os.Create(path_encrypto)
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	encryptoFile.Close()
	encodedFile, err := os.Create(path_encoded)
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	defer encodedFile.Close()

	switch CryptoName {
	case "Магма":
		cryptogost.DecryptK(path, CryptoPassword, path_encrypto)
	case "Кузнечик":
		cryptogost.DecryptK(path, CryptoPassword, path_encrypto)
	default:
		encryptoFile, err := os.OpenFile(path_encrypto, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			fmt.Println("Ошибка при копирование не шифрованного файла:", err)
			return
		}
		_, err = io.Copy(encryptoFile, inputFile)
		if err != nil {
			fmt.Println("Ошибка при передаче файлаfg:", err)
			return
		}
		encryptoFile.Close()
	}
	inputFile.Close()

	encryptoFile, err = os.OpenFile(path_encrypto, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println("Ошибка при открытие перед кодирование:", err)
		return
	}
	if CompressionStat {
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
		err = compression.DecompressOne(encryptoFile, encodedFile)
		if err != nil {
			log.Fatal("Ошибка:", err)
		}
	} else {
		_, err = io.Copy(encodedFile, encryptoFile)
		if err != nil {
			fmt.Println("Ошибка при передаче файла:", err)
			return
		}
	}

}

func CodeFile(path string, filename string, CompressionLevel int, CryptoName string, CryptoPassword string) (string, string) {
	inputFile, err := os.Open(path)
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	defer inputFile.Close()
	path_encoded := "./tmp/client/zt/" + strings.Split(filename, ".")[0] + ".zst"
	path_encrypto := "./tmp/client/crypto/" + strings.Split(filename, ".")[0] + ".gost"
	encryptoFile, err := os.Create(path_encrypto)
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	encryptoFile.Close()
	encodedFile, err := os.Create(path_encoded)
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	defer encodedFile.Close()
	if CompressionLevel == -1 {
		_, err = io.Copy(encodedFile, inputFile)
		if err != nil {
			fmt.Println("Ошибка при передаче файла:", err)
			return "", ""
		}
	} else {
		err = compression.CompressOne(inputFile, encodedFile, zstd.EncoderLevelFromZstd(CompressionLevel))
		if err != nil {
			log.Fatal("Ошибка:", err)
		}

	}
	inputFile.Close()
	encodedFile.Close()
	path_itog := path_encrypto
	name_itog := strings.Split(filename, ".")[0] + ".gost"
	switch CryptoName {
	case "Магма":

		cryptogost.EncryptM(path_encoded, CryptoPassword, path_encrypto)
	case "Кузнечик":

		cryptogost.EncryptK(path_encoded, CryptoPassword, path_encrypto)
	default:
		path_itog = "./tmp/client/crypto/" + strings.Split(filename, ".")[0] + ".zst"
		name_itog = strings.Split(filename, ".")[0] + ".zst"
		CopItog(path_itog, path_encoded)
	}
	return path_itog, name_itog
}
func CopItog(path_itog, path_encoded string) {
	enFile, err := os.Create(path_itog)
	if err != nil {
		fmt.Println("Ошибка при копирование не шифрованного файла:", err)
		return
	}
	defer enFile.Close()
	encFile, err := os.Open(path_encoded)
	if err != nil {
		fmt.Println("Ошибка при копирование не шифрованного файла:", err)
		return
	}
	defer encFile.Close()
	nkl, err := io.Copy(enFile, encFile)
	if err != nil {
		fmt.Println("Ошибка при передаче не шифрованного файла:", err)
		return
	}
	fmt.Print(nkl)
	encFile.Close()
	enFile.Close()
}
