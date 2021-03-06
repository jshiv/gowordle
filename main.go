package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/dustin/go-humanize"
	"github.com/dustin/go-humanize/english"
	"github.com/manifoldco/promptui"

	"github.com/fatih/color"
)

//go:embed sgb-words.txt
var sgb_words string

var fs embed.FS

// knownLetters: letters are known to be in a specific position
type knownLetters struct {
	firstLetter  string
	secondLetter string
	thridLetter  string
	fourthLetter string
	fifthLetter  string
}

// knownNotLetters: letters are known to not be in a specific position
// TODO: should compare to []string because multipule letters can be known to not be.
type knownIsNotLetters struct {
	firstLetterIsNot  []string
	secondLetterIsNot []string
	thridLetterIsNot  []string
	fourthLetterIsNot []string
	fifthLetterIsNot  []string
}

func main() {

	var knl knownIsNotLetters
	sl := strings.Split(string(sgb_words), "\n")

	color.Cyan("GO WORDLE")
	guesses, lettersHas, lettersHasNot, kl := getPrompt()
	for _, guess := range guesses {
		knl.getNotPositions(guess, kl, lettersHas)
	}
	if len(knl.firstLetterIsNot) != 0 {
		color.Cyan("1st letter is not: " + english.WordSeries(knl.firstLetterIsNot, ","))
	}
	if len(knl.secondLetterIsNot) != 0 {
		color.Cyan("2nd letter is not: " + english.WordSeries(knl.secondLetterIsNot, ","))
	}
	if len(knl.thridLetterIsNot) != 0 {
		color.Cyan("3rd letter is not: " + english.WordSeries(knl.thridLetterIsNot, ","))
	}
	if len(knl.fourthLetterIsNot) != 0 {
		color.Cyan("4th letter is not: " + english.WordSeries(knl.fourthLetterIsNot, ","))
	}
	if len(knl.fifthLetterIsNot) != 0 {
		color.Cyan("5th letter is not: " + english.WordSeries(knl.fifthLetterIsNot, ","))
	}

	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond) // Build our new spinner
	s.Start()                                                    // Start the spinner
	time.Sleep(2 * time.Second)                                  // Run for some time to simulate work
	s.Stop()
	hasWords := hasWordle(sl, lettersHas, lettersHasNot, kl, knl)
	if len(hasWords) == 0 {
		color.Red("No Words :-|")
		os.Exit(0)
	}

}

func hasWordle(words []string, lettersHas []string, lettersNot string, kl knownLetters, knl knownIsNotLetters) []string {
	color.Cyan("WORDLE Try...")
	var hasWords []string
	for _, s := range words {
		if hasChars(s, lettersHas) && !strings.ContainsAny(s, lettersNot) && hasPositions(s, kl, knl) {
			if runtime.GOOS == "windows" {
				fmt.Println(s)
			} else {
				color.Magenta(s)
			}
			hasWords = append(hasWords, s)
		}

	}
	return hasWords
}

