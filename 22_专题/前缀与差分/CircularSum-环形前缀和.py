# 环形数组前缀和/环形区间前缀和


from itertools import accumulate
from random import randint
from typing import Callable, List


def circularPresum(nums: List[int]) -> Callable[[int, int], int]:
    """环形数组前缀和
    https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
    """
    n = len(nums)
    preSum = [0] + list(accumulate(nums))

    def cal(r: int) -> int:
        return preSum[n] * (r // n) + preSum[r % n]

    def query(start: int, end: int) -> int:
        """[start,end)的和"""
        return cal(end) - cal(start)

    return query


if __name__ == "__main__":
    nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    cs = circularPresum(nums)
    for _ in range(100):
        # check with bf
        left, right = randint(0, 100), randint(0, 100)
        if left > right:
            left, right = right, left
        sum1 = cs(left, right)
        sum2 = sum(nums[i] for i in range(left, right))
        assert sum1 == sum2, (left, right, sum1, sum2)
