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

	result := make([]model.LetterStatus, len(input))
	solutionUsed := make([]bool, len(solution))

	// Première passe : positions correctes (vert)
	for i := 0; i < len(input); i++ {
		if input[i] == solution[i] {
			result[i] = model.LetterStatus{IsCorrect: true, IsInSolution: true, Letter: input[i]}
			solutionUsed[i] = true
		}
	}

	// Deuxième passe : lettres présentes mais mal placées (jaune)
	// On ne consomme une lettre de la solution que si elle n'est pas déjà utilisée
	for i := 0; i < len(input); i++ {
		if result[i].IsCorrect {
			continue
		}
		letter := input[i]
		found := false
		for j := 0; j < len(solution); j++ {
			if false == solutionUsed[j] && solution[j] == letter {
				solutionUsed[j] = true
				found = true
				break
			}
		}
		result[i] = model.LetterStatus{IsCorrect: false, IsInSolution: found, Letter: letter}
	}

	return result
}

func AppendGuessLetter(m model.State, letter string) model.State {
	if len([]rune(m.Game.CurrentGuess)) >= 5 && model.GamePlaying == m.Game.Status {
		return m
	}
	m.Game.CurrentGuess += strings.ToUpper(letter)
	return m
}

func DeleteGuessLetter(m model.State) model.State {
	r := []rune(m.Game.CurrentGuess)
	if len(r) > 0 && model.GamePlaying == m.Game.Status {
		m.Game.CurrentGuess = string(r[:len(r)-1])
	}
	return m
}

func SubmitGuess(m model.State) model.State {
	if len([]rune(m.Game.CurrentGuess)) != 5 && model.GamePlaying == m.Game.Status {
		return m
	}
	guess := m.Game.CurrentGuess
	result := CompareInputToSolution(guess, m.Game.WordToGuess)

	m.Game.TriedWords = append(m.Game.TriedWords, result)
	m.Game.CurrentGuess = ""

	if m.Game.UsedLetters == nil {
		m.Game.UsedLetters = make(map[byte]model.LetterStatus)
	}
	for _, ls := range result {
		existing, exists := m.Game.UsedLetters[ls.Letter]
		if !exists || (!existing.IsCorrect && ls.IsCorrect) || (!existing.IsInSolution && !existing.IsCorrect && ls.IsInSolution) {
			m.Game.UsedLetters[ls.Letter] = ls
		}
	}

	if strings.EqualFold(guess, m.Game.WordToGuess) {
		m.Game.Status = model.GameWon
	} else if len(m.Game.TriedWords) >= 6 {
		m.Game.Status = model.GameLost
	}
	return m
}
