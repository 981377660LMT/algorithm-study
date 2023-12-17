from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 下标从 0 开始的整数数组 nums 。

# 你可以对 nums 执行特殊操作 任意次 （也可以 0 次）。每一次特殊操作中，你需要 按顺序 执行以下步骤：

# 从范围 [0, n - 1] 里选择一个下标 i 和一个 正 整数 x 。
# 将 |nums[i] - x| 添加到总代价里。
# 将 nums[i] 变为 x 。
# 如果一个正整数正着读和反着读都相同，那么我们称这个数是 回文数 。比方说，121 ，2552 和 65756 都是回文数，但是 24 ，46 ，235 都不是回文数。

# 如果一个数组中的所有元素都等于一个整数 y ，且 y 是一个小于 109 的 回文数 ，那么我们称这个数组是一个 等数数组 。


# 请你返回一个整数，表示执行任意次特殊操作后使 nums 成为 等数数组 的 最小 总代价。

from typing import Generator, Optional, Union


def emumeratePalindrome(
    minLength: int, maxLength: int, reverse=False
) -> Generator[str, None, None]:
    """
    遍历长度在 `[minLength, maxLength]` 之间的回文数字字符串.
    maxLength <= 12.
    """
    if minLength > maxLength:
        return
    if reverse:
        for length in reversed(range(minLength, maxLength + 1)):
            start = 10 ** ((length - 1) >> 1)
            end = start * 10 - 1
            for half in reversed(range(start, end + 1)):
                if length & 1:
                    yield f"{half}{str(half)[:-1][::-1]}"
                else:
                    yield f"{half}{str(half)[::-1]}"
    else:
        for length in range(minLength, maxLength + 1):
            start = 10 ** ((length - 1) >> 1)
            end = start * 10 - 1
            for half in range(start, end + 1):
                if length & 1:
                    yield f"{half}{str(half)[:-1][::-1]}"
                else:
                    yield f"{half}{str(half)[::-1]}"


palindromes = [int(v) for v in emumeratePalindrome(1, 9)]

from typing import Callable, List
from itertools import accumulate
from bisect import bisect_right


def distSum(sortedNums: List[int]) -> Callable[[int], int]:
    """`有序数组`所有点到x=k的距离之和

    排序+二分+前缀和 O(logn)
    """
    preSum = [0] + list(accumulate(sortedNums))

    def query(k: int) -> int:
        pos = bisect_right(sortedNums, k)
        leftSum = k * pos - preSum[pos]
        rightSum = preSum[-1] - preSum[pos] - k * (len(sortedNums) - pos)
        return leftSum + rightSum

    return query


def fibonacciSearch(f: Callable[[int], int], left: int, right: int, min: bool) -> Tuple[int, int]:
    """斐波那契搜索寻找[left,right]中的一个极值点,不要求单峰性质.
    Args:
        f: 目标函数.
        left: 搜索区间左端点(包含).
        right: 搜索区间右端点(包含).
        min: 是否寻找最小值.
    Returns:
        极值点的横坐标x和纵坐标f(x).
    """
    assert left <= right

    a, b, c, d = left, left + 1, left + 2, left + 3
    step = 0
    while d < right:
        b = c
        c = d
        d = b + c - a
        step += 1

    def get(i: int) -> int:
        if right < i:
            return INF
        return f(i) if min else -f(i)

    ya, yb, yc, yd = get(a), get(b), get(c), get(d)
    for _ in range(step):
        if yb < yc:
            d = c
            c = b
            b = a + d - c
            yd = yc
            yc = yb
            yb = get(b)
        else:
            a = b
            b = c
            c = a + d - b
            ya = yb
            yb = yc
            yc = get(c)

    x = a
    y = ya
    if yb < y:
        x = b
        y = yb
    if yc < y:
        x = c
        y = yc
    if yd < y:
        x = d
        y = yd

    return (x, y) if min else (x, -y)


class Solution:
    def minimumCost(self, nums: List[int]) -> int:
        nums.sort()
        Q = distSum(nums)
        return fibonacciSearch(
            lambda i: Q(palindromes[i]), left=0, right=len(palindromes) - 1, min=True
        )[1]
