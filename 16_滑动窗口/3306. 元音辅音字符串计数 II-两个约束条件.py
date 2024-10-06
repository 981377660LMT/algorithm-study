# 3306. 元音辅音字符串计数 II
# https://leetcode.cn/problems/count-of-substrings-containing-every-vowel-and-k-consonants-ii/description/
# 给你一个字符串 word 和一个 非负 整数 k。
# 返回 word 的 子字符串中，每个元音字母（'a'、'e'、'i'、'o'、'u'）至少 出现一次，并且 恰好 包含 k 个辅音字母的子字符串的总数。
#
# !注意这里转换成"至少"比"至多"更容易处理，两个条件的滑窗单调性相同.


from collections import defaultdict


def isVowel(c: str) -> bool:
    return c == "a" or c == "e" or c == "i" or c == "o" or c == "u"


class Solution:
    def countOfSubstrings(self, word: str, k: int) -> int:
        def calc(word: str, k: int) -> int:
            """包含所有元音, 且至少包含k个辅音字母."""
            vowelCounter = defaultdict(int)
            res, left, n = 0, 0, len(word)
            consonantCount = 0
            for right in range(n):
                if isVowel(word[right]):
                    vowelCounter[word[right]] += 1
                else:
                    consonantCount += 1
                while left <= right and len(vowelCounter) == 5 and consonantCount >= k:
                    if isVowel(word[left]):
                        vowelCounter[word[left]] -= 1
                        if vowelCounter[word[left]] == 0:
                            del vowelCounter[word[left]]
                    else:
                        consonantCount -= 1
                    left += 1
                res += left
            return res

        return calc(word, k) - calc(word, k + 1)
