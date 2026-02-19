package stego

import (
	"fmt"
	"image"
)

func ExtractMessage(img image.Image) ([]byte, error) {
	data, err := decode(img)
	if err != nil {
		return nil, err
	}
	message := bitsToMessage(data)
	return message, nil
}

func HideMessage(message []byte, origImg image.Image) (image.Image, error) {
	// Проверяем размер изображения
	capacity := maxCapacity(origImg)
	if len(message) > capacity {
		return nil, fmt.Errorf("сообщение слишком большое (%d байт), макс. вместимость %d байт", len(message), capacity)
	}
	if capacity <= 0 {
		return nil, fmt.Errorf("изображение слишком маленькое для записи данных")
	}

	// Подготавливаем сообщение с заголовком
	secretBits := prepareData(message)

	// Кодируем сообщение в изображение
	stegoImg := encode(origImg, secretBits)

	// Сохраняем результат в новый файл
	return stegoImg, nil
}

// GetMaxCapacity возвращает вместимость изображения в байтах для пользователя.
func GetMaxCapacity(img image.Image) int {
	return maxCapacity(img)
}
