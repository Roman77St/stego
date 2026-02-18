// Получает сообщение ввиде потока бит и кодирует в изображение
// Извлекает из изображения поток бит
package stego

import (
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
// Сначала читаются первые 32 бита для определения длины сообщения,
// затем извлекается само тело сообщения.
func decode(img image.Image) []uint8 {
    bounds := img.Bounds()
    var collectedBits []uint8

    var messageLength uint32
    hasLength := false
    totalBitsNeeded := 32 // Начальная цель — прочитать 32 бита заголовка

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            // Оптимизация: работаем с каналами напрямую без создания слайса
            channels := [3]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}

            for _, ch := range channels {
                if len(collectedBits) < totalBitsNeeded {
                    // Собираем LSB из текущего канала
                    collectedBits = append(collectedBits, getLSB(ch))

                    // Как только получили первые 32 бита — вычисляем общую длину
                    if !hasLength && len(collectedBits) == 32 {
                        messageLength = bitsToUint32(collectedBits)
                        // Рассчитываем итоговое количество бит: заголовок + тело
                        totalBitsNeeded = 32 + int(messageLength*8)
                        hasLength = true
                    }
                } else {
                    // Если собрали всё необходимое, возвращаем биты без заголовка
                    return collectedBits[32:]
                }
            }
        }
    }
    // На случай, если картинка закончилась раньше, чем мы собрали всё сообщение
    if len(collectedBits) > 32 {
        return collectedBits[32:]
    }
    return nil
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

// bitsToUint32 превращает срез из 32 бит обратно в целое число (длину сообщения).
func bitsToUint32(bits []uint8) uint32 {
    var res uint32
    for i := range 32 {
        // Сдвигаем результат влево и добавляем текущий бит
        res = (res << 1) | uint32(bits[i])
    }
    return res
}
