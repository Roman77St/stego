// Получает сообщение ввиде потока бит и кодирует в изображение
// Извлекает из изображения поток бит
package stego

import (
	"image"
	"image/color"
)

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

func Decode(img image.Image) []uint8 {
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
                    return collectedBits[32:]
                }
            }
        }
    }
    return collectedBits[32:]
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

// Вспомогательная функция для превращения 32 бит в число uint32
func bitsToUint32(bits []uint8) uint32 {
    var res uint32
    for i := range 32 {
        res = (res << 1) | uint32(bits[i])
    }
    return res
}
