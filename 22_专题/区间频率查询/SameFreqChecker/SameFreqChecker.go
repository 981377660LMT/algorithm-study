package main

type SameFreqCheckerAddOnly[K comparable] struct {
	counter map[K]int32
	maxFreq int32
	count   int32
}

func NewSameFreqCheckerAddOnly[K comparable]() *SameFreqCheckerAddOnly[K] {
	return &SameFreqCheckerAddOnly[K]{counter: make(map[K]int32)}
}

func (s *SameFreqCheckerAddOnly[K]) Add(num K) {
	pre := s.counter[num]
	s.counter[num] = pre + 1
	s.maxFreq = max32(s.maxFreq, pre+1)
	s.count++
}

func (s *SameFreqCheckerAddOnly[K]) Check() bool {
	return s.maxFreq*int32(len(s.counter)) == s.count
}

type SameFreqChecker[K comparable] struct {
	counter     map[K]int32
	freqCounter map[int32]int32
}

func NewSameFreqChecker[K comparable]() *SameFreqChecker[K] {
	return &SameFreqChecker[K]{counter: make(map[K]int32), freqCounter: make(map[int32]int32)}
}

func (s *SameFreqChecker[K]) Add(num K) {
	preC := s.counter[num]
	s.counter[num] = preC + 1
	s.freqCounter[preC+1] = s.freqCounter[preC+1] + 1
	if preC > 0 {
		preF := s.freqCounter[preC]
		if preF == 1 {
			delete(s.freqCounter, preC)
		} else {
			s.freqCounter[preC] = preF - 1
		}
	}
}

func (s *SameFreqChecker[K]) Discard(num K) bool {
	preC := s.counter[num]
	if preC == 0 {
		return false
	}
	if preC == 1 {
		delete(s.counter, num)
	} else {
		s.counter[num] = preC - 1
	}
	preF := s.freqCounter[preC]
	if preF == 1 {
		delete(s.freqCounter, preC)
	} else {
		s.freqCounter[preC] = preF - 1
	}
	if preC > 1 {
		s.freqCounter[preC-1] = s.freqCounter[preC-1] + 1
	}
	return true
}

func (s *SameFreqChecker[K]) Check() bool {
	return len(s.freqCounter) == 1
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
