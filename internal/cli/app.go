package cli

import (
	"flag"
	"fmt"
	"os"
)

func RunCLI() error {

	helpMessage := `
Stego Tool — утилита для скрытия текста в изображениях (LSB)

Использование:
  stego -mode encode -input <file> -output <file> -msg <text>
  stego -mode decode -input <file>

Флаги:

`
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, helpMessage)
		flag.PrintDefaults()
	}

	mode := flag.String("mode", "", "Режим работы: encode или decode (обязательно)")
	input := flag.String("input", "", "Путь к входному изображению (обязательно)")
	output := flag.String("output", "output.png", "Путь для сохранения результата (только для encode)")
	msg := flag.String("msg", "", "Сообщение для сокрытия (только для encode)")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Используйте флаги. Справка: -h")
		return nil
	}

	switch *mode {
	case "encode":
		if *msg == "" {
			return fmt.Errorf("❌ Ошибка: сообщение (-msg) не может быть пустым для режима encode")
		}
		if err := RunEncode(*input, *output, *msg); err != nil {
			return fmt.Errorf("❌ Ошибка при кодировании: %v", err)
		}
		fmt.Printf("✅ Успех! Сообщение спрятано в: %s\n", *output)

	case "decode":
		res, err := RunDecode(*input)
		if err != nil {
			return fmt.Errorf("❌ Ошибка при декодировании: %v", err)
		}
		fmt.Printf("🔓 Извлечено сообщение:\n---\n%s\n---\n", res)

	default:
		flag.Usage()
		return fmt.Errorf("❌ Неизвестный режим: %s", *mode)
	}
	return nil
}
