from typing import List, Set, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 有一个大型的 (m - 1) x (n - 1) 矩形田地，其两个对角分别是 (1, 1) 和 (m, n) ，田地内部有一些水平栅栏和垂直栅栏，分别由数组 hFences 和 vFences 给出。

# 水平栅栏为坐标 (hFences[i], 1) 到 (hFences[i], n)，垂直栅栏为坐标 (1, vFences[i]) 到 (m, vFences[i]) 。

# 返回通过 移除 一些栅栏（可能不移除）所能形成的最大面积的 正方形 田地的面积，或者如果无法形成正方形田地则返回 -1。

# 由于答案可能很大，所以请返回结果对 109 + 7 取余 后的值。


# 注意：田地外围两个水平栅栏（坐标 (1, 1) 到 (1, n) 和坐标 (m, 1) 到 (m, n) ）以及两个垂直栅栏（坐标 (1, 1) 到 (m, 1) 和坐标 (1, n) 到 (m, n) ）所包围。这些栅栏 不能 被移除。
class Solution:
    def maximizeSquareArea(self, m: int, n: int, hFences: List[int], vFences: List[int]) -> int:
        def cal(arr: List[int]) -> Set[int]:
            res = set()
            for i in range(len(arr)):
                for j in range(i + 1, len(arr)):
                    res.add(arr[j] - arr[i])
            return res

        hFences.sort()
        vFences.sort()
        res1 = cal([1] + hFences[:] + [m])
        res2 = cal([1] + vFences[:] + [n])
        inter = res1 & res2
        if not inter:
            return -1
        return max(inter) ** 2 % MOD


# 3
# 9
# [2]
# [8,6,5,4]
