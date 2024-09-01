package utils

import (
	"log"
	"os"

	"github.com/disintegration/imaging"
)

func CreateFolderIfNotExists(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return os.MkdirAll(folderPath, os.ModePerm)
	}
	return nil
}

func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		return os.Remove(filePath)
	}
	return nil
}

func CompressAndResizeImage(inputPath, outputPath string) error {
	// Открываем изображение
	src, err := imaging.Open(inputPath)
	if err != nil {
		return err
	}

	// Изменяем размер изображения
	dst := imaging.Resize(src, 225, 300, imaging.Lanczos)

	// Сохраняем обработанное изображение
	err = imaging.Save(dst, outputPath)
	if err != nil {
		return err
	}

	log.Printf("Изображение успешно сжато и обрезано: %s", outputPath)
	return nil
}
