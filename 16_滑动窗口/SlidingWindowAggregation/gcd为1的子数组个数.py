# gcd为1的子数组个数


from SlidingWindowAggregation import SlidingWindowAggregation

from typing import List
from math import gcd


def solve(nums: List[int]) -> int:
    n = len(nums)
    S = SlidingWindowAggregation(lambda: 0, gcd)
    res = 0
    right = 0
    for left in range(n):
        right = max(right, left)
        while right < n and gcd(S.query(), nums[right]) != 1:
            S.append(nums[right])
            right += 1
        res += n - right
        S.popleft()
    return res


if __name__ == "__main__":
    from random import randint
    import time

    n = int(1e5)
    nums = [randint(1, int(1e6)) for _ in range(n)]
    time1 = time.time()
    print(solve(nums))
    print(time.time() - time1)
