package cli

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/Roman77St/stego/pkg/stego"
)

// RunEncode выполняет полный цикл скрытия сообщения
func RunEncode(inputPath, outputPath, message string) error {
    img, err := LoadImg(inputPath)
    if err != nil {
        return err
    }

    stegoImg, err := stego.HideMessage([]byte(message), img)
    if err != nil {
        return err
    }

    return SaveImg(outputPath, stegoImg)
}

// RunDecode выполняет извлечение сообщения
func RunDecode(inputPath string) (string, error) {
    img, err := LoadImg(inputPath)
    if err != nil {
        return "", err
    }

    res, err := stego.ExtractMessage(img)
    if err != nil {
        return "", err
    }

    return string(res), nil
}

func LoadImg(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %s: %v", path, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)

	return img, err
}

func SaveImg(path string, img image.Image) error {
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("не удалось создать файл %s: %v", path, err)
	}
	defer outFile.Close()

	return png.Encode(outFile, img)
}
