package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type SingleKeyOption struct {
	Key   rune
	Label string
	Value string
}

var (
	singleKeyCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("239")).
				Padding(0, 1)
	singleKeyTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("110"))
	singleKeyHintStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245"))
	singleKeyKeyStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("215"))
)

// PromptSingleKeyChoice reads a single keypress and returns the matching value.
// It falls back to line-based input when stdin is not a terminal.
func PromptSingleKeyChoice(title string, options ...SingleKeyOption) (string, error) {
	lines := []string{singleKeyTitleStyle.Render(title)}
	for _, option := range options {
		lines = append(lines, fmt.Sprintf("%s %s", singleKeyKeyStyle.Render("["+strings.ToUpper(string(option.Key))+"]"), option.Label))
	}
	lines = append(lines, "", singleKeyHintStyle.Render("Press a single key to choose."))
	fmt.Println(singleKeyCardStyle.Render(strings.Join(lines, "\n")))

	for {
		key, err := readChoiceKey()
		if err != nil {
			return "", err
		}

		for _, option := range options {
			if unicode.ToLower(key) == unicode.ToLower(option.Key) {
				fmt.Println()
				return option.Value, nil
			}
		}
	}
}

func readChoiceKey() (rune, error) {
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		state, err := term.MakeRaw(fd)
		if err != nil {
			return 0, err
		}
		defer term.Restore(fd, state)

		var buf [3]byte
		n, err := os.Stdin.Read(buf[:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		return rune(buf[0]), nil
	}

	reader := bufio.NewReader(os.Stdin)
	key, _, err := reader.ReadRune()
	return key, err
}
