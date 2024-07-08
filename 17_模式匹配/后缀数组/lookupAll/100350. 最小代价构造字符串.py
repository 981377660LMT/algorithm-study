# 单词拼接dp
# 100350. 最小代价构造字符串
# https://leetcode.cn/problems/construct-string-with-minimum-cost/description/
# 给你一个字符串 target、一个字符串数组 words 以及一个整数数组 costs，这两个数组长度相同。
# 设想一个空字符串 s。
# 你可以执行以下操作任意次数（包括零次）：
# 选择一个在范围  [0, words.length - 1] 的索引 i。
# 将 words[i] 追加到 s。
# 该操作的成本是 costs[i]。
# 返回使 s 等于 target 的 最小 成本。如果不可能，返回 -1。
#
# 1 <= target.length <= 5e4
# 所有 words[i].length 的总和小于或等于 5e4
#
# !1.预处理转移 dp from数组.
# !2.dp
# 时空复杂度O(nsqrtn)


from typing import List
from lookupAll import lookupAll, getSA


INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumCost(self, target: str, words: List[str], costs: List[int]) -> int:
        minCost = dict()
        for w, c in zip(words, costs):
            minCost[w] = min2(minCost.get(w, INF), c)

        n = len(target)
        from_ = [[] for _ in range(n + 1)]
        sa = getSA([ord(c) for c in target])
        for w, c in minCost.items():
            starts = lookupAll(target, sa, w)
            for start in starts:
                from_[start + len(w)].append(((start << 20) | c))

        dp = [INF] * (n + 1)
        dp[0] = 0
        MASK = (1 << 20) - 1
        for i in range(1, n + 1):
            for state in from_[i]:
                start, cost = state >> 20, state & MASK
                dp[i] = min2(dp[i], dp[start] + cost)
        return dp[n] if dp[n] < INF else -1
