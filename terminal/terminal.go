package terminal

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func RawTerminal() *term.State {
	// Save the current terminal state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		RestoreTerminal(oldState)
		panic(err)
	}
	return oldState
}

func ResetTerminal(state *term.State) {
	// Restore the terminal to its previous state
	if err := term.Restore(int(os.Stdin.Fd()), state); err != nil {
		panic(err)
	}
}

func AlternateTerminal() {
	fmt.Print("\x1b[?1049h") // alternate screen
	fmt.Print("\x1b[2J")     // clear screen
	fmt.Print("\x1b[H")      // move to top-left
	fmt.Print("\x1b[?25l")   // hide cursor
}

func OriginalTerminal() {
	fmt.Print("\x1b[?1049l")
}

func SetupTerminal() *term.State {
	state := RawTerminal()
	AlternateTerminal()
	return state
}

func RestoreTerminal(state *term.State) {
	OriginalTerminal()
	ResetTerminal(state)
	fmt.Print("\x1b[?25h")
}
