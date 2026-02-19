package cli

import (
	"flag"
	"fmt"
	"os"
)

func RunCLI() error {
	mode := flag.String("mode", "decode", "режим работы: encode или decode")
	input := flag.String("input", "input.png", "путь к входному изображению")
	output := flag.String("output", "output.png", "путь для сохранения (только для encode)")
	msg := flag.String("msg", "", "сообщение для записи (только для encode)")
	// help := flag.String("help", "", "использование утилиты")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Используйте флаги. Справка: -h")
		return nil
	}

	switch *mode {
	case "encode":
		if err := RunEncode(*input, *output, *msg); err != nil {
			return err
		}
		fmt.Println("✅ Готово!")
	case "decode":
		res, err := RunDecode(*input)
		if err != nil {
			return err
		}
		fmt.Printf("🔓 Сообщение: %s\n", res)
	}
	return nil
}
