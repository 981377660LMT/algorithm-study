# 环形数组前缀和/环形区间前缀和


from random import randint
from typing import Callable, List


def groupPreSum(nums: List[int], mod: int) -> Callable[[int, int, int], int]:
    """模分组前缀和
    https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
    """
    preSum = [0] * (len(nums) + mod)
    for i, v in enumerate(nums):
        preSum[i + mod] = preSum[i] + v

    def cal(r: int, k: int) -> int:
        if r % mod <= k:
            return preSum[r // mod * mod + k]
        return preSum[(r + mod - 1) // mod * mod + k]

    def query(start: int, end: int, target: int) -> int:
        """区间[start,end)中下标模mod为target的元素的和"""
        target %= mod
        return cal(end, target) - cal(start, target)

    return query


if __name__ == "__main__":
    nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    gs = groupPreSum(nums, 3)
    for _ in range(100):
        # check with bf
        left, right, target = randint(0, 10), randint(0, 10), randint(0, 2)
        if left > right:
            left, right = right, left
        sum1 = gs(left, right, target)
        sum2 = sum(nums[i] for i in range(left, right) if i % 3 == target)

        assert sum1 == sum2, (left, right, target, sum1, sum2)
