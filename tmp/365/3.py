from bisect import bisect_left
from itertools import accumulate
from typing import Callable, List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的数组 nums 和一个整数 target 。

# 下标从 0 开始的数组 infinite_nums 是通过无限地将 nums 的元素追加到自己之后生成的。


# 请你从 infinite_nums 中找出满足 元素和 等于 target 的 最短 子数组，并返回该子数组的长度。如果不存在满足条件的子数组，返回 -1 。


def circularPresum(nums: List[int]) -> Callable[[int, int], int]:
    """环形数组前缀和."""
    n = len(nums)
    preSum = [0] + list(accumulate(nums))

    def _cal(r: int) -> int:
        return preSum[n] * (r // n) + preSum[r % n]

    def query(start: int, end: int) -> int:
        """[start,end)的和.
        0 <= start < end <= n.
        """
        if start >= end:
            return 0
        return _cal(end) - _cal(start)

    return query
class Solution:
    def minSizeSubarray(self, nums: List[int], target: int) -> int:
        n = len(nums)
        
        m, target = divmod(target, sum(nums))
        
        if target == 0:
            return n * m
        
        res = inf
        subsum, left = 0, 0
        
        for right in range(2 * n):
            subsum += nums[right % n]
            
            while subsum >= target:
                if subsum == target:
                    res = min(res, right - left + 1)
                subsum -= nums[left % n]
                left += 1
                
        return res + n * m if res < inf else -1

作者：我爱志方小姐
链接：https://leetcode.cn/circle/discuss/g1ACpe/view/p5j0bj/
来源：力扣（LeetCode）
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

class Solution:
    def minSizeSubarray(self, nums: List[int], target: int) -> int:
        Q = circularPresum(nums)
        # 枚举起点+二分
        res = INF
        for start in range(len(nums)):
            cand = bisect_left(
                range(int(1e9 + 10)), target, key=lambda mid: Q(start, start + mid) >= target
            )
            if Q(start, cand) == target:
                res = min(res, cand - start)

        return res if res != INF else -1


# nums = [1,2,3], target = 5
print(Solution().minSizeSubarray(nums=[1, 2, 3], target=5))
