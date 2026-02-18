// файл отвечает за низкоуровневую трансформацию данных: преобразование байтов сообщения в поток бит

package stego

// MessageToBits принимает срез байт (сообщение) и разворачивает его в поток бит.
// Каждый байт превращается в 8 элементов среза uint8, где каждый элемент — это 0 или 1.
// Биты извлекаются от старшего к младшему (Most Significant Bit first).
// Если на вход пришла буква 'H' (код 72, в двоичной системе 01001000), цикл for i := 7; i >= 0
// поочередно достанет каждый разряд, и в итоговый срез попадут числа [0, 1, 0, 0, 1, 0, 0, 0]
func messageToBits(data []byte) []uint8 {
	// Предварительно аллоцируем память для ускорения работы: на 1 байт приходится 8 бит.
    bits := make([]uint8, 0, len(data)*8)

    for _, b := range data {
        // Проходим по каждому биту байта от старшего к младшему
        for i := 7; i >= 0; i-- {
            // Операция (b >> i) сдвигает нужный бит в самую правую позицию.
            // Операция & 1 отсекает все остальные биты, оставляя только 0 или 1
            bit := (b >> i) & 1
            bits = append(bits, uint8(bit))
        }
    }
    return bits
}

// BitsToMessage группирует одиночные биты (0 и 1) обратно в байты.
// Каждые 8 бит формируют один символ сообщения.
func bitsToMessage(bits []uint8) []byte {
    var result []byte

    // Проходим по срезу бит с шагом 8 (размер одного байта).
    for i := 0; i < len(bits); i += 8 {
        // Если оставшихся бит меньше 8, значит это неполный байт, игнорируем его.
        if i+8 > len(bits) {
            break
        }

        var currentByte uint8
		// Собираем байт из 8 последовательных бит.
        for j := range 8 {
            // Сдвигаем уже накопленные в currentByte биты влево на одну позицию,
            // освобождая место для нового бита в самом младшем разряде.
            // Затем через операцию ИЛИ (|) добавляем текущий бит
            currentByte = (currentByte << 1) | bits[i+j]
        }
        result = append(result, currentByte)
    }

    return result
}

// LengthToBits преобразует число типа uint32 (длину сообщения) в срез из 32 бит.
// Это позволяет упаковать метаданные о размере данных перед самим сообщением.
func lengthToBits(length uint32) []uint8 {
	// Разбиваем 32-битное число на 4 байта (Big Endian).
    buf := make([]byte, 4)
    // Записываем число в байты (Big Endian - от старшего к младшему)
    buf[0] = uint8(length >> 24) // Старшие 8 бит
    buf[1] = uint8(length >> 16)
    buf[2] = uint8(length >> 8)
    buf[3] = uint8(length)       // Младшие 8 бит

    // Теперь превращаем эти 4 байта в 32 бита
    return messageToBits(buf)
}

// PrepareData формирует итоговый массив бит для записи в изображение.
// Пакет состоит из 32 бит заголовка (длина) и битов самого сообщения.
func prepareData(message []byte) []uint8 {
	length := uint32(len(message))

	// 1. Кодируем длину сообщения в 32 бита (заголовок).
	lengthBits := lengthToBits(length)

	// 2. Кодируем само содержание сообщения в биты.
	messageBits := messageToBits(message)

	// 3. Соединяем их: заголовок всегда идет первым.
	// Оптимизация: заранее выделяем точный объем памяти
    finalBits := make([]uint8, 0, len(lengthBits)+len(messageBits))
    finalBits = append(finalBits, lengthBits...)
    finalBits = append(finalBits, messageBits...)

	return finalBits
}

// TODO Возможно понадобятся другие заголовки.
// Для них можно использовать структуру типа:
//
// type Header struct {
//     Size    uint32
//     Version uint8
//     // Можно добавить еще поля
// }
// 
// // Новая функция-агрегатор
// func (h Header) Serialize() []uint8 {
//     var headerBytes []byte
//     // Упаковываем все поля в байты...
//     return MessageToBits(headerBytes)
// }