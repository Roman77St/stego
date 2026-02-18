package stego

import "image"

func ExtractMessage(img image.Image) []byte {
	data := Decode(img)
	message := BitsToMessage(data)
	return message
}

func HideMessage(message []byte, origImg image.Image) image.Image {
	// Подготавливаем сообщение с заголовком
	secretBits := PrepareData(message)
	// Кодируем сообщение в изображение
	stegoImg := Encode(origImg, secretBits)
	// Сохраняем результат в новый файл
	return stegoImg
}