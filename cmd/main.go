package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/Roman77St/stego/pkg/stego"
)

func main() {
	mode := flag.String("mode", "decode", "режим работы: encode или decode")
	inputPath := flag.String("input", "input.png", "путь к входному изображению")
	outputPath := flag.String("output", "output.png", "путь для сохранения (только для encode)")
	msg := flag.String("msg", "", "сообщение для записи (только для encode)")
	// help := flag.String("help", "", "использование утилиты")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Для использования утилиты необходимо применять флаги")
		os.Exit(0)
	}

	switch *mode {
		case "encode":
			if *msg == "" {
				log.Fatal("Ошибка: необходимо указать сообщение через флаг -msg")
			}
			img, err := loadImg(*inputPath)
			if err != nil {
				log.Fatal(err)
			}
			stegoImg := stego.HideMessage([]byte(*msg), img)
			err = saveImg(*outputPath, stegoImg)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("✅ Сообщение успешно зашито в %s\n", *outputPath)

		case "decode":
			img, err := loadImg(*inputPath)
			if err != nil {
				log.Fatal(err)
			}
			res := stego.ExtractMessage(img)
			fmt.Printf("🔓 Извлеченное сообщение: %s\n", string(res))

		default:
			fmt.Println("Использование:")
			flag.PrintDefaults()
		}

}


func loadImg(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %s: %v", path, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования изображения: %v", err)
	}

	return img, nil
}

func saveImg(path string, img image.Image) error {
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("не удалось создать файл %s: %v", path, err)
	}
	defer outFile.Close()
	err = png.Encode(outFile, img)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении PNG: %v", err)
	}

	return nil
}
