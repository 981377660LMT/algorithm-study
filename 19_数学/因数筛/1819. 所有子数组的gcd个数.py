from typing import Callable, DefaultDict, List, Optional, Tuple
from collections import defaultdict
from math import gcd


def countGcdOfAllSubarray(nums: List[int]) -> int:
    """返回数组所有子数组的 gcd 的不同个数."""
    return len(logTrick(nums, gcd))


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
    print(countGcdOfAllSubarray([6, 10, 3]))
