package constants

type GameWord struct {
	Word string
	Hint string
}

// Максимальное число ошибок, которые может допустить игрок.
// При изменении этого параметра нужно добавить новые кадры виселицы в ui.PrintPicture.
const MaxErrors = 10

// Путь к json со словами, который выберется при неправильном вводе пользователя.
const DefaultJSONPath = "./internal/infrastructure/default/Words.json"
