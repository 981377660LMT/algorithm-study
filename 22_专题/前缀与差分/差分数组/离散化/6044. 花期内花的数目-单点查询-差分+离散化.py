# 10^9值域 差分数组

# 6044. 花期内花的数目-单点查询-差分+离散化
from typing import List
from bisect import bisect_right
from collections import defaultdict
from itertools import accumulate


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


if __name__ == "__main__":
    print(
        Solution().fullBloomFlowers(
            flowers=[[1, 6], [3, 7], [9, 12], [4, 13]], persons=[2, 3, 7, 11]
        )
    )
