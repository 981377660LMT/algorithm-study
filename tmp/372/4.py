from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的正整数数组 heights ，其中 heights[i] 表示第 i 栋建筑的高度。

# 如果一个人在建筑 i ，且存在 i < j 的建筑 j 满足 heights[i] < heights[j] ，那么这个人可以移动到建筑 j 。

# 给你另外一个数组 queries ，其中 queries[i] = [ai, bi] 。第 i 个查询中，Alice 在建筑 ai ，Bob 在建筑 bi 。


# 请你能返回一个数组 ans ，其中 ans[i] 是第 i 个查询中，Alice 和 Bob 可以相遇的 最左边的建筑 。如果对于查询 i ，Alice 和 Bob 不能相遇，令 ans[i] 为 -1 。

# function useBlock(
#   n: number,
#   blockSize = Math.sqrt(n + 1) | 0
# ): {
#   /** 下标所属的块. */
#   belong: Uint16Array
#   /** 每个块的起始下标(包含). */
#   blockStart: Uint32Array
#   /** 每个块的结束下标(不包含). */
#   blockEnd: Uint32Array
#   /** 块的数量. */
#   blockCount: number
# } {
#   const blockCount = 1 + ((n / blockSize) | 0)
#   const blockStart = new Uint32Array(blockCount)
#   const blockEnd = new Uint32Array(blockCount)
#   const belong = new Uint16Array(n)
#   for (let i = 0; i < blockCount; i++) {
#     blockStart[i] = i * blockSize
#     blockEnd[i] = Math.min((i + 1) * blockSize, n)
#   }
#   for (let i = 0; i < n; i++) {
#     belong[i] = (i / blockSize) | 0
#   }

#   return {
#     belong,
#     blockStart,
#     blockEnd,
#     blockCount
#   }
# }

INF = int(1e18)


def min(a: int, b: int) -> int:
    return a if a < b else b


def max(a: int, b: int) -> int:
    return a if a > b else b


def createBlock(
    n: int, blockSize: Optional[int] = None
) -> Tuple[int, List[int], List[int], List[int]]:
    if blockSize is None:
        blockSize = int(n**0.5) + 1
    blockCount = 1 + (n // blockSize)
    belong = [i // blockSize for i in range(n)]
    blockStart = [i * blockSize for i in range(blockCount)]
    blockEnd = [min((i + 1) * blockSize, n) for i in range(blockCount)]
    return blockCount, belong, blockStart, blockEnd


class Solution:
    def leftmostBuildingQueries(self, heights: List[int], queries: List[List[int]]) -> List[int]:
        def findRightNearestHigher(start: int, target: int) -> int:
            """找到[start, n)中第一个严格大于target的下标.如果不存在,返回-1."""
            bid = belong[start]
            for i in range(start, blockEnd[bid]):  # 散块
                if heights[i] > target:
                    return i
            for i in range(bid + 1, blockCount):  # 整块
                if blockMax[i] > target:
                    for j in range(blockStart[i], blockEnd[i]):
                        if heights[j] > target:
                            return j
            return -1

        n = len(heights)
        blockCount, belong, blockStart, blockEnd = createBlock(len(heights), 2 * int(n**0.5) + 1)
        blockMax = [INF] * blockCount
        for i, v in enumerate(heights):
            bid = belong[i]
            blockMax[bid] = max(blockMax[bid], v)
        res = [-1] * len(queries)
        for qi, (alice, bob) in enumerate(queries):
            if alice == bob:
                res[qi] = alice
                continue
            if alice > bob:
                alice, bob = bob, alice
            if heights[alice] < heights[bob]:
                res[qi] = bob
                continue
            res[qi] = findRightNearestHigher(bob, heights[alice])
        return res


# [1,4,2,3,5]
#

print(
    Solution().leftmostBuildingQueries(
        [1, 4, 2, 3, 5],
        [
            [1, 2],
        ],
    )
)