func getPrompt() ([]string, []string, string, knownLetters) {

	var guesses []string
	for i := 1; i <= 6; i++ {
		validateGuess := func(input string) error {
			if len(input) != 5 && input != "" {
				return errors.New("Wordle Guess should be 5 letters")
			}
			return nil
		}

		pg := promptui.Prompt{
			Label:    humanize.Ordinal(i) + " Guess",
			Default:  "",
			Validate: validateGuess,
		}

		guess, err := pg.Run()
		if guess == "" {
			break
		}
		guesses = append(guesses, guess)

		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}
	}

	pa := promptui.Prompt{
		Label:   "Wordle has the letters",
		Default: "",
	}

	ra, err := pa.Run()

	if err == promptui.ErrInterrupt {
		os.Exit(0)
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	lettersHas := strings.Split(ra, "")

	lettersHasNot := strings.Join(guesses, "")
	for _, l := range lettersHas {
		lettersHasNot = strings.ReplaceAll(lettersHasNot, l, "")
	}

	lettersHasNot = removeDuplicates(lettersHasNot)

	color.Cyan("Wordle does not have the letters: " + strings.Join(strings.Split(lettersHasNot, ""), ","))

	p1 := promptui.Prompt{
		Label:   "Known First Letter",
		Default: "",
	}

	l1, err := p1.Run()

	if err == promptui.ErrInterrupt {
		os.Exit(0)
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	p2 := promptui.Prompt{
		Label:   "Known Second Letter",
		Default: "",
	}

	l2, err := p2.Run()

	if err == promptui.ErrInterrupt {
		os.Exit(0)
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	p3 := promptui.Prompt{
		Label:   "Known Third Letter",
		Default: "",
	}

	l3, err := p3.Run()

	if err == promptui.ErrInterrupt {
		os.Exit(0)
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	p4 := promptui.Prompt{
		Label:   "Known Fourth Letter",
		Default: "",
	}

	l4, err := p4.Run()

	if err == promptui.ErrInterrupt {
		os.Exit(0)
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	p5 := promptui.Prompt{
		Label:   "Known Fifth Letter",
		Default: "",
	}

	l5, err := p5.Run()

	if err == promptui.ErrInterrupt {
		os.Exit(0)
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	return guesses, lettersHas, lettersHasNot, knownLetters{firstLetter: l1, secondLetter: l2, thridLetter: l3, fourthLetter: l4, fifthLetter: l5}

}

func removeDuplicates(s string) string {
	set := make(map[string]bool)
	for _, c := range s {
		set[string(c)] = true
	}

	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}

	return strings.Join(keys, "")
}

func hasChars(s string, chars []string) bool {
	b := true
	for _, c := range chars {
		b = b && strings.Contains(s, c)
	}
	return b
}

func hasPositions(word string, kl knownLetters, knl knownIsNotLetters) bool {
	b := true
	if kl.firstLetter != "" {
		c := string(word[0])
		b = b && (c == kl.firstLetter)
	}

	if len(knl.firstLetterIsNot) > 0 {
		c := string(word[0])
		b = b && !strings.Contains(strings.Join(knl.firstLetterIsNot, ""), c)
	}

	if kl.secondLetter != "" {
		c := string(word[1])
		b = b && (c == kl.secondLetter)
	}

	if len(knl.secondLetterIsNot) > 0 {
		c := string(word[1])
		b = b && !strings.Contains(strings.Join(knl.secondLetterIsNot, ""), c)
	}

	if kl.thridLetter != "" {
		c := string(word[2])
		b = b && (c == kl.thridLetter)
	}

	if len(knl.thridLetterIsNot) > 0 {
		c := string(word[2])
		b = b && !strings.Contains(strings.Join(knl.thridLetterIsNot, ""), c)
	}

	if kl.fourthLetter != "" {
		c := string(word[3])
		b = b && (c == kl.fourthLetter)
	}

	if len(knl.fourthLetterIsNot) > 0 {
		c := string(word[3])
		b = b && !strings.Contains(strings.Join(knl.fourthLetterIsNot, ""), c)
	}

	if kl.fifthLetter != "" {
		c := string(word[4])
		b = b && (c == kl.fifthLetter)
	}

	if len(knl.fifthLetterIsNot) > 0 {
		c := string(word[4])
		b = b && !strings.Contains(strings.Join(knl.fifthLetterIsNot, ""), c)
	}

	return b

}

func (knl *knownIsNotLetters) getNotPositions(guess string, kl knownLetters, hasLetters []string) {

	for _, l := range hasLetters {
		hasLetterIndex := strings.Index(guess, l)
		switch hasLetterIndex {
		case -1:
			continue
		// the first letter we know of
		case 0:
			// does not match the position we know it is in
			if kl.firstLetter != l {
				knl.firstLetterIsNot = append(knl.firstLetterIsNot, l)
			}
		case 1:

			if kl.secondLetter != l {
				knl.secondLetterIsNot = append(knl.secondLetterIsNot, l)
			}

		case 2:

			if kl.thridLetter != l {
				knl.thridLetterIsNot = append(knl.thridLetterIsNot, l)
			}

		case 3:

			if kl.fourthLetter != l {
				knl.fourthLetterIsNot = append(knl.fourthLetterIsNot, l)
			}

		case 4:

			if kl.fifthLetter != l {
				knl.fifthLetterIsNot = append(knl.fifthLetterIsNot, l)
			}
		}
	}

}
