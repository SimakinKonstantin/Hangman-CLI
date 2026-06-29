package ui

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/constants"
)

// UI, который используются на этапе инициализации игры, когда выбираются категория слов, слово и уровень сложности.
type InitUI struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewInitUI(reader io.Reader, writer io.Writer) InitUI {
	return InitUI{bufio.NewReader(reader), bufio.NewWriter(writer)}
}

// Запрашивает число - кол-во разрешенных ошибок.
// При некорректном вводе генерирует случайную сложность.
func (ui *InitUI) InputDifficulty() (int, error) {
	_, err := fmt.Fprintf(ui.writer, "Введите уровень сложности (кол-во попыток) от 1 до %d: ", constants.MaxErrors)
	if err != nil {
		return 0, err
	}

	err = ui.writer.Flush()
	if err != nil {
		return 0, err
	}

	byteDifficulty, _, err := ui.reader.ReadLine()
	if err != nil {
		return 0, err
	}

	if len(byteDifficulty) == 0 {
		_, err = fmt.Fprintln(ui.writer, "Указано пустое значение, уровень сложности выбран случайно")
		if err != nil {
			return 0, err
		}

		err = ui.writer.Flush()
		if err != nil {
			return 0, err
		}

		randInt, err := getRandomInt(1, constants.MaxErrors+1)
		if err != nil {
			return 0, err
		}

		return randInt, nil
	}

	difficulty, err := strconv.Atoi(string(byteDifficulty))
	if err != nil {
		return 0, err
	}

	if difficulty < 1 || difficulty > constants.MaxErrors {
		_, err = fmt.Fprintln(ui.writer, "Указано некорректное значение, уровень сложности выбран случайно")
		if err != nil {
			return 0, err
		}

		err = ui.writer.Flush()
		if err != nil {
			return 0, err
		}

		return getRandomInt(1, constants.MaxErrors+1)
	}

	return difficulty, err
}

// Запрашивает игрока ввода названия категории и возвращет его.
// При некорректном вводе выбирает случайную сложность.
func (ui *InitUI) InputCategory(categories map[string][]constants.GameWord) (string, error) {
	_, err := fmt.Fprintln(ui.writer, "Доступные категории слов:")
	if err != nil {
		return "", err
	}

	err = ui.writer.Flush()
	if err != nil {
		return "", err
	}

	for key := range categories {
		_, err = fmt.Fprint(ui.writer, key, " ")
		if err != nil {
			return "", err
		}

		err = ui.writer.Flush()
		if err != nil {
			return "", err
		}
	}

	_, err = fmt.Fprint(ui.writer, "\nВведите категорию слов из списка: ")
	if err != nil {
		return "", err
	}

	err = ui.writer.Flush()
	if err != nil {
		return "", err
	}

	byteChosenCategory, _, err := ui.reader.ReadLine()
	if err != nil {
		return "", err
	}

	chosenCategory := strings.ToUpper(string(byteChosenCategory))
	_, ok := categories[chosenCategory]

	if !ok {
		_, err = fmt.Fprintln(ui.writer, "Введена некорректная категория, она будет выбрана случайно")
		if err != nil {
			return "", err
		}

		err = ui.writer.Flush()
		if err != nil {
			return "", err
		}

		randomPosition, err := getRandomInt(0, len(categories))
		if err != nil {
			return "", err
		}

		// Поиск случайной категории. Зависит от сгенерированного randomPosition и порядка обхода map.
		i := 0
		for key := range categories {
			if i == randomPosition {
				chosenCategory = key
			}

			i++
		}
	}

	return chosenCategory, nil
}

// Возвращает случайнок слово и его подсказку из текущей категории.
func (InitUI) ChooseWord(category string, categories map[string][]constants.GameWord) (word, hint string, err error) {
	randomIndex, err := getRandomInt(0, len(categories[category]))
	if err != nil {
		return "", "", err
	}

	return categories[category][randomIndex].Word, categories[category][randomIndex].Hint, nil
}

// Запрашивает ввода пути к json со словами.
// При пустом вводе берет Words.json.
func (ui *InitUI) InputJSONPath() (string, error) {
	_, err := fmt.Fprint(ui.writer, "Введите путь к json-файлу со словами либо будут выбраны слова по умолчанию: ")
	if err != nil {
		return "", err
	}

	err = ui.writer.Flush()
	if err != nil {
		return "", err
	}

	bytePath, _, err := ui.reader.ReadLine()
	if err != nil {
		return "", err
	}

	path := string(bytePath)
	if path == "" {
		path = constants.DefaultJSONPath
	}

	return path, nil
}

// Используется для остановки консоли.
// Чтобы перед ее очисткой игрок мог все прочитать.
func (ui InitUI) InputAccept() error {
	_, err := fmt.Fprint(ui.writer, "Нажмите Enter для продолжения...")
	if err != nil {
		return err
	}

	err = ui.writer.Flush()
	if err != nil {
		return err
	}

	_, _, err = ui.reader.ReadLine()

	return err
}

// Возвращает случайное число в полуинтервале [lower_bound, upper_bound).
func getRandomInt(lowerBound, upperBound int) (int, error) {
	diff := big.NewInt(int64(upperBound - lowerBound))
	randomValue, err := rand.Int(rand.Reader, diff)

	return int(randomValue.Int64()) + lowerBound, err
}
