package session

import "fmt"

type ReadJSONError struct {
	msg string
	err error
}

func (readErr ReadJSONError) Error() string {
	return fmt.Sprintf("%s: %v", readErr.msg, readErr.err)
}

func (readErr ReadJSONError) Unwrap() error { return readErr.err }

// Ошибка при неправильном формате json-файла, не оборачивает ошибки из других функций.
type ParseJSONError struct {
	msg string
}

func (parseErr ParseJSONError) Error() string {
	return parseErr.msg
}

type InputError struct {
	msg string
	err error
}

func (inputErr InputError) Error() string {
	return fmt.Sprintf("%s: %v", inputErr.msg, inputErr.err)
}

func (inputErr InputError) Unwrap() error { return inputErr.err }

type OutputError struct {
	msg string
	err error
}

func (outputErr OutputError) Error() string {
	return fmt.Sprintf("%s: %v", outputErr.msg, outputErr.err)
}

func (outputErr OutputError) Unwrap() error { return outputErr.err }

type ProcessInputError struct {
	msg string
	err error
}

func (processInputErr ProcessInputError) Error() string {
	return fmt.Sprintf("%s: %v", processInputErr.msg, processInputErr.err)
}

func (processInputErr ProcessInputError) Unwrap() error { return processInputErr.err }
