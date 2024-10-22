package main

func main() {
	type E = int32
	counter := make(map[E]int32)

	add := func(x E) {
		preFreq := counter[x]
		counter[x]++
		if preFreq != 0 {
			// remove (preFreq, x)
		}
		// add (preFreq+1, x)
	}
	_ = add

	remove := func(x E) bool {
		preFreq := counter[x]
		if preFreq == 0 {
			return false
		}
		counter[x]--
		// remove (preFreq, x)
		if preFreq > 1 {
			// add (preFreq-1, x)
		} else {
			delete(counter, x)
		}
		return true
	}
	_ = remove
}
