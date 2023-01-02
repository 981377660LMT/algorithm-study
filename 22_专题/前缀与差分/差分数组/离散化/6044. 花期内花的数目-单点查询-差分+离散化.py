from typing import List
from bisect import bisect_right
from collections import defaultdict
from itertools import accumulate
from 紧离散化模板 import Discretizer


# 10^9值域 差分数组


class Solution:
    def fullBloomFlowers(self, flowers: List[List[int]], persons: List[int]) -> List[int]:
        """单点查询时:只对flowers离散化,开字典+二分查找query值被映射成啥"""
        diff = defaultdict(int)
        for left, right in flowers:
            diff[left] += 1
            diff[right + 1] -= 1

        # 离散化的keys、原数组前缀和
        keys = sorted(diff)
        diff = list(accumulate((diff[key] for key in keys), initial=0))
        return [diff[bisect_right(keys, p)] for p in persons]

    def fullBloomFlowers2(self, flowers: List[List[int]], persons: List[int]) -> List[int]:
        """单点查询时：如果同时也把person添加到离散化,就不用二分查找了/不用开字典了"""
        D = Discretizer()
        for left, right in flowers:
            D.add(left)
            D.add(right)
        for p in persons:
            D.add(p)
        D.build()

        diff = [0] * (len(D) + 10)
        for left, right in flowers:
            diff[D.get(left)] += 1
            diff[D.get(right) + 1] -= 1
        diff = list(accumulate(diff))
        return [diff[D.get(p)] for p in persons]


if __name__ == "__main__":

    print(
        Solution().fullBloomFlowers(
            flowers=[[1, 6], [3, 7], [9, 12], [4, 13]], persons=[2, 3, 7, 11]
        )
    )
    print(
        Solution().fullBloomFlowers2(
            flowers=[[1, 6], [3, 7], [9, 12], [4, 13]], persons=[2, 3, 7, 11]
        )
    )
