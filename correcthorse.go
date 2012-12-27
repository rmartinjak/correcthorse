package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	wordlistDir = "/usr/share/correcthorse"
)

// default options
var (
	optCount = 1
	optChars = 12
	optWords = 4
	optLists = stringSliceOpt{[]string{"english"}, false}
	optIncs  = stringSliceOpt{[]string{}, false}
	optSep   = ""
	optCamel = false
)

type stringSliceOpt struct {
	vals  []string
	isSet bool
}

// String and Set implement flag.Value
func (s *stringSliceOpt) String() string {
	if len(s.vals) == 0 {
		return "[none]"
	}
	return fmt.Sprint(strings.Join(s.vals, ","))
}
func (s *stringSliceOpt) Set(value string) error {
	// "clear" the slice if it hasn't been set before
	if !s.isSet {
		s.isSet = true
		s.vals = []string{}
	}

	// append value to the slice
	for _, val := range strings.Split(value, ",") {
		s.vals = append(s.vals, strings.TrimSpace(val))
	}
	return nil
}

// read all lines in a file into a string slice
func readLines(filename string) ([]string, error) {
	lines := make([]string, 0, 100)

	file, err := os.Open(filename)
	if err != nil {
		return []string{}, err
	}
	reader := bufio.NewReader(file)
	for err == nil {
		line, e := reader.ReadString('\n')
		if (e == nil || e == io.EOF) && line != "" {
			lines = append(lines, strings.TrimSpace(line))
		}
		err = e
	}
	file.Close()
	if err != io.EOF {
		return []string{}, err
	}
	return lines, nil
}

func loadWords(paths []string) ([][]string, error) {
	wordLists := make([][]string, len(paths))

	// read words from all provided pathnames
	// relative pathnames are searched below wordlistDir
	for i, p := range paths {
		if !filepath.IsAbs(p) {
			p = filepath.Join(wordlistDir, p)
		}
		lines, err := readLines(p)
		if err != nil {
			return [][]string{}, err
		}
		wordLists[i] = lines
	}
	return wordLists, nil
}

// shuffle a string slice
func shuffleStrings(words []string) []string {
	nWords := len(words)
	indices := rand.Perm(len(words))
	w := make([]string, nWords)
	for i, j := range indices {
		w[i] = words[j]
	}
	return w
}

func makePassphrase(wordLists [][]string) string {
	nChars := 0
	words := make([]string, 0)

	// add user-specified words
	for _, word := range optIncs.vals {
		nChars += len(word)
		words = append(words, word)
	}

	// add random word from random list until enough words and total characters
	for len(words) < optWords || nChars < optChars {
		list := wordLists[rand.Intn(len(wordLists))]
		word := list[rand.Intn(len(list))]
		nChars += len(word)
		words = append(words, word)
	}

	// capitalize first letters
	if optCamel {
		for i, w := range words {
			words[i] = strings.ToUpper(w[0:1]) + w[1:]
		}
	}

	return strings.Join(shuffleStrings(words), optSep)
}

func init() {
	// initialize argument parsing
	flag.Var(&optIncs, "inc", "word(s) to include in passphrase")
	flag.Var(&optLists, "list", "wordlist(s) to use (lower case L)")
	flag.IntVar(&optChars, "chars", optChars, "minimum number of characters")
	flag.IntVar(&optWords, "words", optWords, "minimum number of words")
	flag.BoolVar(&optCamel, "camel", optCamel, "use CamelCase")
	flag.StringVar(&optSep, "sep", optSep, "word separator")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	// last non-option argument is the number of passphrases to generate
	lastarg := flag.Arg(flag.NArg() - 1)
	c, err := strconv.Atoi(lastarg)
	if err == nil {
		optCount = c
	}

	// load wordlists
	wordLists, err := loadWords(optLists.vals)
	if err != nil {
		fmt.Println(err)
		return
	}

	// remove empty wordlists
	for i := 0; i < len(wordLists); i++ {
		if len(wordLists[i]) == 0 {
			L := len(wordLists) - 1
			wordLists[i] = wordLists[L]
			wordLists = wordLists[:L]
		}
	}

	// at least one non-empty list is needed
	if len(wordLists) == 0 {
		fmt.Println("no non-empty wordlist found")
	}

	// generate passphrases
	ch := make(chan string, optCount)
	gen := func() {
		for i := 0; i < optCount; i++ {
			ch <- makePassphrase(wordLists)
		}
		close(ch)
	}
	go gen()

	// print all passphrases
	for p := range ch {
		fmt.Println(p)
	}

	return
}
