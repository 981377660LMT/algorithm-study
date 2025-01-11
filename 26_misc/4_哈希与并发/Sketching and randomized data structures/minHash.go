package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"strings"
	"time"
)

// Constants
const (
	numHashFunctions = 100 // Number of hash functions
	maxHashValue     = math.MaxUint32
)

// MinHash struct
type MinHash struct {
	numHash uint32
	hashA   []uint32
	hashB   []uint32
}

// NewMinHash initializes a MinHash with a given number of hash functions
func NewMinHash(num uint32) *MinHash {
	mh := &MinHash{
		numHash: num,
		hashA:   make([]uint32, num),
		hashB:   make([]uint32, num),
	}

	// Initialize hash function coefficients with random values
	rand.Seed(time.Now().UnixNano())
	for i := uint32(0); i < num; i++ {
		mh.hashA[i] = rand.Uint32() | 1 // Ensure it's odd
		mh.hashB[i] = rand.Uint32()
	}

	return mh
}

// hashToken hashes a token to a uint32 using FNV hash
func hashToken(token string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(token))
	return h.Sum32()
}

// ComputeSignature computes the MinHash signature for a given set of tokens
func (mh *MinHash) ComputeSignature(tokens []string) []uint32 {
	signature := make([]uint32, mh.numHash)
	for i := range signature {
		signature[i] = maxHashValue
	}

	for _, token := range tokens {
		tokenHash := hashToken(token)
		for i := uint32(0); i < mh.numHash; i++ {
			combinedHash := (mh.hashA[i]*tokenHash + mh.hashB[i]) % maxHashValue
			if combinedHash < signature[i] {
				signature[i] = combinedHash
			}
		}
	}

	return signature
}

// EstimateJaccard estimates the Jaccard similarity between two signatures
func EstimateJaccard(sig1, sig2 []uint32) float64 {
	if len(sig1) != len(sig2) {
		return 0.0
	}
	match := 0
	for i := 0; i < len(sig1); i++ {
		if sig1[i] == sig2[i] {
			match++
		}
	}
	return float64(match) / float64(len(sig1))
}

// Helper function to convert text to tokens (e.g., word-level shingles)
func tokenize(text string) []string {
	// Simple whitespace tokenizer; can be replaced with more sophisticated tokenization
	return strings.Fields(text)
}

func main() {
	// Initialize MinHash
	mh := NewMinHash(numHashFunctions)

	// Example documents
	doc1 := "The quick brown fox jumps over the lazy dog"
	doc2 := "The quick brown fox leaps over the lazy dog"
	doc3 := "Lorem ipsum dolor sit amet, consectetur adipiscing elit"

	// Tokenize documents
	tokens1 := tokenize(doc1)
	tokens2 := tokenize(doc2)
	tokens3 := tokenize(doc3)

	// Compute signatures
	sig1 := mh.ComputeSignature(tokens1)
	sig2 := mh.ComputeSignature(tokens2)
	sig3 := mh.ComputeSignature(tokens3)

	// Estimate similarities
	sim12 := EstimateJaccard(sig1, sig2)
	sim13 := EstimateJaccard(sig1, sig3)
	sim23 := EstimateJaccard(sig2, sig3)

	// Output results
	fmt.Printf("Jaccard similarity between Doc1 and Doc2: %.4f\n", sim12)
	fmt.Printf("Jaccard similarity between Doc1 and Doc3: %.4f\n", sim13)
	fmt.Printf("Jaccard similarity between Doc2 and Doc3: %.4f\n", sim23)
}
