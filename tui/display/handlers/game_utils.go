package handlers

import (
	"strings"

	"charm.land/lipgloss/v2"
	"gowordle.com/display/model"
)

func ToWordleDisplayString(s string) string {
	return strings.ToUpper(strings.Join(strings.Split(s, ""), " "))
}

func AddPaddingsToCurrentGuessedWord(s string, length ...int) string {
	l := 5
	if len(length) > 0 {
		l = length[0]
	}

	finalStrArray := []string{}
	for i := 0; i < l; i++ {
		if i >= len(s) || s[i] == 0 {
			finalStrArray = append(finalStrArray, "_")
		} else {
			finalStrArray = append(finalStrArray, string(s[i]))
		}
	}
	return strings.Join(finalStrArray, "")
}

func RenderWordleColoredLetter(letter model.LetterStatus, correct, misplaced, incorrect lipgloss.Style) string {
	if letter.IsCorrect {
		return correct.Render(string(letter.Letter))
	}
	if letter.IsInSolution {
		return misplaced.Render(string(letter.Letter))
	}
	return incorrect.Render(string(letter.Letter))
}
