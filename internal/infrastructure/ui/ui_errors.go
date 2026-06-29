package ui

import "fmt"

type TerminalClearError struct {
	msg string
	err error
}

func (terminalClearerr TerminalClearError) Error() string {
	return fmt.Sprintf("%s: %v", terminalClearerr.msg, terminalClearerr.err)
}

func (terminalClearerr TerminalClearError) Unwrap() error { return terminalClearerr.err }
