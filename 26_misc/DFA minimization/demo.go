package main

import (
	"fmt"
)

// DFA 结构体：states, alphabet, transition, startState, acceptStates
type DFA struct {
	States       []int                // 状态集合，用 int 表示状态编号
	Alphabet     []rune               // 字母表
	Transition   map[int]map[rune]int // Transition[s][a] = t
	Start        int                  // 初始态
	AcceptStates map[int]bool         // 接受态集合
}

// HopcroftMinimize 实现 Hopcroft 算法，返回一个新的最小化 DFA
func HopcroftMinimize(dfa *DFA) *DFA {
	// Step 0: 如果需要，先去除不可达状态（可选）
	// 为简化，假设传入的 dfa 已经去除了不可达状态

	// 1. 初始划分：P = {AcceptStates, NonAcceptStates}
	partition := make([][]int, 0) // 存储每个 block 的状态列表
	blockID := make(map[int]int)  // state -> 哪个 block 的索引

	acceptBlock := make([]int, 0)
	nonAcceptBlock := make([]int, 0)

	for _, s := range dfa.States {
		if dfa.AcceptStates[s] {
			acceptBlock = append(acceptBlock, s)
		} else {
			nonAcceptBlock = append(nonAcceptBlock, s)
		}
	}
	if len(acceptBlock) > 0 {
		partition = append(partition, acceptBlock)
		for _, st := range acceptBlock {
			blockID[st] = 0
		}
	}
	if len(nonAcceptBlock) > 0 {
		partition = append(partition, nonAcceptBlock)
		for _, st := range nonAcceptBlock {
			blockID[st] = len(partition) - 1
		}
	}

	// 2. 建立 Reverse Transition: revTrans[symbol][t] = list of states that go to t on symbol
	revTrans := make(map[rune]map[int][]int)
	for _, a := range dfa.Alphabet {
		revTrans[a] = make(map[int][]int)
	}
	for _, s := range dfa.States {
		for _, a := range dfa.Alphabet {
			t := dfa.Transition[s][a]
			revTrans[a][t] = append(revTrans[a][t], s)
		}
	}

	// 3. 初始化队列 Q：将最小的 block 加入队列
	type blockSymbolPair struct {
		blockIndex int
		symbol     rune
	}
	Q := make([]blockSymbolPair, 0)

	// 选出最小的 blockIndex
	// （如果大小相同，可随机，这里先简单处理：先把 acceptBlock 加入队列）
	if len(acceptBlock) <= len(nonAcceptBlock) && len(acceptBlock) > 0 {
		for _, a := range dfa.Alphabet {
			Q = append(Q, blockSymbolPair{0, a})
		}
	} else if len(nonAcceptBlock) > 0 {
		// nonAcceptBlock
		idx := 0
		if len(acceptBlock) > 0 {
			idx = 1
		}
		for _, a := range dfa.Alphabet {
			Q = append(Q, blockSymbolPair{idx, a})
		}
	}

	// 4. 主循环：当队列不空时
	for len(Q) > 0 {
		// 取队首
		pair := Q[0]
		Q = Q[1:]
		currentBlock := pair.blockIndex
		a := pair.symbol

		// 找到将会转到 currentBlock 的那些状态
		// 即对 block partition[currentBlock] 的每个状态 t， revTrans[a][t] 是它的前驱状态集
		// 这些前驱状态散落在各个 block，需要进行拆分
		var involvedStates []int
		for _, t := range partition[currentBlock] {
			// 找到前驱
			if list, ok := revTrans[a][t]; ok {
				involvedStates = append(involvedStates, list...)
			}
		}
		// 以 block 为单位进行拆分
		// 我们要对 involvedStates 所在的块进行拆分
		blockChanged := make(map[int][]int) // oldBlockIndex -> statesInThatBlockNeedSplit
		for _, s := range involvedStates {
			oldBid := blockID[s]
			blockChanged[oldBid] = append(blockChanged[oldBid], s)
		}

		for oldB, subset := range blockChanged {
			// subset 的状态需要从 oldB 对应的 block 中拆分出来
			if len(subset) == len(partition[oldB]) {
				// 整个 block 都是这些状态，无需拆分
				continue
			}
			// 否则，我们拆出 subset 形成新块
			newBlockIndex := len(partition)
			newBlock := subset
			// 在原 block 中删除 subset
			oldBlock := partition[oldB]

			remainBlock := make([]int, 0, len(oldBlock)-len(newBlock))
			inSubset := make(map[int]bool)
			for _, stt := range newBlock {
				inSubset[stt] = true
			}
			for _, stt := range oldBlock {
				if !inSubset[stt] {
					remainBlock = append(remainBlock, stt)
				}
			}
			partition[oldB] = remainBlock
			partition = append(partition, newBlock)

			// 更新 blockID
			for _, stt := range newBlock {
				blockID[stt] = newBlockIndex
			}

			// 将“较小的那个” block 加入队列
			if len(newBlock) < len(remainBlock) {
				for _, x := range dfa.Alphabet {
					Q = append(Q, blockSymbolPair{newBlockIndex, x})
				}
			} else {
				for _, x := range dfa.Alphabet {
					Q = append(Q, blockSymbolPair{oldB, x})
				}
			}
		}
	}

	// 5. 构造最小化 DFA
	// partition 中每个 block 视为一个新状态
	minStates := make([]int, len(partition))
	for i := range minStates {
		minStates[i] = i
	}

	// 找到新初始态
	// 原来的初始态 dfa.Start 所在 blockID 即是最小化后的初始态
	newStart := blockID[dfa.Start]

	// 找到新接受态集
	newAccept := make(map[int]bool)
	for i, blk := range partition {
		for _, s := range blk {
			if dfa.AcceptStates[s] {
				newAccept[i] = true
				break
			}
		}
	}

	// 构造新的转移
	newTrans := make(map[int]map[rune]int)
	for i := range partition {
		newTrans[i] = make(map[rune]int)
	}
	// 对 partition[i] 中的任意一个代表状态 s，来构造转移即可
	for i, blk := range partition {
		if len(blk) == 0 {
			// 空块，理论上不会出现，若出现可视为死状态
			continue
		}
		rep := blk[0] // 代表状态
		for _, a := range dfa.Alphabet {
			t := dfa.Transition[rep][a]
			newTrans[i][a] = blockID[t]
		}
	}

	// 返回新 DFA
	minDFA := &DFA{
		States:       minStates,
		Alphabet:     dfa.Alphabet,
		Transition:   newTrans,
		Start:        newStart,
		AcceptStates: newAccept,
	}
	return minDFA
}

