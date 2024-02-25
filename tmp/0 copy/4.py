from bisect import bisect_left, bisect_right
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个下标从 1 开始的整数数组 nums 和 changeIndices ，数组的长度分别为 n 和 m 。

# 一开始，nums 中所有下标都是未标记的，你的任务是标记 nums 中 所有 下标。

# 从第 1 秒到第 m 秒（包括 第 m 秒），对于每一秒 s ，你可以执行以下操作 之一 ：


# 选择范围 [1, n] 中的一个下标 i ，并且将 nums[i] 减少 1 。
# 将 nums[changeIndices[s]] 设置成任意的 非负 整数。
# 选择范围 [1, n] 中的一个下标 i ， 满足 nums[i] 等于 0, 并 标记 下标 i 。
# 什么也不做。
# 请你返回范围 [1, m] 中的一个整数，表示最优操作下，标记 nums 中 所有 下标的 最早秒数 ，如果无法标记所有下标，返回 -1 。
def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def earliestSecondToMarkIndices(self, nums: List[int], changeIndices: List[int]) -> int:
        def check(mid: int) -> bool:
            curIndices = changeIndices[:mid]
            need = nums[:]
            marked = [False] * n
            score = 0
            counter = [0] * n
            last = [0] * n
            for i, v in enumerate(curIndices):
                counter[v] += 1
                last[v] = i
            if any(v == 0 for v in counter):
                return False
            print(1212)

            indexScore = set()

            for i, v in enumerate(curIndices):
                counter[v] -= 1
                if counter[v] == 0:  # last
                    if v in indexScore:
                        indexScore.remove(v)
                        marked[v] = True
                    else:
                        if need[v] <= score:
                            score -= need[v]
                            marked[v] = True
                        else:  # need to borrow indexScore
                            keys = sorted(indexScore, key=lambda x: -rangeCount(x, i, last[x]))
                            needV = need[v]
                            for k in keys:
                                score += 1
                                indexScore.remove(k)
                                if score >= needV:
                                    break
                            if score < needV:
                                return False
                else:
                    if need[v] <= 1:
                        score += 1
                    else:
                        if v not in indexScore:
                            indexScore.add(v)
                        else:
                            score += 1

            return all(marked)

        mp = defaultdict(list)
        for i, v in enumerate(changeIndices):
            mp[v].append(i)

        def rangeCount(v: int, left: int, right: int):
            indexes = mp[v]
            return bisect_right(indexes, right) - bisect_left(indexes, left)

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


# nums = [3,2,3], changeIndices = [1,3,2,2,2,2,3]
print(Solution().earliestSecondToMarkIndices([3, 2, 3], [1, 3, 2, 2, 2, 2, 3]))  # 3
