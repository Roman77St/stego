package main

import (
	// "fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
)

func main() {
	// 1. Загружаем исходное изображение
	inputFile := "output.png"
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Не удалось открыть файл %s: %v", inputFile, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Ошибка декодирования изображения: %v", err)
	}

	// // 2. Подготавливаем сообщение (длина + текст)
	// message := "Hello, Go Steganography!"
	// secretBits := PrepareDataForEmbedding(message)

	// // 3. Кодируем сообщение в изображение
	// // Функция Encode возвращает новую картинку *image.RGBA
	// stegoImg := Encode(img, secretBits)

	// // 4. Сохраняем результат в новый файл
	// outputFile := "output.png"
	// out, err := os.Create(outputFile)
	// if err != nil {
	// 	log.Fatalf("Не удалось создать файл %s: %v", outputFile, err)
	// }
	// defer out.Close()

	// err = png.Encode(out, stegoImg)
	// if err != nil {
	// 	log.Fatalf("Ошибка при сохранении PNG: %v", err)
	// }
	// log.Printf("Сообщение успешно спрятано в %s", outputFile)

	// 5. Проверяем результат: пробуем прочитать сообщение из только что созданного файла
	decodedMessage := Decode(img)
	log.Printf("Проверка декодирования: [%s]", decodedMessage)
}

// MessageToBits принимает срез байт и возвращает срез бит (0 или 1)
func MessageToBits(data []byte) []uint8 {
    var bits []uint8
    for _, b := range data {
        // Проходим по каждому биту байта от старшего к младшему
        for i := 7; i >= 0; i-- {
            // Сдвигаем байт и проверяем последний бит
            bit := (b >> i) & 1
            bits = append(bits, bit)
        }
    }
    return bits
}

// BitsToMessage группирует биты по 8 и превращает их в строку
func BitsToMessage(bits []uint8) string {
    var result []byte

    // Проходим по срезу бит с шагом 8
    for i := 0; i < len(bits); i += 8 {
        // Проверка: если бит меньше 8 (остаток), не обрабатываем
        if i+8 > len(bits) {
            break
        }

        var currentByte uint8
        for j := range 8 {
            // Сдвигаем накопленное значение влево на 1
            // И добавляем текущий бит в самый конец
            currentByte = (currentByte << 1) | bits[i+j]
        }
        result = append(result, currentByte)
    }

    return string(result)
}

func LengthToBits(length uint32) []uint8 {
    // Создаем срез из 4 байт
    buf := make([]byte, 4)
    // Записываем число в байты (Big Endian - от старшего к младшему)
    buf[0] = uint8(length >> 24)
    buf[1] = uint8(length >> 16)
    buf[2] = uint8(length >> 8)
    buf[3] = uint8(length)

    // Теперь превращаем эти 4 байта в 32 бита
    return MessageToBits(buf)
}

func PrepareDataForEmbedding(message string) []uint8 {
    data := []byte(message)
    length := uint32(len(data))

    // 1. Получаем биты длины (32 бита)
    lengthBits := LengthToBits(length)

    // 2. Получаем биты сообщения
    messageBits := MessageToBits(data)

    // 3. Склеиваем их
    return append(lengthBits, messageBits...)
}

func setLSB(colorChannel uint8, bit uint8) uint8 {
    // 0xFE это 11111110. Операция & обнуляет последний бит.
    // Затем | (OR) устанавливает нужный нам бит (0 или 1).
    return (colorChannel & 0xFE) | bit
}

func getLSB(color uint8) uint8 {
	// 1 в двоичном виде это 00000001
    // Если последний бит канала 1: (XXXXXXX1 & 00000001) = 1
    // Если последний бит канала 0: (XXXXXXX0 & 00000001) = 0
    return color & 1
}


func Encode(img image.Image, messageBits []uint8) *image.RGBA {
    bounds := img.Bounds()
    newImg := image.NewRGBA(bounds)

    bitIndex := 0
    totalBits := len(messageBits)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            // Получаем оригинальный цвет
            r, g, b, a := img.At(x, y).RGBA()
            r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)

            // Внедряем в Красный канал
            if bitIndex < totalBits {
                r8 = setLSB(r8, messageBits[bitIndex])
                bitIndex++
            }
            // Внедряем в Зеленый канал
            if bitIndex < totalBits {
                g8 = setLSB(g8, messageBits[bitIndex])
                bitIndex++
            }
            // Внедряем в Синий канал
            if bitIndex < totalBits {
                b8 = setLSB(b8, messageBits[bitIndex])
                bitIndex++
            }

            // Устанавливаем новый пиксель (Альфа-канал не трогаем!)
            newImg.Set(x, y, color.RGBA{r8, g8, b8, a8})
        }
    }
    return newImg
}

func Decode(img image.Image) string {
    bounds := img.Bounds()
    var collectedBits []uint8

    // Переменные для контроля длины
    var messageLength uint32
    hasLength := false
    totalBitsNeeded := 32 // Сначала ищем только заголовок

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            channels := []uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}

            for _, ch := range channels {
                if len(collectedBits) < totalBitsNeeded {
                    collectedBits = append(collectedBits, getLSB(ch))

                    // Как только получили первые 32 бита — вычисляем общую длину
                    if !hasLength && len(collectedBits) == 32 {
                        messageLength = bitsToUint32(collectedBits)
                        totalBitsNeeded = 32 + int(messageLength*8)
                        hasLength = true
                    }
                } else {
                    // Все нужные биты собраны
                    return BitsToMessage(collectedBits[32:])
                }
            }
        }
    }
    return BitsToMessage(collectedBits[32:])
}

// Вспомогательная функция для превращения 32 бит в число uint32
func bitsToUint32(bits []uint8) uint32 {
    var res uint32
    for i := range 32 {
        res = (res << 1) | uint32(bits[i])
    }
    return res
}