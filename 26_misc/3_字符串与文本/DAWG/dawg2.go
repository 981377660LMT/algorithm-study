// https://github.com/ftbe/dawg
// Package dawg implements a Directed Acyclic Word Graph, with fuzzy search of words in the graph.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// DAWG is used to store the representation of the Directly Acyclic Word Graph
type DAWG struct {
	initialState *state
	nodesCount   uint64
}

type letter struct {
	char  rune // Yay ! Unicode !
	state *state

	// Tree, allow for O(log(n)) search operations
	left  *letter
	right *letter

	// Linked list, allow for a quick iteration on all the sub-letters of a state
	next *letter
}

type state struct {
	final bool

	letters      *letter // Root of the letter tree and the letter linked list
	lettersCount int     // Number of letters in the tree/linked list

	next   *state  // Linked list of all the state on the same level (used to merge duplicate nodes)
	letter *letter // The letter this state comes from (used to merge duplicate nodes)
	number uint64  // The number of this state (used to save the DAWG to a file)
}

// Linked list of words
type word struct {
	content  string
	nextWord *word
}

// Check if two states are equals.
// Two states are equals :
//   - if they are on the same level
//   - if they are both final/non-final
//   - if they have the same sub-letters going to the same states.
func (state *state) equals(otherState *state) (equals bool) {
	if state.final != otherState.final || state.lettersCount != otherState.lettersCount {
		return false
	}

	for curLetter := state.letters; curLetter != nil; curLetter = curLetter.next {
		if !otherState.containsLetter(curLetter) {
			return false
		}
	}

	return true
}

// Check if this state contains this letter (in O(log(n)) time)
func (state *state) containsLetter(letter *letter) (containsLetter bool) {
	curLetter := state.letters
	for curLetter != nil && curLetter.char != letter.char {
		if curLetter.char < letter.char {
			curLetter = curLetter.left
		} else {
			curLetter = curLetter.right
		}
	}
	return curLetter != nil && curLetter.state == letter.state
}

// Get a letter from the state (in O(log(n)) time)
func (state *state) getletter(letter rune) *letter {
	curLetter := state.letters
	for curLetter != nil && curLetter.char != letter {
		if curLetter.char < letter {
			curLetter = curLetter.left
		} else {
			curLetter = curLetter.right
		}
	}
	return curLetter
}

// Create a new DAWG by loading the words from a file.
// The file must be UTF-8 encoded, one word per line.
func CreateDAWGFromFile(fileName string) (dawg *DAWG, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	initialState := &state{final: false}
	var nbNodes uint64 = 1
	maxWordSize := 0
	for scanner.Scan() {
		_, size, createdNodes := addWord(initialState, scanner.Text())
		if size > maxWordSize {
			maxWordSize = size
		}
		nbNodes += createdNodes
	}
	if err = scanner.Err(); err != nil {
		return
	}
	nbNodes -= compressTrie(initialState, maxWordSize)
	return &DAWG{initialState: initialState, nodesCount: nbNodes}, nil
}

// Create a new DAWG by loading the words from an array.
func CreateDAWG(words []string) *DAWG {
	initialState := &state{final: false}
	var nbNodes uint64 = 1
	maxWordSize := 0
	for _, word := range words {
		_, size, createdNodes := addWord(initialState, word)
		if size > maxWordSize {
			maxWordSize = size
		}
		nbNodes += createdNodes
	}
	nbNodes -= compressTrie(initialState, maxWordSize)
	return &DAWG{initialState: initialState, nodesCount: nbNodes}
}

