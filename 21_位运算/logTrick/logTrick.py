from itertools import accumulate
from math import gcd
from typing import Callable, DefaultDict, List, Optional, Tuple
from collections import defaultdict


def logTrick(
    nums: List[int],
    op: Callable[[int, int], int],
    f: Optional[Callable[[List[Tuple[int, int, int]], int], None]] = None,
) -> DefaultDict[int, int]:
    """
    将 `nums` 的所有非空子数组的元素进行 `op` 操作，返回所有不同的结果和其出现次数.

    Args:
        nums: 1 <= len(nums) <= 1e5.
        op: 与/或/gcd/lcm 中的一种操作，具有单调性.
        f: (interval: List[Tuple[int, int, int]], right: int) -> None
        数组的右端点为right.
        interval 的 leftStart/leftEnd 表示子数组的左端点left的范围.
        interval 的 value 表示该子数组 arr[left,right] 的 op 结果.

    Returns:
        所有不同的结果和其出现次数
    """
    res = defaultdict(int)
    dp = []
    for pos, cur in enumerate(nums):
        for v in dp:
            v[2] = op(v[2], cur)
        dp.append([pos, pos + 1, cur])

        ptr = 0
        for v in dp[1:]:
            if dp[ptr][2] != v[2]:
                ptr += 1
                dp[ptr] = v
            else:
                dp[ptr][1] = v[1]
        dp = dp[: ptr + 1]

        for v in dp:
            res[v[2]] += v[1] - v[0]
        if f is not None:
            f(dp, pos)

    return res


if __name__ == "__main__":
    # 1521. 找到最接近目标值的函数值
    class Solution2:
        def closestToTarget(self, arr: List[int], target: int) -> int:
            counter = logTrick(arr, lambda x, y: x & y)
            return min(abs(k - target) for k in counter)

    # 2941. 子数组的最大 GCD-Sum
    # https://leetcode.cn/problems/maximum-gcd-sum-of-a-subarray/description/
    class Solution:
        def maxGcdSum(self, nums: List[int], k: int) -> int:
            def f(interval: List[Tuple[int, int, int]], right: int) -> None:
                nonlocal res
                for start, _, gcd_ in interval:
                    len_ = right - start + 1
                    if len_ >= k:
                        res = max(res, gcd_ * (preSum[right + 1] - preSum[start]))

            res = 0
            preSum = [0] + list(accumulate(nums))
            logTrick(nums, gcd, f)
            return res

    # 3605. 数组的最小稳定性因子
    # 如果一个 子数组 的所有元素的最大公因数（gcd） 大于或等于 2，则称该子数组是稳定的。
    # 给定上界 mid，能否通过至多 maxC 次修改，让 nums 的稳定性因子（最长稳定子数组的长度）不超过 mid？
    class Solution:
        def minStable(self, nums: List[int], maxC: int) -> int:
            n = len(nums)
            # !对于每个 i，以 i 为右端点的子数组 GCD ≥2 时，子数组的左端点的最小值
            leftMin = [n] * n

            def f(interval: List[Tuple[int, int, int]], right: int) -> None:
                for start, _, gcd_ in interval:
                    if gcd_ >= 2:
                        leftMin[right] = min(leftMin[right], start)

            logTrick(nums, gcd, f)

            def check(mid: int) -> bool:
                remain = maxC
                right = mid
                while right < n:
                    if right - leftMin[right] + 1 > mid:
                        if remain == 0:
                            return False
                        remain -= 1
                        right += mid + 1
                    else:
                        right += 1
                return True

            left, right = 0, n
            while left <= right:
                mid = (left + right) // 2
                if check(mid):
                    right = mid - 1
                else:
                    left = mid + 1
            return left
