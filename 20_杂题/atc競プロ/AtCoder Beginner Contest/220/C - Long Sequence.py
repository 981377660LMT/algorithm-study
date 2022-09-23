from itertools import accumulate
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def findMinK(nums: List[int], target: int) -> int:
    """找到最小的k使得nums[:k]的和大于target
    其中nums是一个无限长的循环正整数数组
    """

    def check(mid: int) -> bool:
        div, mod = divmod(mid, len(nums))
        return div * preSum[-1] + preSum[mod] > target

    preSum = [0] + list(accumulate(nums))
    left, right = 0, target
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
        else:
            left = mid + 1
    return left


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    target = int(input())
    print(findMinK(nums, target))
