from typing import List
from sortedcontainers import SortedList


class BitTree:
    def __init__(self, n: int):
        self.tree = [0 for _ in range(n + 1)]
        self.n = n

    def add(self, x: int, diff: int) -> None:
        if x <= 0:
            raise ValueError

        while x <= self.n:
            self.tree[x] += diff
            x += self.lowbit(x)

    def query(self, x: int) -> int:
        presum = 0
        while x > 0:
            presum += self.tree[x]
            x -= self.lowbit(x)

        return presum

    def lowbit(self, x: int) -> int:
        return x & (-x)


# 1 <= N <= 20000
# 最初，所有灯泡都关闭。每天只打开一个灯泡，直到 N 天后所有灯泡都打开。
# ，其中 bulls[i] = x 意味着在第 (i+1) 天，我们会把在位置 x 的灯泡打开，其中 i 从 0 开始，x 从 1 开始。

# 请你输出在第几天恰好有两个打开的灯泡，使得它们中间 正好 有 K 个灯泡且这些灯泡 全部是关闭的 。

# 有序集合找左/右边第一个亮灯泡就可以了
class Solution:
    def kEmptySlots(self, bulbs: List[int], k: int) -> int:
        lights = SortedList()
        for i, cur in enumerate(bulbs):
            lights.add(cur)
            index = lights.bisect_left(cur)

            if index + 1 < len(lights):
                rightPos = lights[index + 1]
                if rightPos - cur == k + 1:
                    return i + 1

            if index > 0:
                leftPos = lights[index - 1]
                if cur - leftPos == k + 1:
                    return i + 1

        return -1

    # 树状数组 查询区间范围值
    def kEmptySlots2(self, bulbs: List[int], k: int) -> int:
        n = len(bulbs)
        bit = BitTree(n + 1)
        isLighted = [False] * (n + 1)

        for i, pos in enumerate(bulbs):
            bit.add(pos, 1)
            isLighted[pos] = True

            leftCand = pos - (k + 1)
            rightCand = pos + (k + 1)

            if 0 <= leftCand and isLighted[leftCand]:
                if bit.query(pos - 1) == bit.query(leftCand):
                    return i + 1

            if rightCand <= n and isLighted[rightCand]:
                if bit.query(rightCand - 1) == bit.query(pos):
                    return i + 1

        return -1


print(Solution().kEmptySlots2([1, 3, 2], 1))
print(Solution().kEmptySlots2([6, 5, 8, 9, 7, 1, 10, 2, 3, 4], 2))

# 输出：2
# 解释：
# 第一天 bulbs[0] = 1，打开第一个灯泡 [1,0,0]
# 第二天 bulbs[1] = 3，打开第三个灯泡 [1,0,1]
# 第三天 bulbs[2] = 2，打开第二个灯泡 [1,1,1]
# 返回2，因为在第二天，两个打开的灯泡之间恰好有一个关闭的灯泡。

