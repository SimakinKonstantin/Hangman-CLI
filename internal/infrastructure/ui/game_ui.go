package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/constants"
)

type GameUI struct {

	// Текущее слова для вывода, включая разгаданные, неразгаданные символы.
	wordProgress string
	reader       *bufio.Reader
	writer       *bufio.Writer
}

// Создает новый экземпляр GameUI.
func NewGameUI(rightWord string, reader io.Reader, writer io.Writer) GameUI {
	size := utf8.RuneCountInString(rightWord)
	return GameUI{wordProgress: strings.Repeat("_", size), reader: bufio.NewReader(reader), writer: bufio.NewWriter(writer)}
}

// Возврашает текущий игровой прогресс.
func (ui GameUI) GetProgress() string {
	return ui.wordProgress
}

// Выводит сообщение о проигрыше.
func (ui GameUI) PrintLose() error {
	err := ui.ClearScreen()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(ui.writer, "Вы проиграли :(")
	if err != nil {
		return err
	}

	err = ui.writer.Flush()

	return err
}

// Выводит сообщение о выигрыше.
func (ui GameUI) PrintWin() error {
	err := ui.ClearScreen()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(ui.writer, "Вы победили !")
	if err != nil {
		return err
	}

	err = ui.writer.Flush()

	return err
}

// Обновляет прогресс после хода игрока.
// Возвращает кол-во букв, которые осталось угадать.
// indexes - массив индексов букв, которые правильно угадал пользователь в последнем шаге.
func (ui *GameUI) UpdateProgress(indexes []int, rightWord []rune) int {
	newProgress := []rune(ui.wordProgress)
	for _, index := range indexes {
		newProgress[index] = rightWord[index]
	}

	ui.wordProgress = string(newProgress)

	lettersLeft := 0

	for _, char := range newProgress {
		if char == '_' {
			lettersLeft++
		}
	}

	return lettersLeft
}

// Печатает блок информации во время игры.
func (ui GameUI) Print(prefix string, value any) error {
	_, err := fmt.Fprint(ui.writer, prefix, value, "\n")
	if err != nil {
		return err
	}

	return ui.writer.Flush()
}

// Принимает ввод значения для игрока.
func (ui GameUI) InputLetter() (rune, error) {
	for {
		_, err := fmt.Fprint(ui.writer, "Введите букву или '?' для получения подсказки: ")
		if err != nil {
			return ' ', err
		}

		err = ui.writer.Flush()
		if err != nil {
			return ' ', err
		}

		byteInput, _, err := ui.reader.ReadLine()
		if err != nil {
			return 0, err
		}

		input := string(byteInput)

		// Обработка подсказки.
		if input == "?" {
			err = ui.ClearScreen()
			return '?', err
		}

		// Обработка буквы.
		if utf8.RuneCountInString(input) == 1 {
			runeInput, _ := utf8.DecodeRuneInString(input)
			res := unicode.ToUpper(runeInput)
			err = ui.ClearScreen()

			return res, err
		}
	}
}

// Печатает текущее состояние виселицы.
func (ui GameUI) PrintPicture(errorCount int) error {
	pictures := []string{"_________\n",
		"|\n|\n|\n|\n|\n|\n",
		"|       |\n|       |\n|\n|\n|\n|\n",
		"|       |\n|       |\n|       ()\n|\n|\n|\n",
		"|       |\n|       |\n|       ()\n|       |\n|\n|\n",
		"|       |\n|       |\n|       ()\n|      /|\n|\n|\n",
		"|       |\n|       |\n|       ()\n|      /|\\\n|\n|\n",
		"|       |\n|       |\n|       ()\n|      /|\\\n|      /\n|\n",
		"|       |\n|       |\n|       ()\n|      /|\\\n|      /\\\n|\n",
		"-----\n"}

	switch errorCount {
	case 0:
		_, err := fmt.Fprint(ui.writer, "\n\n\n\n\n\n")
		if err != nil {
			return err
		}

		return ui.writer.Flush()
	case 1:
		_, err := fmt.Fprint(ui.writer, "\n\n\n\n\n\n", pictures[len(pictures)-1])
		if err != nil {
			return err
		}

		return ui.writer.Flush()
	case 2:
		_, err := fmt.Fprint(ui.writer, pictures[1], pictures[len(pictures)-1])
		if err != nil {
			return err
		}

		return ui.writer.Flush()
	}

	if errorCount > 2 && errorCount < constants.MaxErrors+1 {
		// Индекс элемента массива, который будет использоваться как вертикальная часть виселицы.
		verticalIndex := errorCount - 2

		_, err := fmt.Fprint(ui.writer, pictures[0], pictures[verticalIndex], pictures[len(pictures)-1])
		if err != nil {
			return err
		}

		return ui.writer.Flush()
	}

	return nil
}

// Очищает экран с игрой. Для нормальной работы этой функции, программу нужно вызывать в терминале windows, linux.
func (ui GameUI) ClearScreen() error {
	osName := runtime.GOOS

	switch osName {
	case "windows":
		command := exec.Command("cmd", "/c", "cls")
		command.Stdout = os.Stdout

		err := command.Run()
		if err != nil {
			return TerminalClearError{"ошибка очистки терминала", err}
		}
	case "linux":
		command := exec.Command("clear")
		command.Stdout = os.Stdout

		err := command.Run()
		if err != nil {
			return TerminalClearError{"ошибка очистки терминала", err}
		}
	default:
		_, err := fmt.Fprint(ui.writer, "\n\n\n\n\n\n\n\n\n\n\n\n")
		if err != nil {
			return err
		}

		err = ui.writer.Flush()
		if err != nil {
			return err
		}
	}

	return nil
}
