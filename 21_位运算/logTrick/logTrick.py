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
            v[0] = op(v[0], cur)
        dp.append([cur, pos, pos + 1])

        ptr = 0
        for v in dp[1:]:
            if dp[ptr][0] != v[0]:
                ptr += 1
                dp[ptr] = v
            else:
                dp[ptr][2] = v[2]

        dp = dp[: ptr + 1]
        for v in dp:
            res[v[0]] += v[2] - v[1]
        if f is not None:
            f(dp, pos)

    return res


if __name__ == "__main__":
    # 1521. 找到最接近目标值的函数值
    class Solution:
        def closestToTarget(self, arr: List[int], target: int) -> int:
            counter = logTrick(arr, lambda x, y: x & y)
            return min(abs(k - target) for k in counter)
