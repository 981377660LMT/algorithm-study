// https://github.com/siongui/go-succinct-data-structure-trie?tab=readme-ov-file
//
// 字典树结构非常适合快速查找字典单词，但如果字典的词汇量很大，构建的字典树可能会占用大量空间。
// 因此，简洁数据结构被应用于字典树结构，使我们能够同时实现快速查找和较小的空间需求。

// #region alphabet

package main

import "strings"

//var allowedCharacters = "abcdeghijklmnoprstuvyāīūṁṃŋṇṅñṭḍḷ…'’° -"
var allowedCharacters = "abcdefghijklmnopqrstuvwxyz "
var mapCharToUint = getCharToUintMap(allowedCharacters)
var mapUintToChar = getUintToCharMap(mapCharToUint)

/**
 * Write the data for each node, call getDataBits() to calculate how many bits
 * for one node.
 * 1 bit stores the "final" indicator. The other bits store one of the
 * characters of the alphabet.
 */
var dataBits = getDataBits(allowedCharacters)

func SetAllowedCharacters(alphabet string) {
	allowedCharacters = alphabet
	mapCharToUint = getCharToUintMap(alphabet)
	mapUintToChar = getUintToCharMap(mapCharToUint)
	dataBits = getDataBits(alphabet)
}

func getCharToUintMap(alphabet string) map[string]uint {
	result := map[string]uint{}

	var i uint = 0
	chars := strings.Split(alphabet, "")
	for _, char := range chars {
		result[char] = i
		i++
	}

	return result
}

func getUintToCharMap(c2ui map[string]uint) map[uint]string {
	result := map[uint]string{}
	for k, v := range c2ui {
		result[v] = k
	}
	return result
}

func getDataBits(alphabet string) uint {
	numOfChars := len(strings.Split(alphabet, ""))
	var i uint = 0

	for (1 << i) < numOfChars {
		i++
	}

	// one more bit for the "final" indicator
	return (i + 1)
}

// #endregion
