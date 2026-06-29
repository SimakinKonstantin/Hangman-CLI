package game_logic_test

import (
	"io"
	"os"
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/gamelogic"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/ui"
	"github.com/stretchr/testify/assert"
)

func TestCheckLetter(t *testing.T) {
	l1 := gamelogic.NewGameLogic("AbcAAbA", "", 0)
	l2 := gamelogic.NewGameLogic("Q", "", 0)
	l3 := gamelogic.NewGameLogic("", "", 0)

	ind1 := l1.CheckLetter('A')
	ind2 := l2.CheckLetter('Q')
	ind3 := l3.CheckLetter('X')

	assert.Equal(t, []int{0, 3, 4, 6}, ind1, "Стандартное слово с повторяющимимся буквами")
	assert.Equal(t, []int{0}, ind2, "Слово из одного символа")
	assert.Equal(t, []int{}, ind3, "Некорректное по размеру слово")
}

func TestProcessInput(t *testing.T) {
	l1 := gamelogic.NewGameLogic("12345", "", 0)

	// В данном случае не важно, какой reader использовать, т.к. он нигде не применяется
	ui1 := ui.NewGameUI("12345", os.Stdin, io.Discard)

	isRight1, willContinue1, _ := l1.ProcessInput([]int{0, 2, 4}, &ui1)

	assert.Equal(t, true, isRight1, "Были угаданы некоторые буквы")
	assert.Equal(t, true, willContinue1, "Остались неоткрытые буквы")

	isRight2, willContinue2, _ := l1.ProcessInput([]int{0, 1, 2, 3, 4}, &ui1)

	assert.Equal(t, true, isRight2, "Были угаданы некоторые буквы")
	assert.Equal(t, false, willContinue2, "Все буквы рагаданы")
}
