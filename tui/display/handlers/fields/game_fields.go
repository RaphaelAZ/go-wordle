package fields

import (
	"strings"

	"gowordle.com/display/model"
)

func CompareInputToSolution(input string, solution string) []model.LetterStatus {
	input = strings.ToUpper(input)
	solution = strings.ToUpper(solution)

	if len(solution) == 0 || len(input) != len(solution) {
		return nil
	}

	var final []model.LetterStatus
	for i := 0; i < len(input); i++ {
		letter := input[i]
		if letter == solution[i] {
			final = append(final, model.LetterStatus{IsCorrect: true, IsInSolution: true, Letter: letter})
		} else if strings.ContainsRune(solution, rune(letter)) {
			final = append(final, model.LetterStatus{IsCorrect: false, IsInSolution: true, Letter: letter})
		} else {
			final = append(final, model.LetterStatus{IsCorrect: false, IsInSolution: false, Letter: letter})
		}
	}
	return final
}

func AppendGuessLetter(m model.State, letter string) model.State {
	if len([]rune(m.Game.CurrentGuess)) >= 5 {
		return m
	}
	m.Game.CurrentGuess += strings.ToUpper(letter)
	return m
}

func DeleteGuessLetter(m model.State) model.State {
	r := []rune(m.Game.CurrentGuess)
	if len(r) > 0 {
		m.Game.CurrentGuess = string(r[:len(r)-1])
	}
	return m
}

func SubmitGuess(m model.State) model.State {
	if len([]rune(m.Game.CurrentGuess)) != 5 {
		return m
	}
	guess := m.Game.CurrentGuess

	m.Game.TriedWords = append(m.Game.TriedWords, CompareInputToSolution(guess, m.Game.WordToGuess))
	m.Game.CurrentGuess = ""

	if strings.EqualFold(guess, m.Game.WordToGuess) {
		m.Game.Status = model.GameWon
	} else if len(m.Game.TriedWords) >= 6 {
		m.Game.Status = model.GameLost
	}
	return m
}
