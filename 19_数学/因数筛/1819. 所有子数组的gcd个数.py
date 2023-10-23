from typing import Callable, DefaultDict, List, Optional
from collections import defaultdict
from math import gcd


def countGcdOfAllSubarray(nums: List[int]) -> int:
    """返回数组所有子数组的 gcd 的不同个数."""
    return len(logTrick(nums, gcd))


def logTrick(
    nums: List[int],
    op: Callable[[int, int], int],
    f: Optional[Callable[[int, DefaultDict[int, int]], None]] = None,
) -> DefaultDict[int, int]:
    """
    将 `nums` 的所有非空子数组的元素进行 `op` 操作，返回所有不同的结果和其出现次数.

    Args:
        nums: 1 <= len(nums) <= 1e5.
        op: 与/或/gcd/lcm 中的一种操作，具有单调性.
        f: `nums[:end]` 中所有子数组的结果为 `preCounter`.

    Returns:
        所有不同的结果和其出现次数
    """
    res = defaultdict(int)
    dp = []
    for pos, cur in enumerate(nums):
        for i in range(len(dp)):
            dp[i] = op(dp[i], cur)
        dp.append(cur)

        # 去重
        ptr = 0
        for i in range(1, len(dp)):
            if dp[i] != dp[ptr]:
                ptr += 1
                dp[ptr] = dp[i]

        dp = dp[: ptr + 1]
        for v in dp:
            res[v] += 1
        if f is not None:
            f(pos + 1, res)

    return res


if __name__ == "__main__":
    print(countGcdOfAllSubarray([6, 10, 3]))
