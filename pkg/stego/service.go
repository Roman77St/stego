package stego

import "image"

func ExtractMessage(img image.Image) []byte {
	data := decode(img)
	message := bitsToMessage(data)
	return message
}

func HideMessage(message []byte, origImg image.Image) image.Image {
	// Подготавливаем сообщение с заголовком
	secretBits := prepareData(message)
	// Кодируем сообщение в изображение
	stegoImg := encode(origImg, secretBits)
	// Сохраняем результат в новый файл
	return stegoImg
}