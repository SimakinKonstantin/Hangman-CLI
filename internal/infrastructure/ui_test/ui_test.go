package ui_test

import (
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/constants"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/ui"
	"github.com/stretchr/testify/assert"
)

// Тесты ui игры.
func TestNewGameUI(t *testing.T) {
	ui1 := ui.NewGameUI("12345", strings.NewReader("\n"), io.Discard)
	ui2 := ui.NewGameUI("1", strings.NewReader("\n"), io.Discard)
	ui3 := ui.NewGameUI("", strings.NewReader("\n"), io.Discard)

	assert.Equal(t, "_____", ui1.GetProgress(), "Проверка создания строки с прогрессом")
	assert.Equal(t, "_", ui2.GetProgress(), "Проверка создания строки с прогрессом")
	assert.Equal(t, "", ui3.GetProgress(), "Проверка создания строки с прогрессом при некорректном слове")
}

func TestUpdateProgress(t *testing.T) {
	rightWord1 := "12345"
	ui1 := ui.NewGameUI(rightWord1, strings.NewReader("\n"), io.Discard)
	indexes1 := []int{0, 2, 3}
	lettersLeft1 := ui1.UpdateProgress(indexes1, []rune(rightWord1))

	rightWord2 := "AAA"
	ui2 := ui.NewGameUI(rightWord2, strings.NewReader("\n"), io.Discard)
	indexes2 := []int{0, 1, 2}
	lettersLeft2 := ui2.UpdateProgress(indexes2, []rune(rightWord2))

	assert.Equal(t, 2, lettersLeft1, "Должны остаться неугаданные буквы")
	assert.Equal(t, "1_34_", ui1.GetProgress(), "Должны остаться открыться часть слова")

	assert.Equal(t, 0, lettersLeft2, "Не должно остаться неугаданных букв")
	assert.Equal(t, "AAA", ui2.GetProgress(), "Все буквы должны открыться")
}

func TestInputLetter(t *testing.T) {
	ui1 := ui.NewGameUI("", strings.NewReader("W"), io.Discard)
	ui2 := ui.NewGameUI("", strings.NewReader("ъ"), io.Discard)

	val1, _ := ui1.InputLetter()
	val2, _ := ui2.InputLetter()

	assert.Equal(t, 'W', val1, "Должна быть в верхнем регистре")
	assert.Equal(t, 'Ъ', val2, "Должна быть в верхнем регистре")
}

// Тесты ui инициализации.
func TestInputDifficulty(t *testing.T) {
	ui1 := ui.NewInitUI(strings.NewReader(strconv.Itoa(constants.MaxErrors)), io.Discard)
	ui2 := ui.NewInitUI(strings.NewReader("0"), io.Discard)
	ui3 := ui.NewInitUI(strings.NewReader("qwerty"), io.Discard)
	ui4 := ui.NewInitUI(strings.NewReader("?"), io.Discard)

	val1, _ := ui1.InputDifficulty()
	val2, _ := ui2.InputDifficulty()
	val3, _ := ui3.InputDifficulty()
	val4, _ := ui4.InputDifficulty()

	assert.Equal(t, constants.MaxErrors, val1, "Корректное значение на входе")
	assert.LessOrEqual(t, val2, constants.MaxErrors, "Проверка получения корректного случайного значения")
	assert.GreaterOrEqual(t, val2, 0, "Проверка получения корректного случайного значения")
	assert.Equal(t, 0, val3, "Строка вместо числа")
	assert.Equal(t, 0, val4, "Символ вместо числа")
}

func TestInputCategory(t *testing.T) {
	categories := map[string][]constants.GameWord{"cat1": {}}
	ui1 := ui.NewInitUI(strings.NewReader("cat1"), io.Discard)
	ui2 := ui.NewInitUI(strings.NewReader("false_cat"), io.Discard)
	ui3 := ui.NewInitUI(strings.NewReader("\n"), io.Discard)

	val1, _ := ui1.InputCategory(categories)
	val2, _ := ui2.InputCategory(categories)
	val3, _ := ui3.InputCategory(categories)

	assert.Equal(t, "cat1", val1, "Корректная категория из списка")
	assert.Equal(t, "cat1", val2, "Несуществующая категория из списка (должна замениться на случайную из списка)")
	assert.Equal(t, "cat1", val3, "Пустая строка вместо ввода")
}

func TestInputJsonPath(t *testing.T) {
	ui1 := ui.NewInitUI(strings.NewReader("./path.json"), io.Discard)
	ui2 := ui.NewInitUI(strings.NewReader("\n"), io.Discard)

	val1, _ := ui1.InputJSONPath()
	val2, _ := ui2.InputJSONPath()

	assert.Equal(t, "./path.json", val1, "Введен корректный путь")
	assert.Equal(t, constants.DefaultJSONPath, val2, "Неверный ввод, должен выбраться путь по умолчанию")
}
