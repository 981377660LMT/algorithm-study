// https://github.com/ftbe/dawg
// 带模糊搜索功能的 DAWG 实现.
//
// 优点：节省空间，提高搜索效率.
// 缺点：构建过程复杂且耗时，不适合频繁变动的数据集.
//
// Package dawg implements a Directed Acyclic Word Graph, with fuzzy search of words in the graph.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	words := []string{"ag", "agi", "agin", "aging", "aging", "aging2"}
	graph := CreateDAWG(words)
	words, err := graph.Search("aging", 1, 50, true, true)
	if err != nil {
		// Do something
		return
	}
	for _, word := range words {
		fmt.Println(word)
	}
}

// DAWG is used to store the representation of the Directly Acyclic Word Graph
type DAWG struct {
	initialState *state // root
	nodesCount   uint64
}

// letterNode 是一颗二叉搜索树（BST）。
// 一方面使用二叉搜索树做“精确查找”子字符，另一方面又用单链表方便“顺序遍历”。
type letterNode struct {
	char  rune // 该字母边上存储的字符（`rune` 以兼容 Unicode）
	state *state

	// Tree, allow for O(log(n)) search operations
	left  *letterNode
	right *letterNode

	// Linked list, allow for a quick iteration on all the sub-letters of a state
	next *letterNode
}

// 用来描述 DAWG 中的一个节点（state）。
type state struct {
	final bool // 是某个单词的结束节点（类似 Trie 中的“终止标志”）

	letter       *letterNode // 指向一棵由 `letter` 组成的**二叉搜索树**（BST），同时也通过 `next` 串成**链表**。这在查找时可以 O(log n) 搜索，也可以 O(n) 遍历
	lettersCount int         // 当前 `state` 下有多少个子边（子节点）

	next       *state      // 用于将相同“层级”的 state 串成链表，后续合并时会用到
	fromLetter *letterNode // 指向从哪个 letter 来到这个 state（在合并时使用）
	number     uint64      // 序列化/反序列化用的唯一标识
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

	for curLetter := state.letter; curLetter != nil; curLetter = curLetter.next {
		if !otherState.containsLetter(curLetter) {
			return false
		}
	}

	return true
}

// Check if this state contains this letter (in O(log(n)) time)
func (state *state) containsLetter(letter *letterNode) (containsLetter bool) {
	curLetter := state.letter
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
func (state *state) getletter(letter rune) *letterNode {
	curLetter := state.letter
	for curLetter != nil && curLetter.char != letter {
		if curLetter.char < letter {
			curLetter = curLetter.left
		} else {
			curLetter = curLetter.right
		}
	}
	return curLetter
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

// !模糊搜索.
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

// !构建完 Trie 后，需要进行最重要的一步：**合并重复子树**（压缩成一个 DAWG）。
// 压缩 Trie，合并重复的节点，返回删除的节点数.
func compressTrie(initialState *state, maxWordSize int) (deletedNodes uint64) {
	levels := make([]*state, maxWordSize)
	// !1) analyseSubTrie: 把各层节点串起来，结果保存在 levels 中
	// !2) 在每层里，合并等价的 state
	//    for curState := levels[i]; curState != nil && curState.next != nil; curState = curState.next {
	//        for previousState, sameState := curState, curState.next; sameState != nil; sameState = sameState.next {
	//            if curState.equals(sameState) {
	//                // 合并
	//                previousState.next = sameState.next
	//                sameState.letter.state = curState
	//                deletedNodes++
	//            } else {
	//                previousState = sameState
	//            }
	//        }
	//    }

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

		for curLetter := initialState.letter; curLetter != nil; curLetter = curLetter.next {
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
					// 将 sameState 合并到 curState
					previousState.next = sameState.next
					sameState.fromLetter.state = curState
					deletedNodes++
				} else {
					previousState = sameState
				}
			}
		}
	}
	return
}

