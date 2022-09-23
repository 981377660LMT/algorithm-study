# 975. 奇偶跳.py
from typing import List, Tuple
from bisect import bisect_right
from collections import defaultdict

from sortedcontainers import SortedList


class Solution:
    def helper1(self, arr: List[int]) -> Tuple[List[int], List[int]]:
        """寻找每个元素右侧比自己大的里最小的和右侧比自己小的里最大的

        如果有多个符合题意，取右侧第一个
        """
        n = len(arr)
        ids = list(range(n))
        ids.sort(key=lambda i: (arr[i], i))

        # 右侧比自己大的里最小的（相同大的，取index小的(第一个)）index
        nextBigger = [-1] * n
        stack = []
        for id in ids:
            while stack and stack[-1] < id:
                nextBigger[stack.pop()] = id
            stack.append(id)

        ids.sort(key=lambda i: (-arr[i], i))
        # 右侧比自己小的里最大的（相同大的，取index小的(第一个)）index
        nextSmaller = [-1] * n
        stack = []
        for id in ids:
            while stack and stack[-1] < id:
                nextSmaller[stack.pop()] = id
            stack.append(id)

        return nextSmaller, nextBigger

    def helper2(self, nums: List[int]) -> Tuple[List[int], List[int]]:
        """有序集合寻找每个元素右侧比自己大的里最小的和右侧比自己小的里最大的

        相同大的,取index小的
        """
        nextSmaller, nextBigger = [-1] * len(nums), [-1] * len(nums)
        sl = SortedList((num, i) for i, num in enumerate(nums))
        indexMap = defaultdict(list)
        for i, num in enumerate(nums):
            indexMap[num].append(i)

        for i, num in enumerate(nums):
            sl.discard((num, i))

            pos1 = sl.bisect_right((num, int(1e20))) - 1  # 比自己小的里最大的，注意如果有相同元素，取index小的
            if pos1 >= 0:
                value = sl[pos1][0]
                pos = bisect_right(indexMap[value], i)
                nextSmaller[i] = indexMap[value][pos]

            pos2 = sl.bisect_left((num, int(-1e20)))  # 比自己大的里最小的，注意如果有相同元素，取index小的
            if pos2 < len(sl):
                value = sl[pos2][0]
                pos = bisect_right(indexMap[value], i)
                nextBigger[i] = indexMap[value][pos]

        return nextSmaller, nextBigger


if __name__ == "__main__":
    assert (
        Solution().helper1([10, 13, 12, 14, 15])
        == Solution().helper2([10, 13, 12, 14, 15])
        == ([-1, 2, -1, -1, -1], [2, 3, 3, 4, -1])
    )
    print(Solution().helper1([2, 3, 1, 1, 4]))
    print(Solution().helper2([2, 3, 1, 1, 4]))
    assert Solution().helper1([2, 3, 1, 1, 4]) == Solution().helper2([2, 3, 1, 1, 4])
