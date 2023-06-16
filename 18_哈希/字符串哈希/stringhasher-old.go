package main

import "fmt"

const MOD1, MOD2 int = 1e8 + 7, 1e9 + 7
const BASE1, BASE2 int = 131, 13131
const OFFSET1, OFFSET2 int = 0, 0

// Hasher returns a function that can be used to hash a slice of the string.
// The returned function takes two indices, left and right,
// and returns the hash of the slice [left, right).
//
// It is based on the Rabin-Karp algorithm.
// The hash function is:
//   hash(s[left:right]) = ((s[left]-offset)*base^(right-left-1) +
//   (s[left+1]-offset)*base^(right-left-2) + ... + (s[right-1]-offset)) % mod
// where base is a prime number and mod is a prime number larger than the maximum value of a rune.
// offset is a constant that is subtracted from each rune to make it non-negative.
func StringHasher(ords []int, mod int, base int, offset int) func(left int, right int) int {
	prePow := make([]int, len(ords)+1)
	prePow[0] = 1
	preHash := make([]int, len(ords)+1)
	for i, v := range ords {
		prePow[i+1] = (prePow[i] * base) % mod
		preHash[i+1] = (preHash[i]*base + v - offset) % mod
	}

	sliceHash := func(left, right int) int {
		if left >= right {
			return 0
		}
		return (preHash[right] - preHash[left]*prePow[right-left]%mod + mod) % mod
	}

	return sliceHash
}

// In order to avoid hash collision, we can use two hash functions.
// Two strings are equal if and only if two hashes are equal.
func BiStringHasher(ords []int, mod1, mod2, base1, base2, offset1, offset2 int) func(left int, right int) (hash1, hash2 int) {
	hasher1 := StringHasher(ords, mod1, base1, offset1)
	hasher2 := StringHasher(ords, mod2, base2, offset2)

	sliceHash := func(left, right int) (hash1, hash2 int) {
		if left >= right {
			return 0, 0
		}

		hash1 = hasher1(left, right)
		hash2 = hasher2(left, right)
		return
	}

	return sliceHash
}

// Returns the hash of the concatenation of two strings.
func concatHash(hash1, hash2, len2, mod, base int) int {
	pow := func(base, exp, mod int) int {
		base %= mod
		res := 1
		for ; exp > 0; exp >>= 1 {
			if exp&1 == 1 {
				res = res * base % mod
			}
			base = base * base % mod
		}
		return res
	}

	return (hash1*pow(base, len2, mod) + hash2) % mod
}

func main() {
	s := "abcabc"
	ords := make([]int, len(s))
	for i, v := range s {
		ords[i] = int(v)
	}
	hasher := BiStringHasher(ords, MOD1, MOD2, BASE1, BASE2, OFFSET1, OFFSET2)
	fmt.Println(hasher(0, 3)) // 69609650
	fmt.Println(hasher(3, 6)) // 69609650
}
