# 1 <= words.length <= 1000
# 1 <= words[i].length <= 1000
# words[i] 由小写英文字母组成


from collections import defaultdict
from typing import List

from Trie import Trie


class Solution:
    def sumPrefixScores(self, words: List[str]) -> List[int]:
        """1. 字典树"""
        trie = Trie(words)
        res = []
        for w in words:
            res.append(sum(trie.countWordStartsWith(w)))
        return res

    def sumPrefixScores2(self, words: List[str]) -> List[int]:
        """2. 暴力切片

        # !python 暴力 4400ms js 暴力 9000ms  python 字符串切片比js快很多
        # !字符串长度<=1000时 字符串切片可以当作常数时间
        """
        mp = defaultdict(int)
        for word in words:
            for i in range(len(word)):
                mp[word[: i + 1]] += 1

        res = []
        for word in words:
            count = 0
            for i in range(len(word)):
                count += mp[word[: i + 1]]
            res.append(count)
        return res

    def sumPrefixScores3(self, words: List[str]) -> List[int]:
        """3. 双哈希优化暴力切片"""
        base1, mod1, base2, mod2 = 131, 13331, 13331, int(1e9 + 7)
        mp = defaultdict(int)
        for word in words:
            hash1, hash2 = 0, 0
            for char in word:
                ord_ = ord(char)
                hash1 = (hash1 * base1 + ord_) % mod1
                hash2 = (hash2 * base2 + ord_) % mod2
                mp[(hash1, hash2)] += 1

        res = []
        for word in words:
            hash1, hash2 = 0, 0
            count = 0
            for char in word:
                ord_ = ord(char)
                hash1 = (hash1 * base1 + ord_) % mod1
                hash2 = (hash2 * base2 + ord_) % mod2
                count += mp[(hash1, hash2)]
            res.append(count)
        return res
