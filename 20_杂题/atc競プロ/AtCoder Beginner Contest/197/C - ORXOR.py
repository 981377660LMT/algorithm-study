from functools import reduce
from operator import xor
from typing import List

INF = int(1e18)


def orXor(nums: List[int]) -> int:
    """各个区间或的异或最小值"""

    def bt(index: int, path: List[int], curOr: int) -> None:
        nonlocal res
        if index == n:
            res = min(res, reduce(xor, path, curOr))
            return

        bt(index + 1, path, curOr | nums[index])
        path.append(curOr)
        bt(index + 1, path, nums[index])
        path.pop()

    n = len(nums)
    res = INF
    bt(0, [], 0)
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums = list(map(int, input().split()))
    print(orXor(nums))