func compressTrie(initialState *state, maxWordSize int) (deletedNodes uint64) {
	// First, analyse the trie recursively to create a linked list of all the state on the same level
	levels := make([]*state, maxWordSize)
	if initialState.lettersCount != 0 {
		channels := make([]chan int, maxWordSize) // To synchronize the access to levels
		done := make(chan int, initialState.lettersCount)
		for i := 0; i < maxWordSize; i++ {
			channels[i] = make(chan int, 1)
			i := i
			go func() {
				channels[i] <- 1
			}()
		}
		for curLetter := initialState.letters; curLetter != nil; curLetter = curLetter.next {
			// Parallelize the treatment
			go func(curState *state) {
				analyseSubTrie(curState, levels, channels)
				done <- 1
			}(curLetter.state)
		}
		// Wait for the end of all goroutines
		for i := 0; i < initialState.lettersCount; i++ {
			<-done
		}
	}

	// For each level, merge the duplicates states
	for i := 0; i < maxWordSize; i++ {
		for curState := levels[i]; curState != nil && curState.next != nil; curState = curState.next {
			for previousState, sameState := curState, curState.next; sameState != nil; sameState = sameState.next {
				if curState.equals(sameState) {
					previousState.next = sameState.next
					sameState.letter.state = curState
					deletedNodes++
				} else {
					previousState = sameState
				}
			}
		}
	}
	return
}

func analyseSubTrie(curState *state, levels []*state, channels []chan int) (subLevels int) {
	var curLevel int = 0
	if curState.lettersCount != 0 {
		for curLetter := curState.letters; curLetter != nil; curLetter = curLetter.next {
			curSubLevels := analyseSubTrie(curLetter.state, levels, channels)
			if curSubLevels > curLevel {
				curLevel = curSubLevels
			}
		}
	}

	<-channels[curLevel]
	curState.next = levels[curLevel]
	levels[curLevel] = curState
	channels[curLevel] <- 1

	return curLevel + 1
}

// Add a new word to the Trie
func addWord(initialState *state, word string) (newEndState bool, wordSize int, createdNodes uint64) {
	curState := initialState
	for _, l := range word {
		var curLetter *letter
		if curState.letters == nil {
			curLetter = &letter{char: l}
			curState.letters = curLetter
		} else {
			for curLetter = curState.letters; curLetter.char != l; {
				if curLetter.char < l {
					if curLetter.left == nil {
						curLetter.left = &letter{char: l}
					}
					curLetter = curLetter.left
				} else {
					if curLetter.right == nil {
						curLetter.right = &letter{char: l}
					}
					curLetter = curLetter.right
				}
			}
		}
		if curLetter.state == nil {
			curLetter.state = &state{final: false, letter: curLetter}
			createdNodes++
			curState.lettersCount++
			if curState.final == false && curState.lettersCount == 1 || curState.lettersCount > 1 {
				newEndState = true
			}
			if curLetter != curState.letters {
				curLetter.next = curState.letters.next
				curState.letters.next = curLetter
			}
		}
		curState = curLetter.state
		wordSize++ // We can't use len() on UTF-8 strings
	}
	curState.final = true
	return
}

// Approximate string searching in the DAWG.
// levenshteinDistance is the maximum Levenshtein distance allowed beetween word and the words found in the DAWG.
// maxResults allow to limit the number of returned results (to reduce the time needed by the search)
// allowAdd and allowDelete specify if the returned words can have insertions/deletions of letters
func (dawg *DAWG) Search(word string, levenshteinDistance int, maxResults int, allowAdd bool, allowDelete bool) (words []string, err error) {
	wordsFound, _, wordsSize, err := searchSubString(dawg.initialState, *bytes.NewBufferString(""), *bytes.NewBufferString(word), levenshteinDistance, maxResults, allowAdd, allowDelete, 0)
	if err != nil {
		return
	}
	// Truncate if we have found more words than we need
	for ; wordsSize > maxResults; wordsSize-- {
		wordsFound = wordsFound.nextWord
	}
	// Transform to an array of strings
	words = make([]string, wordsSize)
	for ; wordsSize > 0; wordsSize-- {
		words[wordsSize-1] = wordsFound.content
		wordsFound = wordsFound.nextWord
	}
	return
}

func mergeWords(words1 *word, lastWord1 *word, wordsSize1 int, words2 *word, lastWord2 *word, wordsSize2 int) (words *word, lastWord *word, wordsSize int) {
	if words1 == nil {
		return words2, lastWord2, wordsSize2
	}
	lastWord1.nextWord = words2
	if lastWord2 == nil {
		return words1, lastWord1, wordsSize1 + wordsSize2
	} else {
		return words1, lastWord2, wordsSize1 + wordsSize2
	}
}