// 测试函数
func main() {
	// 举个小示例
	// 构造一个简单的 DFA
	// 状态 0,1,2,3,4； 初始态 0； 接受态 2,4
	// 字母表 {0,1}
	dfa := &DFA{
		States:       []int{0, 1, 2, 3, 4},
		Alphabet:     []rune{'0', '1'},
		Transition:   make(map[int]map[rune]int),
		Start:        0,
		AcceptStates: map[int]bool{2: true, 4: true},
	}
	// 初始化转移
	for _, s := range dfa.States {
		dfa.Transition[s] = make(map[rune]int)
	}
	// 这里随便写一些转移关系（仅示例）
	dfa.Transition[0]['0'] = 1
	dfa.Transition[0]['1'] = 2
	dfa.Transition[1]['0'] = 1
	dfa.Transition[1]['1'] = 3
	dfa.Transition[2]['0'] = 1
	dfa.Transition[2]['1'] = 2
	dfa.Transition[3]['0'] = 4
	dfa.Transition[3]['1'] = 1
	dfa.Transition[4]['0'] = 4
	dfa.Transition[4]['1'] = 2

	fmt.Println("Original DFA:")
	fmt.Println("States:", dfa.States)
	fmt.Println("Start:", dfa.Start)
	fmt.Println("Accept:", dfa.AcceptStates)
	fmt.Println("Transition:", dfa.Transition)
	fmt.Println()

	// 最小化
	minDFA := HopcroftMinimize(dfa)
	fmt.Println("Minimized DFA:")
	fmt.Println("States:", minDFA.States)
	fmt.Println("Start:", minDFA.Start)
	fmt.Println("Accept:", minDFA.AcceptStates)
	fmt.Println("Transition:", minDFA.Transition)
}
