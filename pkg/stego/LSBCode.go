// Получает сообщение ввиде потока бит и кодирует в изображение
// Извлекает из изображения поток бит
package stego

import (
	"fmt"
	"image"
	"image/color"
)

// Encode записывает поток бит в младшие биты (LSB) цветовых каналов изображения.
// Возвращает новое изображение в формате RGBA с внедренными данными.
func encode(img image.Image, messageBits []uint8) *image.RGBA {
    bounds := img.Bounds()
    // Создаем новое пустое изображение того же размера
    newImg := image.NewRGBA(bounds)

    bitIndex := 0
    totalBits := len(messageBits)
    // Проходим по картинке построчно (Y), затем попиксельно (X)
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            // Получаем цвета текущего пикселя (значения 0-65535)
            r, g, b, a := img.At(x, y).RGBA()
            // Приводим к стандарту 8-бит (0-255)
            r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)

            // Если биты сообщения еще не закончились, внедряем их в R, G и B каналы
            if bitIndex < totalBits {
                r8 = setLSB(r8, messageBits[bitIndex])
                bitIndex++
            }
            if bitIndex < totalBits {
                g8 = setLSB(g8, messageBits[bitIndex])
                bitIndex++
            }
            if bitIndex < totalBits {
                b8 = setLSB(b8, messageBits[bitIndex])
                bitIndex++
            }

            // Записываем измененный пиксель в новое изображение
            // Альфа-канал (прозрачность) оставляем без изменений
            newImg.Set(x, y, color.RGBA{r8, g8, b8, a8})
        }
    }
    return newImg
}

// Decode извлекает скрытые биты из изображения.
// Сначала читается заголовок для определения наличия и длины сообщения,
// затем извлекается само тело сообщения.
func decode(img image.Image) ([]uint8, error) {
    bounds := img.Bounds()
    var collectedBits []uint8

    var header stegoHeader
    hasHeader := false
    totalBitsNeeded := HeaderTotalSize // Начальная цель — прочитать заголовк

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            // Оптимизация: работаем с каналами напрямую без создания слайса
            channels := [3]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}

            for _, ch := range channels {
                if len(collectedBits) < totalBitsNeeded {
                    // Собираем LSB из текущего канала
                    collectedBits = append(collectedBits, getLSB(ch))

                    // Как только получили заголовок — проверяем наличие сообщения и вычисляем общую длину
                    if !hasHeader && len(collectedBits) == HeaderTotalSize {
                        header.Magic = uint16(bitsToUint(collectedBits[:HeaderMagicSize]))
						header.Length = uint32(bitsToUint(collectedBits[HeaderMagicSize:HeaderTotalSize]))
                        // 1. Проверка флага: если магическое число не совпало, это не наше сообщение
						if header.Magic != MagicValue {
							return nil, fmt.Errorf("steganography header not found (invalid magic)")
						}
                        // 2. Валидация длины по емкости картинки
						if int(header.Length) > maxCapacity(img) {
							return nil, fmt.Errorf("invalid message length in header")
						}
                        // Рассчитываем итоговое количество бит: заголовок + тело
                        totalBitsNeeded = HeaderTotalSize + int(header.Length*8)
                        hasHeader = true
                    }
                } else {
                    // Если собрали всё необходимое, возвращаем биты без заголовка
                    return collectedBits[HeaderTotalSize:], nil
                }
            }
        }
    }
    return nil, fmt.Errorf("unexpected end of image data")
}

// setLSB заменяет последний бит в байте цветового канала на нужный бит сообщения.
func setLSB(colorChannel uint8, bit uint8) uint8 {
    // 0xFE это 11111110. Операция & обнуляет последний бит.
    // Затем | (OR) устанавливает нужный нам бит (0 или 1).
    return (colorChannel & 0xFE) | bit
}

// getLSB возвращает значение последнего бита в байте.
func getLSB(color uint8) uint8 {
	// 1 в двоичном виде это 00000001
    // Если последний бит канала 1: (XXXXXXX1 & 00000001) = 1
    // Если последний бит канала 0: (XXXXXXX0 & 00000001) = 0
    return color & 1
}

// MaxBytesCapacity возвращает максимальное количество байт,
// которое можно скрыть в изображении (за вычетом заголовка).
func maxCapacity(img image.Image) int {
    bounds := img.Bounds()
    // Всего доступно бит: W * H * 3 канала
    totalBits := bounds.Dx() * bounds.Dy() * 3

    // Вычитаем длину заголовока и переводим в байты
    if totalBits < HeaderTotalSize {
        return 0
    }
    return (totalBits - HeaderTotalSize) / 8
}