// 串联各层节点，返回当前层的深度.
func analyseSubTrie(curState *state, levels []*state, channels []chan int) (subLevels int) {
	var curLevel int = 0
	if curState.lettersCount != 0 {
		for curLetter := curState.letter; curLetter != nil; curLetter = curLetter.next {
			curSubLevels := analyseSubTrie(curLetter.state, levels, channels)
			if curSubLevels > curLevel {
				curLevel = curSubLevels
			}
		}
	}

	<-channels[curLevel]
	curState.next = levels[curLevel] // !串联同层节点
	levels[curLevel] = curState
	channels[curLevel] <- 1

	return curLevel + 1
}

// Add a new word to the Trie
func addWord(initialState *state, word string) (newEndState bool, wordSize int, createdNodes uint64) {
	curState := initialState

	for _, l := range word {
		// 1) 在 curState.letters（BST）里找到对应字符 'l' 的 letter
		//    若无则新建 letter
		// 2) 若 letter.state == nil 表示还没有下一层 state，就新建一个
		//    并记录 createdNodes++
		// 3) 移动 curState = letter.state
		// 4) wordSize++ (统计当前单词的字符数)

		// !将字符 'l' 添加到 curState 的 letters 中
		var curLetter *letterNode
		if curState.letter == nil {
			curLetter = &letterNode{char: l}
			curState.letter = curLetter
		} else {
			for curLetter = curState.letter; curLetter.char != l; {
				if curLetter.char < l {
					if curLetter.left == nil {
						curLetter.left = &letterNode{char: l}
					}
					curLetter = curLetter.left
				} else {
					if curLetter.right == nil {
						curLetter.right = &letterNode{char: l}
					}
					curLetter = curLetter.right
				}
			}
		}

		// !新建 state
		if curLetter.state == nil {
			curLetter.state = &state{final: false, fromLetter: curLetter}
			createdNodes++
			curState.lettersCount++
			if !curState.final && curState.lettersCount == 1 || curState.lettersCount > 1 {
				newEndState = true
			}
			if curLetter != curState.letter {
				curLetter.next = curState.letter.next
				curState.letter.next = curLetter
			}
		}

		curState = curLetter.state
		wordSize++ // !We can't use len() on UTF-8 strings
	}

	curState.final = true
	return
}

// 把不同分支递归得到的结果链表合并在一起，记录返回结果.
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
		// 依次读取每行，构建对应的 `state` 对象，并把 `char -> linkedNodeNumber` 的关系恢复到各自的 `letter`
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

				if states[nodeNumber].letter == nil {
					states[nodeNumber].letter = &letterNode{char: char, state: states[linkedNodeNumber]}
				} else {
					for curLetter := states[nodeNumber].letter; curLetter.char != char; {
						if curLetter.char < char {
							if curLetter.left == nil {
								curLetter.left = &letterNode{char: char, state: states[linkedNodeNumber]}
								curLetter.left.next = states[nodeNumber].letter.next
								states[nodeNumber].letter.next = curLetter.left
							}
							curLetter = curLetter.left
						} else {
							if curLetter.right == nil {
								curLetter.right = &letterNode{char: char, state: states[linkedNodeNumber]}
								curLetter.right.next = states[nodeNumber].letter.next
								states[nodeNumber].letter.next = curLetter.right
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

// 每个 `state`，依次为其分配一个 `number`（递增），并将 `final` 标志、子 `letter`（字符 + 状态编号）写入文件
func saveSubTrieToFile(file *os.File, curState *state, nodeNumber *uint64) (err error) {
	for curLetter := curState.letter; curLetter != nil; curLetter = curLetter.next {
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
		for curLetter := curState.letter; curLetter != nil; curLetter = curLetter.next {
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
			letter := state.letter
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

// 模糊搜索的递归实现.
//
//	`start`: 当前已匹配的前缀（在递归过程中不断增长或回退）。
//	`end`: 还没匹配的剩余字符（从 `word` 中读取）。
//	`levenshteinDistance`: 还剩多少可“容错”编辑的机会。
//	`ignoreChar`: 递归里处理“替换字符”或“跳过字符”时，用来记住被替换/跳过的字符。
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
			for letter := state.letter; letter != nil; letter = letter.next {
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
		for letter := state.letter; letter != nil; letter = letter.next {
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
