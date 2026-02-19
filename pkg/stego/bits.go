// файл отвечает за низкоуровневую трансформацию данных: преобразование байтов сообщения в поток бит
package stego

const (
	HeaderMagicSize  = 16 // 16 бит для флага "есть сообщение"
	HeaderLengthSize = 32 // 32 бита для длины
	HeaderTotalSize  = HeaderMagicSize + HeaderLengthSize
	MagicValue       = 0x4454 // "DT" в HEX
)

type stegoHeader struct {
	Magic  uint16
	Length uint32
}

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

// PrepareData формирует итоговый массив бит для записи в изображение.
// Собирает пакет: [Magic (16bit)] + [Length (32bit)] + [Payload]
func prepareData(message []byte) []uint8 {
	length := uint32(len(message))

	// Кодируем Magic Value (16 бит)
	magicBits := uintToBits(uint64(MagicValue), HeaderMagicSize)

	// Кодируем длину сообщения в 32 бита.
	lengthBits := uintToBits(uint64(length), HeaderLengthSize)

	// Кодируем само содержание сообщения в биты.
	messageBits := messageToBits(message)

	// Соединяем их: заголовок всегда идет первым.
	// Оптимизация: заранее выделяем точный объем памяти
    finalBits := make([]uint8, 0, HeaderTotalSize + len(messageBits))
    finalBits = append(finalBits, magicBits...)
    finalBits = append(finalBits, lengthBits...)
    finalBits = append(finalBits, messageBits...)

	return finalBits
}

// Универсальная функция для превращения бит в число любого размера
func bitsToUint(bits []uint8) uint64 {
	var res uint64
	for _, b := range bits {
		res = (res << 1) | uint64(b)
	}
	return res
}

// Функция для превращения любого uint в биты заданного размера
func uintToBits(val uint64, size int) []uint8 {
	bits := make([]uint8, size)
	for i := range size {
		bits[size-1-i] = uint8((val >> i) & 1)
	}
	return bits
}
