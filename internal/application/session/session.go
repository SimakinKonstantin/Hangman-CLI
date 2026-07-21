package session

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/gamelogic"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/constants"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/ui"
)

type Session struct {

	// Графический интерфейс игры.
	ui ui.GameUI

	// Логика игры.
	gameLogic gamelogic.GameLogic
}

// Возвращет новый экземпляр Session.
// В этой функции происходит инициализация сессии, когда выбираются категория слов, слово и уровень сложности.
func NewSession() (Session, error) {
	s := Session{}
	initializingUI := ui.NewInitUI(os.Stdin, os.Stdout)

	jsonPath, err := initializingUI.InputJSONPath()
	if err != nil {
		return s, InputError{"ошибка во время ввода пути к файлу со словами", err}
	}

	categoriesJSON, err := os.ReadFile(jsonPath)
	if err != nil {
		return s, ReadJSONError{"ошибка чтения json-файла", err}
	}

	categories, err := parseJSON(categoriesJSON)
	if err != nil {
		return s, fmt.Errorf("ошибка парсинга json: %v", err)
	}

	category, err := initializingUI.InputCategory(categories)
	if err != nil {
		return s, InputError{"ошибка при вводе названия категории", err}
	}

	rightAnswer, hint, err := initializingUI.ChooseWord(category, categories)
	if err != nil {
		return s, InputError{"ошибка при случайном выборе слова", err}
	}

	difficulty, err := initializingUI.InputDifficulty()
	if err != nil {
		return s, InputError{"ошибка при вводе уровня сложности", err}
	}

	s.gameLogic = gamelogic.NewGameLogic(rightAnswer, hint, constants.MaxErrors-difficulty)
	s.ui = ui.NewGameUI(rightAnswer, os.Stdin, os.Stdout)

	err = initializingUI.InputAccept()
	if err != nil {
		return s, err
	}

	err = s.ui.ClearScreen()

	return s, err
}

// Контролирует игру, взаимодействуя с логикой и графическим интерфейсом игры.
func (s Session) Run() error {
	attemptsLeft := constants.MaxErrors - s.gameLogic.GetErrorsCounter()

	// Правильной ли была буква. Когда = false, показывается оставшееся кол-во попыток.
	inputWasRight := false

	// Нужно ли продолжать игру.
	continueGame := true

	// Показывать ли подсказку.
	showHint := false

	for {
		if !continueGame {
			return nil
		}

		if showHint {
			err := s.ui.Print("Подсказка: ", s.gameLogic.GetHint())
			if err != nil {
				return OutputError{"ошибка вывода подсказки", err}
			}
		}

		err := s.ui.Print("Прогресс: ", s.ui.GetProgress())
		if err != nil {
			return OutputError{"ошибка вывода прогресса", err}
		}

		err = s.ui.PrintPicture(s.gameLogic.GetErrorsCounter())
		if err != nil {
			return OutputError{"ошибка отображения виселицы", err}
		}

		if !inputWasRight {
			err = s.ui.Print("Осталось попыток: ", attemptsLeft)
			if err != nil {
				return OutputError{"ошибка вывода кол-ва попыток", err}
			}
		}

		letter, err := s.ui.InputLetter()
		if err != nil {
			return InputError{"ошибка при вводе буквы", err}
		}

		if letter == '?' {
			showHint = true
			continue
		}

		indexes := s.gameLogic.CheckLetter(letter)

		inputWasRight, continueGame, err = s.gameLogic.ProcessInput(indexes, &s.ui)
		if err != nil {
			return ProcessInputError{"ошибка при обработке введенной буквы", err}
		}

		if !inputWasRight {
			attemptsLeft--
		}
	}
}

// Парсит json со словами\категориями, возвращает результат в виде map {категория : [слово1, слово2 ...]}.
func parseJSON(jsonString []byte) (map[string][]constants.GameWord, error) {
	// Преобразование файла, чтобы нигде не было проблем с регистром.
	tmpString := string(jsonString)
	jsonString = []byte(strings.ToUpper(tmpString))

	var parsedJSON map[string][]constants.GameWord
	err := json.Unmarshal(jsonString, &parsedJSON)

	if len(parsedJSON) == 0 {
		return parsedJSON, ParseJSONError{"в jsone нет слов\\категорий"}
	}

	// Проверка, что в Json нет пустых слов, категорий
	for key, val := range parsedJSON {
		// Некорректная категория
		if key == "" {
			return parsedJSON, ParseJSONError{"в jsone содержится категория, у которой название - пустое"}
		}

		// Пустая категория
		if len(val) == 0 {
			return parsedJSON, ParseJSONError{"в jsone содержится пустая категория"}
		}

		// Некорректное слово в категории
		for _, v := range val {
			if v.Word == "" || v.Hint == "" {
				return parsedJSON, ParseJSONError{"одно из слов в json неверно заполнено"}
			}
		}
	}

	return parsedJSON, err
}
