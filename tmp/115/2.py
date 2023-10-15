from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def getWordsInLongestSubsequence(
        self, n: int, words: List[str], groups: List[int]
    ) -> List[str]:
        # dp[i] 表示以 words[i] 结尾的最长子序列长度
        dp = [1] * n
        pre = [-1] * n
        for i in range(1, n):
            for j in range(i):
                if groups[i] != groups[j]:
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


# ["wx","h"]
