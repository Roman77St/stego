// транслирует байты в поток бит и обратно
// подготавливает заголовок

package stego

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

// BitsToMessage группирует биты по 8 и превращает их в срез байт
func BitsToMessage(bits []uint8) []byte {
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

    return result
}

func PrepareData(message []byte) []uint8 {
    length := uint32(len(message))

    // 1. Получаем биты длины (32 бита)
    lengthBits := LengthToBits(length)

    // 2. Получаем биты сообщения
    messageBits := MessageToBits(message)

    // 3. Склеиваем их
    return append(lengthBits, messageBits...)
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