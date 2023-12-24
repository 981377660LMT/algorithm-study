from functools import lru_cache
from heapq import heappop, heappush
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


INF = int(1e20)

# 给你两个下标从 0 开始的字符串 source 和 target ，它们的长度均为 n 并且由 小写 英文字母组成。

# 另给你两个下标从 0 开始的字符串数组 original 和 changed ，以及一个整数数组 cost ，其中 cost[i] 代表将字符串 original[i] 更改为字符串 changed[i] 的成本。

# 你从字符串 source 开始。在一次操作中，如果 存在 任意 下标 j 满足 cost[j] == z  、original[j] == x 以及 changed[j] == y ，你就可以选择字符串中的 子串 x 并以 z 的成本将其更改为 y 。 你可以执行 任意数量 的操作，但是任两次操作必须满足 以下两个 条件 之一 ：

# 在两次操作中选择的子串分别是 source[a..b] 和 source[c..d] ，满足 b < c  或 d < a 。换句话说，两次操作中选择的下标 不相交 。
# 在两次操作中选择的子串分别是 source[a..b] 和 source[c..d] ，满足 a == c 且 b == d 。换句话说，两次操作中选择的下标 相同 。
# 返回将字符串 source 转换为字符串 target 所需的 最小 成本。如果不可能完成转换，则返回 -1 。

# 注意，可能存在下标 i 、j 使得 original[j] == original[i] 且 changed[j] == changed[i] 。


def min2(a, b):
    return a if a < b else b


# 将区间拆成若干段，求每一段转换的最小距离之和
from typing import Sequence

MOD = 10**11 + 7
BASE = 1313131


def useStringHasher(s: Sequence[str], mod=MOD, base=BASE):
    n = len(s)
    prePow = [1] * (n + 1)
    preHash = [0] * (n + 1)
    for i in range(1, n + 1):
        prePow[i] = (prePow[i - 1] * base) % mod
        preHash[i] = (preHash[i - 1] * base + ord(s[i - 1])) % mod

    def sliceHash(left: int, right: int):
        """切片 `s[left:right]` 的哈希值"""
        if left >= right:
            return 0
        left += 1
        return (preHash[right] - preHash[left - 1] * prePow[right - left + 1]) % mod

    return sliceHash


def getHash(s: str, mod=MOD, base=BASE):
    res = 0
    for c in s:
        res = (res * base + ord(c)) % mod
    return res


class Solution:
    def minimumCost(
        self, source: str, target: str, original: List[str], changed: List[str], cost: List[int]
    ) -> int:
        originId = [0] * len(original)
        changedId = [0] * len(changed)
        hashToId = defaultdict(lambda: len(hashToId))
        for i, s in enumerate(original):
            hash_ = getHash(s)
            originId[i] = hashToId[hash_]
        for i, s in enumerate(changed):
            hash_ = getHash(s)
            changedId[i] = hashToId[hash_]
        gSize = len(hashToId)
        adjList = [[INF] * gSize for _ in range(gSize)]
        for i in range(gSize):
            adjList[i][i] = 0
        for a, b, c in zip(originId, changedId, cost):
            adjList[a][b] = min2(adjList[a][b], c)
        for k in range(gSize):
            for i in range(gSize):
                for j in range(gSize):
                    adjList[i][j] = min2(adjList[i][j], adjList[i][k] + adjList[k][j])
        H1 = useStringHasher(source)
        H2 = useStringHasher(target)

        n = len(source)

        @lru_cache(None)
        def dfs(index: int) -> int:
            if index == n:
                return 0
            res = INF
            for j in range(index + 1, n + 1):
                hash_ = H1(index, j)
                targetHash = H2(index, j)
                if hash_ == targetHash:
                    res = min2(res, dfs(j))
                elif hash_ in hashToId and targetHash in hashToId:
                    res = min2(res, dfs(j) + adjList[hashToId[hash_]][hashToId[targetHash]])
            return res

        res = dfs(0)
        dfs.cache_clear()
        return res if res != INF else -1


# ：source = "abcd", target = "acbe", original = ["a","b","c","c","e","d"], changed = ["b","c","b","e","b","e"], cost = [2,5,5,1,2,20]
print(
    Solution().minimumCost(
        source="abcd",
        target="acbe",
        original=["a", "b", "c", "c", "e", "d"],
        changed=["b", "c", "b", "e", "b", "e"],
        cost=[2, 5, 5, 1, 2, 20],
    )
)
