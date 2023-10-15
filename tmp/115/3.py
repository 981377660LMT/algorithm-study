from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个整数 n 和一个下标从 0 开始的字符串数组 words ，和一个下标从 0 开始的 二进制 数组 groups ，两个数组长度都是 n 。

# 两个长度相等字符串的 汉明距离 定义为对应位置字符 不同 的数目。

# 你需要从下标 [0, 1, ..., n - 1] 中选出一个 最长子序列 ，将这个子序列记作长度为 k 的 [i0, i1, ..., ik - 1] ，它需要满足以下条件：

# 相邻 下标对应的 groups 值 不同。即，对于所有满足 0 < j + 1 < k 的 j 都有 groups[ij] != groups[ij + 1] 。
# 对于所有 0 < j + 1 < k 的下标 j ，都满足 words[ij] 和 words[ij + 1] 的长度 相等 ，且两个字符串之间的 汉明距离 为 1 。
# 请你返回一个字符串数组，它是下标子序列 依次 对应 words 数组中的字符串连接形成的字符串数组。如果有多个答案，返回任意一个。

# 子序列 指的是从原数组中删掉一些（也可能一个也不删掉）元素，剩余元素不改变相对位置得到的新的数组。

# 注意：words 中的字符串长度可能 不相等 。


class Solution:
    def getWordsInLongestSubsequence(
        self, n: int, words: List[str], groups: List[int]
    ) -> List[str]:
        # dp[i] 表示以 words[i] 结尾的最长子序列长度
        dp = [1] * n  # !dp初始值搞错了,TODO: LCS怎么dp的?
        pre = [-1] * n
        for i in range(1, n):
            for j in range(i):
                if (
                    groups[i] != groups[j]
                    and len(words[i]) == len(words[j])
                    and sum(a != b for a, b in zip(words[i], words[j])) == 1
                ):
                    if dp[j] + 1 > dp[i]:
                        dp[i] = dp[j] + 1
                        pre[i] = j

        maxLen = max(dp)
        maxIndex = dp.index(maxLen)
        res = []
        while maxIndex != -1:
            res.append(words[maxIndex])
            maxIndex = pre[maxIndex]
        return res[::-1]


# n = 3, words = ["bab","dab","cab"], groups = [1,2,2]
print(Solution().getWordsInLongestSubsequence(3, ["bab", "dab", "cab"], [1, 2, 2]))
print(Solution().getWordsInLongestSubsequence(3, ["bab", "dab", "c"], [1, 2, 2]))
print(Solution().getWordsInLongestSubsequence(1, ["bab"], [1]))
# 3
# ["bdb","aaa","ada"]
# [2,1,3]

# ["aaa","ada"]
print(Solution().getWordsInLongestSubsequence(3, ["bdb", "aaa", "ada"], [2, 1, 3]))
