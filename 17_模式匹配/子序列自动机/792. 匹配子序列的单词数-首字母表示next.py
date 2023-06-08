# !求 words[i] 中是 S 的子序列的单词个数。


from typing import List
from collections import defaultdict
from SubsequenceAutomaton import SubsequenceAutomaton1, SubsequenceAutomaton2


class Solution:
    def numMatchingSubseq(self, s: str, words: List[str]) -> int:
        """子序列自动机解法."""
        SA = SubsequenceAutomaton1(s)
        return sum(SA.includes(word) for word in words)

    def numMatchingSubseq2(self, s: str, words: List[str]) -> int:
        # !指向下一个字母的指针
        # 因为 S 很长，所以寻找一种只需遍历一次 S 的方法，避免暴力解法的多次遍历。
        # !将所有单词根据首字母不同放入不同的桶中
        # 每个桶中的单词就是该单词正在等待匹配的下一个字母
        wordsByHead = defaultdict(list)
        for word in words:
            wordsByHead[word[0]].append(word)

        res = 0
        for char in s:
            matches = wordsByHead[char]
            wordsByHead[char] = []
            for word in matches:
                if len(word) == 1:
                    res += 1
                else:
                    wordsByHead[word[1]].append(word[1:])

        return res


print(Solution().numMatchingSubseq(s="abcde", words=["a", "bb", "acd", "ace"]))
