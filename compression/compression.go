package compression

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

func DecompressOne(in io.Reader, out io.Writer) error {
	dec, err := zstd.NewReader(in)
	if err != nil {
		return err
	}
	defer dec.Close()
	_, err = io.Copy(out, dec)
	if err != nil {
		return err
	}
	return nil
}

func CompressOne(in io.Reader, out io.Writer, level zstd.EncoderLevel) error {
	enc, err := zstd.NewWriter(out, zstd.WithEncoderLevel(level))
	if err != nil {
		return err
	}
	_, err = io.Copy(enc, in)
	if err != nil {
		enc.Close()
		return err
	}
	return enc.Close()
}
func CompressTar(name string, files []string, level zstd.EncoderLevel) {
	outFile, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Создание zstd.Writer
	zstdWriter, err := zstd.NewWriter(outFile, zstd.WithEncoderLevel(1))
	if err != nil {
		return
	}
	// Создание tar.Writer для записи файлов в архив
	tarWriter := tar.NewWriter(zstdWriter)
	defer tarWriter.Close()

	for _, file := range files {
		err = addFileToArchive(tarWriter, file)
		if err != nil {
			panic(err)
		}
	}

	// Закрытие zstd.Writer, чтобы завершить запись архива.
	// Это также закроет tar.Writer
	err = zstdWriter.Close()
	if err != nil {
		panic(err)
	}
}
func addFileToArchive(tarWriter *tar.Writer, filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return nil
	}

	sourceFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	tarHeader, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}

	tarHeader.Name = filepath.Base(filePath)
	err = tarWriter.WriteHeader(tarHeader)
	if err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
