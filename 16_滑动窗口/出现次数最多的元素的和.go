// MaxFreq

package main

import "fmt"

func main() {
	ms := NewMajorSum()
	ms.Add(1)
	ms.Add(2)
	ms.Add(2)
	ms.Add(3)
	fmt.Println(ms.Query())
	ms.Add(3)
	ms.Add(3)
	fmt.Println(ms.Query())
	ms.Discard(3)
	fmt.Println(ms.Query())

	mf := NewMajorFreq()
	mf.Add(1)
	mf.Add(2)
	mf.Add(2)
	mf.Add(3)
	fmt.Println(mf.MaxFreq())
	mf.Discard(2)
	fmt.Println(mf.MaxFreq())
}

// 维护出现次数最多的元素的`出现次数`.
type MajorFreq struct {
	maxFreq   int         // 最大出现次数
	counter   map[int]int // 每个元素出现的次数
	freqTypes map[int]int // 每个出现次数的元素的种类数
}

func NewMajorFreq() *MajorFreq {
	return &MajorFreq{
		counter:   make(map[int]int),
		freqTypes: make(map[int]int),
	}
}

func (mf *MajorFreq) Add(x int) *MajorFreq {
	mf.counter[x]++
	xFreq := mf.counter[x]
	mf.freqTypes[xFreq]++
	mf.freqTypes[xFreq-1]--
	if xFreq > mf.maxFreq {
		mf.maxFreq = xFreq
	}
	return mf
}

func (mf *MajorFreq) Discard(x int) bool {
	if mf.counter[x] == 0 {
		return false
	}
	mf.counter[x]--
	xFreq := mf.counter[x]
	mf.freqTypes[xFreq]++
	mf.freqTypes[xFreq+1]--
	if xFreq+1 == mf.maxFreq && mf.freqTypes[mf.maxFreq] == 0 {
		mf.maxFreq--
	}
	if mf.counter[x] == 0 {
		delete(mf.counter, x)
	}
	return true
}

func (mf *MajorFreq) MaxFreq() int {
	return mf.maxFreq
}

// 维护出现次数最多的元素的和(多种元素出现次数一样最多也算).
type MajorSum struct {
	MaxFreq    int         // 最大出现次数
	MaxFreqSum int         // 最大出现次数的元素的和
	counter    map[int]int // 每个元素出现的次数
	freqSum    map[int]int // 每个出现次数的元素的和
	freqTypes  map[int]int // 每个出现次数的元素的种类数
}

func NewMajorSum() *MajorSum {
	return &MajorSum{
		counter:   make(map[int]int),
		freqSum:   make(map[int]int),
		freqTypes: make(map[int]int),
	}
}

func (ms *MajorSum) Add(x int) {
	ms.counter[x]++
	xFreq := ms.counter[x]
	ms.freqSum[xFreq] += x
	ms.freqSum[xFreq-1] -= x
	ms.freqTypes[xFreq]++
	ms.freqTypes[xFreq-1]--
	if xFreq > ms.MaxFreq {
		ms.MaxFreq = xFreq
		ms.MaxFreqSum = x
	} else if xFreq == ms.MaxFreq {
		ms.MaxFreqSum += x
	}
}

func (ms *MajorSum) Discard(x int) {
	if ms.counter[x] == 0 {
		return
	}
	ms.counter[x]--
	xFreq := ms.counter[x]
	ms.freqSum[xFreq] += x
	ms.freqSum[xFreq+1] -= x
	ms.freqTypes[xFreq]++
	ms.freqTypes[xFreq+1]--
	if xFreq+1 == ms.MaxFreq {
		ms.MaxFreqSum -= x
		if ms.freqTypes[ms.MaxFreq] == 0 {
			ms.MaxFreq--
			ms.MaxFreqSum = ms.freqSum[ms.MaxFreq]
		}
	}
	if ms.counter[x] == 0 {
		delete(ms.counter, x)
	}
}

func (ms *MajorSum) Query() int {
	return ms.MaxFreqSum
}