// Load from a file a DAWG saved by SaveToFile
func LoadDAWGFromFile(fileName string) (dawg *DAWG, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	var nbNodes uint64
	var initialState *state
	if scanner.Scan() {
		nbNodes, err = strconv.ParseUint(scanner.Text(), 10, 64)
		if err != nil {
			return
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}
	states := make([]*state, nbNodes)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) < 2 {
			err = errors.New("Incorrect node format : at least 2 fields expected.")
			return
		}
		var nodeNumber uint64
		nodeNumber, err = strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return
		}
		var finalNode bool
		finalNode, err = strconv.ParseBool(fields[1])
		if err != nil {
			return
		}

		states[nodeNumber] = &state{final: finalNode}
		initialState = states[nodeNumber]
		var char rune = 0
		for i, str := range fields[2:] {
			if i%2 == 0 {
				// It seems that char, _, _, err = strconv.UnquoteChar(str, 0) doesn't work, so we have to use Unquote before UnquoteChar
				var unquoted string
				unquoted, err = strconv.Unquote(str)
				if err != nil {
					return
				}
				char, _, _, err = strconv.UnquoteChar(unquoted, 0)
				if err != nil {
					return
				}
			} else {
				var linkedNodeNumber uint64
				linkedNodeNumber, err = strconv.ParseUint(str, 10, 64)
				if err != nil {
					return
				}

				states[nodeNumber].lettersCount = (i + 1) / 2

				if states[nodeNumber].letters == nil {
					states[nodeNumber].letters = &letter{char: char, state: states[linkedNodeNumber]}
				} else {
					for curLetter := states[nodeNumber].letters; curLetter.char != char; {
						if curLetter.char < char {
							if curLetter.left == nil {
								curLetter.left = &letter{char: char, state: states[linkedNodeNumber]}
								curLetter.left.next = states[nodeNumber].letters.next
								states[nodeNumber].letters.next = curLetter.left
							}
							curLetter = curLetter.left
						} else {
							if curLetter.right == nil {
								curLetter.right = &letter{char: char, state: states[linkedNodeNumber]}
								curLetter.right.next = states[nodeNumber].letters.next
								states[nodeNumber].letters.next = curLetter.right
							}
							curLetter = curLetter.right
						}
					}
				}
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}
	return &DAWG{initialState: initialState, nodesCount: nbNodes}, nil
}

// Save the DAWG to a file, usefull if you want to load it later without re-computing anything
func (dawg *DAWG) SaveToFile(fileName string) (err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return
	}

	if _, err = file.WriteString(strconv.FormatUint(dawg.nodesCount, 10)); err != nil {
		return
	}
	if _, err = file.WriteString("\n"); err != nil {
		return
	}

	var nodeNumber uint64 = 0
	// FIXME: if dawg.initialState.number != 0, the file was already saved, reinit all the numbers
	err = saveSubTrieToFile(file, dawg.initialState, &nodeNumber)
	return
}

func saveSubTrieToFile(file *os.File, curState *state, nodeNumber *uint64) (err error) {
	for curLetter := curState.letters; curLetter != nil; curLetter = curLetter.next {
		if curLetter.state.number == 0 {
			err = saveSubTrieToFile(file, curLetter.state, nodeNumber)
			if err != nil {
				return
			}
		}
	}
	if curState.number == 0 {
		(*nodeNumber)++
		curState.number = (*nodeNumber)
		if _, err = file.WriteString(strconv.FormatUint(curState.number-1, 10)); err != nil {
			return
		}
		if _, err = file.WriteString(" "); err != nil {
			return
		}
		if _, err = file.WriteString(strconv.FormatBool(curState.final)); err != nil {
			return
		}
		for curLetter := curState.letters; curLetter != nil; curLetter = curLetter.next {
			if _, err = file.WriteString(" "); err != nil {
				return
			}
			if _, err = file.WriteString(strconv.QuoteRune(curLetter.char)); err != nil {
				return
			}
			if _, err = file.WriteString(" "); err != nil {
				return
			}
			if _, err = file.WriteString(strconv.FormatUint(curLetter.state.number-1, 10)); err != nil {
				return
			}
		}
		if _, err = file.WriteString("\n"); err != nil {
			return
		}
	}
	return
}

func (dawg *DAWG) FindRandomWord(wordSize int) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// FIXME : infinite loop if no word of size wordSize
	// FIXME : highly inefficient
