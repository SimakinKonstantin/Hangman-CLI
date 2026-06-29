package gamelogic

import (
	"unicode/utf8"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/constants"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/ui"
)

type GameLogic struct {

	// Правильное слово, которое отгадывает пользователь.
	rightAnswer string

	// Подсказка к слову.
	hint string

	// Текущее кол-во ошибок.
	errorsCounter int
}

// Возвращает новый экземпляр GameLogic.
func NewGameLogic(rightAnswer, hint string, errorsCounter int) GameLogic {
	return GameLogic{rightAnswer, hint, errorsCounter}
}

// Возвращает текущее кол-во ошибок.
func (gl GameLogic) GetErrorsCounter() int {
	return gl.errorsCounter
}

// Возвращает текст подсказки.
func (gl GameLogic) GetHint() string {
	return gl.hint
}

// Проверяет введенную пользователем букву возвращает срез с индексами букв, которые пользователь угадал.
func (gl GameLogic) CheckLetter(letter rune) []int {
	indexes := make([]int, 0)
	length := utf8.RuneCountInString(gl.rightAnswer)

	for i := 0; i < length; i++ {
		var curLetter, size = utf8.DecodeRuneInString(gl.rightAnswer)
		if curLetter == letter {
			indexes = append(indexes, i)
		}

		gl.rightAnswer = gl.rightAnswer[size:]
	}

	return indexes
}

// Обрабатывает индексы слова, где находится угаданная игроком буква.
// Первое возвращаемое значение - правильно ли угадал букву, второе - нужно ли продолжать игру.
func (gl *GameLogic) ProcessInput(indexes []int, gameUI *ui.GameUI) (isCorrect, willContinue bool, err error) {
	if len(indexes) > 0 {
		if gameUI.UpdateProgress(indexes, []rune(gl.rightAnswer)) == 0 {
			err := gameUI.PrintWin()
			return true, false, err
		}

		return true, true, nil
	}

	if len(indexes) == 0 {
		gl.errorsCounter++
		if gl.errorsCounter >= constants.MaxErrors {
			err := gameUI.PrintLose()
			return false, false, err
		}

		return false, true, nil
	}

	return false, false, nil
}
