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
# dp[i]=min(dp[i],dp[j+1]+cost) 当target[i:j+1]出现在字典里.
# !解法：短串插入到trie中匹配(哈希表也可以)，长串个数不多，使用kmp预处理.


from typing import List, Optional, Sequence, TypeVar

T = TypeVar("T", int, str)


def indexOfAll(
    longer: Sequence[T], shorter: Sequence[T], start=0, nexts: Optional[List[int]] = None
) -> List[int]:
    """kmp O(n+m)求搜索串 `longer` 中所有匹配 `shorter` 的位置."""
    if not shorter:
        return []
    if len(longer) < len(shorter):
        return []
    res = []
    next = getNext(shorter) if nexts is None else nexts
    hitJ = 0
    for i in range(start, len(longer)):
        while hitJ > 0 and longer[i] != shorter[hitJ]:
            hitJ = next[hitJ - 1]
        if longer[i] == shorter[hitJ]:
            hitJ += 1
        if hitJ == len(shorter):
            res.append(i - len(shorter) + 1)
            hitJ = next[hitJ - 1]
    return res


def getNext(needle: Sequence[T]) -> List[int]:
    """kmp O(n)求 `needle`串的 `next`数组
    `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度
    """
    next = [0] * len(needle)
    j = 0
    for i in range(1, len(needle)):
        while j and needle[i] != needle[j]:  # 1. fallback后前进：匹配不成功j往右走
            j = next[j - 1]
        if needle[i] == needle[j]:  # 2. 匹配：匹配成功j往右走一步
            j += 1
        next[i] = j
    return next


class TrieNode:
    __slots__ = ("children", "cost")

    def __init__(self):
        self.children = dict()
        self.cost = INF

    def insert(self, word: str, cost: int) -> None:
        node = self
        for c in word:
            if c not in node.children:
                node.children[c] = TrieNode()
            node = node.children[c]
        node.cost = min2(node.cost, cost)


INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumCost(self, target: str, words: List[str], costs: List[int]) -> int:
        THRESHOLD = 500

        # !1.1 小串维护trie
        smallTrieRoot = TrieNode()
        for word, cost in zip(words, costs):
            if len(word) <= THRESHOLD:
                smallTrieRoot.insert(word, cost)

        # !1.2 大串kmp预处理
        n = len(target)
        bigTo = [[] for _ in range(n + 1)]
        nexts = getNext(target)
        for word, cost in zip(words, costs):
            if len(word) > THRESHOLD:
                starts = indexOfAll(target, word, 0, nexts)
                for start in starts:
                    bigTo[start].append((start + len(word), cost))

        dp = [INF] * (n + 1)
        dp[-1] = 0
        for i in range(n - 1, -1, -1):
            # !2.1 小串匹配
            curRoot = smallTrieRoot
            for j in range(i, min2(n, i + THRESHOLD)):
                if target[j] not in curRoot.children:
                    break
                curRoot = curRoot.children[target[j]]
                dp[i] = min2(dp[i], dp[j + 1] + curRoot.cost)

            # !2.2 大串匹配
            for to, cost in bigTo[i]:
                dp[i] = min2(dp[i], dp[to] + cost)

        return dp[0] if dp[0] != INF else -1
