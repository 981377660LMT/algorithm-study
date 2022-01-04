from typing import List
from collections import defaultdict
from bisect import bisect_left

# 求 words[i] 中是 S 的子序列的单词个数。

# 经典题
# 2种解法


class Solution:
    def numMatchingSubseq(self, S: str, words: List[str]) -> int:
        # 392。判断word是否是S的子序列 哈希表记录出现位置
        # 每次寻找最左插入位置
        def is_sub(word: str, w_i: int, s_i: int) -> bool:
            if w_i == len(word):
                return True
            idx = idx_by_char[word[w_i]]
            if not idx or s_i > idx[-1]:
                return False

            left = idx[bisect_left(idx, s_i)]
            return is_sub(word, w_i + 1, left + 1)

        idx_by_char = defaultdict(list)
        for i, char in enumerate(S):
            idx_by_char[char].append(i)
        return sum(is_sub(word, 0, 0) for word in words)

    # 指向下一个字母的指针
    # 因为 S 很长，所以寻找一种只需遍历一次 S 的方法，避免暴力解法的多次遍历。
    # 将所有单词根据首字母不同放入不同的桶中
    # 每个桶中的单词就是该单词正在等待匹配的下一个字母
    def numMatchingSubseq2(self, S: str, words: List[str]) -> int:
        word_by_head = defaultdict(list)
        res = 0

        for word in words:
            word_by_head[word[0]].append(word)

        for char in S:
            match_word = word_by_head[char]
            word_by_head[char] = []
            for word in match_word:
                if len(word) == 1:
                    res += 1
                else:
                    word_by_head[word[1]].append(word[1:])

        return res


print(Solution().numMatchingSubseq(S="abcde", words=["a", "bb", "acd", "ace"]))
