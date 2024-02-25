from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个下标从 1 开始的整数数组 nums 和 changeIndices ，数组的长度分别为 n 和 m 。

# 一开始，nums 中所有下标都是未标记的，你的任务是标记 nums 中 所有 下标。

# 从第 1 秒到第 m 秒（包括 第 m 秒），对于每一秒 s ，你可以执行以下操作 之一 ：


# 选择范围 [1, n] 中的一个下标 i ，并且将 nums[i] 减少 1 。
# 如果 nums[changeIndices[s]] 等于 0 ，标记 下标 changeIndices[s] 。
# 什么也不做。
# 请你返回范围 [1, m] 中的一个整数，表示最优操作下，标记 nums 中 所有 下标的 最早秒数 ，如果无法标记所有下标，返回 -1 。


class Solution:
    def earliestSecondToMarkIndices(self, nums: List[int], changeIndices: List[int]) -> int:
        # 二分+倒序遍历
        def check(mid: int) -> bool:
            curIndices = changeIndices[:mid]
            curNeed = nums[:]
            marked = [False] * n
            score = 0
            mp = defaultdict(int)
            for v in curIndices:
                mp[v] += 1
            for v in curIndices:
                mp[v] -= 1
                if mp[v] == 0:  # last
                    if curNeed[v] > score:
                        return False
                    score -= curNeed[v]
                    marked[v] = True
                else:
                    score += 1

            return all(marked)

        changeIndices = [x - 1 for x in changeIndices]
        n, m = len(nums), len(changeIndices)
        left, right = 0, m
        ok = False
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
                ok = True
            else:
                left = mid + 1
        return left if ok else -1


# [0,2,3,0]
# [2,4,1,3,3,3,3,3,3,2,1]
print(Solution().earliestSecondToMarkIndices([0, 2, 3, 0], [2, 4, 1, 3, 3, 3, 3, 3, 3, 2, 1]))
