from typing import List
from bisect import bisect_right

INF = 2 ** 63 - 1
MOD = int(1e9 + 7)


# 请你选择 最优 箱子供应商，使得 总浪费空间最小 。
# 1 <= n <= 105
# 模拟，对每个box找能装多少个packages


class Solution:
    def minWastedSpace(self, packages: List[int], boxes: List[List[int]]) -> int:
        packages.sort()
        res = INF

        for box in boxes:
            box.sort()
            if box[-1] < packages[-1]:
                continue

            # indexof模式搜索;找到这个盒子能装几个package
            pos = 0
            curSum = 0
            for b in box:
                prePos = pos
                pos = bisect_right(packages, b, lo=prePos)
                curSum += (pos - prePos) * b
                if pos == len(packages):
                    break

            res = min(res, curSum)

        return (res - sum(packages)) % MOD if res != INF else -1


print(Solution().minWastedSpace(packages=[2, 3, 5], boxes=[[4, 8], [2, 8]]))
# 输出：6
# 解释：选择第一个供应商最优，用两个尺寸为 4 的箱子和一个尺寸为 8 的箱子。
# 总浪费空间为 (4-2) + (4-3) + (8-5) = 6 。