INFINITE:
	for {
		word := new(bytes.Buffer)
		state := dawg.initialState
		for i := 0; i < wordSize; i++ {
			if state.lettersCount == 0 { // That's bad
				continue INFINITE
			}
			var numLetter int
			if state.lettersCount == 1 {
				numLetter = 0
			} else {
				numLetter = r.Intn(state.lettersCount)
			}
			letter := state.letters
			for j := 0; j < numLetter; j++ {
				letter = letter.next
			}
			_, err := word.WriteRune(letter.char)
			if err != nil {
				return "", err
			}
			state = letter.state
		}
		if state.final {
			return word.String(), nil
		}
	}
}

func searchSubString(state *state, start bytes.Buffer, end bytes.Buffer, levenshteinDistance int, maxResults int, allowAdd bool, allowDelete bool, ignoreChar rune) (words *word, lastWord *word, wordsSize int, er error) {
	var char rune
	if end.Len() > 0 {
		char, _, er = end.ReadRune()
		if er != nil {
			return
		}
		if char != ignoreChar {
			if letter := state.getletter(char); letter != nil {
				runeLen, err := start.WriteRune(letter.char)
				if err != nil {
					return nil, nil, 0, err
				}
				foundWords, foundLastWord, foundWordsSize, err := searchSubString(letter.state, start, end, levenshteinDistance, maxResults, allowAdd, allowDelete, 0)
				if err != nil {
					return nil, nil, 0, err
				}
				words, lastWord, wordsSize = mergeWords(foundWords, foundLastWord, foundWordsSize, words, lastWord, wordsSize)
				if maxResults > 0 && wordsSize > maxResults {
					return
				}
				start.Truncate(start.Len() - runeLen) // Revert the WriteRune
			}
		}

		if levenshteinDistance > 0 {
			for letter := state.letters; letter != nil; letter = letter.next {
				if letter.char != char && letter.char != ignoreChar { // Change one letter
					runeLen, err := start.WriteRune(letter.char)
					if err != nil {
						return nil, nil, 0, err
					}
					foundWords, foundLastWord, foundWordsSize, err := searchSubString(letter.state, start, end, levenshteinDistance-1, maxResults, allowAdd, allowDelete, char)
					if err != nil {
						return nil, nil, 0, err
					}
					words, lastWord, wordsSize = mergeWords(foundWords, foundLastWord, foundWordsSize, words, lastWord, wordsSize)
					if maxResults > 0 && wordsSize > maxResults {
						return
					}
					start.Truncate(start.Len() - runeLen) // Revert the WriteRune
				}
			}
			if allowDelete {
				foundWords, foundLastWord, foundWordsSize, err := searchSubString(state, start, end, levenshteinDistance-1, maxResults, allowAdd, allowDelete, char) // Remove one letter
				if err != nil {
					return nil, nil, 0, err
				}
				words, lastWord, wordsSize = mergeWords(foundWords, foundLastWord, foundWordsSize, words, lastWord, wordsSize)
				if maxResults > 0 && wordsSize > maxResults {
					return
				}
			}
		}

		if err := end.UnreadRune(); err != nil { // Revert the ReadRune
			return nil, nil, 0, err
		}
	} else if state.final {
		words = &word{content: start.String(), nextWord: words}
		lastWord = words
		wordsSize = 1
	}

	if levenshteinDistance > 0 && allowAdd {
		for letter := state.letters; letter != nil; letter = letter.next {
			if letter.char != char && letter.char != ignoreChar { // Add one letter
				runeLen, err := start.WriteRune(letter.char)
				if err != nil {
					return nil, nil, 0, err
				}
				foundWords, foundLastWord, foundWordsSize, err := searchSubString(letter.state, start, end, levenshteinDistance-1, maxResults, allowAdd, allowDelete, 0)
				if err != nil {
					return nil, nil, 0, err
				}
				words, lastWord, wordsSize = mergeWords(foundWords, foundLastWord, foundWordsSize, words, lastWord, wordsSize)
				if maxResults > 0 && wordsSize > maxResults {
					return
				}
				start.Truncate(start.Len() - runeLen) // Revert the WriteRune
			}
		}
	}

	return
}
