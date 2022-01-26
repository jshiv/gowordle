package main

import (
	"embed"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"

	"github.com/fatih/color"
)

//go:embed sgb-words.txt
var sgb_words string

var fs embed.FS

type knownLetters struct {
	firstLetter  string
	secondLetter string
	thridLetter  string
	fourthLetter string
	fifthLetter  string
}

func main() {

	sl := strings.Split(string(sgb_words), "\n")

	color.Cyan("GO WORDLE")
	lettersHas, lettersHasNot, kl := getPrompt()
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond) // Build our new spinner
	s.Start()                                                    // Start the spinner
	time.Sleep(2 * time.Second)                                  // Run for some time to simulate work
	s.Stop()
	hasWords := hasWordle(sl, lettersHas, lettersHasNot, kl)
	if len(hasWords) == 0 {
		color.Red("No Words :-|")
		os.Exit(0)
	}

}

func hasWordle(words []string, lettersHas []string, lettersNot string, kl knownLetters) []string {
	color.Cyan("WORDLE Try...")
	var hasWords []string
	for _, s := range words {
		if hasChars(s, lettersHas) && (strings.ContainsAny(s, lettersNot) == false) && hasPositions(s, kl) {
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

func getPrompt() ([]string, string, knownLetters) {

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

	pb := promptui.Prompt{
		Label:   "Wordle does not have the letters",
		Default: "",
	}

	rb, err := pb.Run()

	if err == promptui.ErrInterrupt {
		os.Exit(0)
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	lettersHasNot := rb

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

	return lettersHas, lettersHasNot, knownLetters{firstLetter: l1, secondLetter: l2, thridLetter: l3, fourthLetter: l4, fifthLetter: l5}

}

func hasChars(s string, chars []string) bool {
	b := true
	for _, c := range chars {
		b = b && strings.Contains(s, c)
	}
	return b
}

func hasPositions(s string, kl knownLetters) bool {
	b := true
	if kl.firstLetter != "" {
		b = b && (string(s[0]) == kl.firstLetter)
	}

	if kl.secondLetter != "" {
		b = b && (string(s[1]) == kl.secondLetter)
	}

	if kl.thridLetter != "" {
		b = b && (string(s[2]) == kl.thridLetter)
	}

	if kl.fourthLetter != "" {
		b = b && (string(s[3]) == kl.fourthLetter)
	}

	if kl.fifthLetter != "" {
		b = b && (string(s[4]) == kl.fifthLetter)
	}

	return b

}
